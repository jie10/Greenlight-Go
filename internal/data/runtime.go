package data

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32

func (r *Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := strconv.Itoa(int(*r)) + " mins"

	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	unquotedJSONValue = strings.TrimSuffix(unquotedJSONValue, " mins")

	i, err := strconv.ParseInt(unquotedJSONValue, 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)
	return nil
}
