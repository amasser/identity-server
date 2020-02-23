package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

// InvalidArgumentError indicates that an operation failed because
// of an invalid argument.
type InvalidArgumentError struct {
	Description string
}

func (iae *InvalidArgumentError) Error() string {
	return iae.Description
}

// MarshalJSON implements the json.Marshaler interface and is used
// by http.DefaultErrorEncoder
func (iae *InvalidArgumentError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"error": iae.Error(),
	})
}

// StatusCode returns http.StatusBadRequest and implements the
// StatusCoder interface of go-kit's http transport.
func (*InvalidArgumentError) StatusCode() int {
	return http.StatusBadRequest
}

// NewInvalidArgumentError returns a new invalid argument error
func NewInvalidArgumentError(descr string) error {
	return &InvalidArgumentError{Description: descr}
}

// IsInvalidArgument reports whether err is an invalid argument error or
// not. IsInvalidArgument returns false if err is nil.
func IsInvalidArgument(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(*InvalidArgumentError); ok {
		return true
	}

	return errors.Is(err, &InvalidArgumentError{})
}
