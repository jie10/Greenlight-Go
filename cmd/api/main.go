package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jie10/greenlight-go/internal/data"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn         string
		maxOpenCons int
		maxIdleCons int
		maxIdleTime time.Duration
	}
}

type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:admin123@localhost/greenlight?sslmode=disable", "PostgresSQL DSN")
	flag.IntVar(&cfg.db.maxOpenCons, "db-max-open-cons", 25, "PostgresSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleCons, "db-max-idle-cons", 25, "PostgresSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgresSQL max connection idle time")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}(db)

	logger.Info("database connection pool established")

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://./migrations", "greenlight", migrationDriver)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database migrations applied")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	db.SetMaxOpenConns(cfg.db.maxOpenCons)
	db.SetMaxIdleConns(cfg.db.maxIdleCons)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
