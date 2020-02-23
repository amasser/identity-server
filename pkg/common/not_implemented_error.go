package common

import (
	"encoding/json"
	"net/http"
)

type notImplementedError struct{}

func (nie *notImplementedError) Error() string {
	return "not implemented"
}

// MarshalJSON implements the json.Marshaler interface and is used
// by http.DefaultErrorEncoder.
func (nie *notImplementedError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"error": nie.Error(),
	})
}

// StatusCode returns http.StatusNotImplemented and implements the
// StatusCoder interface of go-kit's http transport.
func (*notImplementedError) StatusCode() int {
	return http.StatusNotImplemented
}

// ErrNotImplemented indicates that a given operation has not been
// implemented.
var ErrNotImplemented = &notImplementedError{}
