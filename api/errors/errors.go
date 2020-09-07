package errors

import "fmt"

// ServerError defines a custom error for a server including a message and a status code.
type ServerError struct {
	message string
	code    int
}

// Error implements the error interface
func (se *ServerError) Error() string {
	return fmt.Sprintf("message: %q, code: %d", se.message, se.code)
}

// Message gets the ServerError message.
func (se *ServerError) Message() string {
	return se.message
}

// Code gets the ServerError status code.
func (se *ServerError) Code() int {
	return se.code
}

// NewServerError creates a ServerError struct with the given parameters.
func NewServerError(message string, code int) *ServerError {
	return &ServerError{message: message, code: code}
}
