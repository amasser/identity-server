package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// NotFoundError is returned when a requested resource was not found
// or does not exist.
type NotFoundError struct {
	// ResourceKind is the kind of resource that was requested.
	ResourceKind string
}

func (nfe *NotFoundError) Error() string {
	return fmt.Sprintf("Request %s not found", nfe.ResourceKind)
}

// MarshalJSON implements the json.Marshaler interface and is
// used by http.DefaultErrorEncoder.
func (nfe *NotFoundError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"error": nfe.Error(),
	})
}

// StatusCode returns http.StatusNotFound and implements the
// StatusCoder interface of the http transport of go-kit.
func (*NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// IsNotFound reports whether err is a NotFoundError or not. IsNotFound
// returns false if err is nil.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	if _, ok := err.(*NotFoundError); ok {
		return true
	}

	return errors.Is(err, &NotFoundError{})
}

// NewNotFoundError returns new not-found error for resourceKind
func NewNotFoundError(resourceKind string) error {
	return &NotFoundError{ResourceKind: resourceKind}
}
