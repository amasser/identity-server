package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ConflictError is returned when a requested operation failed because
// of a resource conflict.
type ConflictError struct {
	// ConflictType is the kind of resource that caused a conflict.
	ConflictType string
}

func (ce *ConflictError) Error() string {
	return fmt.Sprintf("Detected %s conflict", ce.ConflictType)
}

// MarshalJSON implements the json.Marshaller interface and is
// used by http.DefaultErrorEncoder.
func (ce *ConflictError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"error": ce.Error(),
	})
}

// StatusCode returns http.StatusConflict and implements the
// StatusCoder interface of go-kit's http transport.
func (*ConflictError) StatusCode() int {
	return http.StatusConflict
}

// IsConflict reports whether err is a ConflictError or not. IsConflict
// returns false if err is nil.
func IsConflict(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, &ConflictError{})
}

// NewConflictError returns a new conflict error for the given
// type.
func NewConflictError(conflictType string) error {
	return &ConflictError{ConflictType: conflictType}
}
