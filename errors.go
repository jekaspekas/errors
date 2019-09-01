package errors

import (
	"fmt"

	pkgerrors "github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

const (
	// NoType error
	NoType ErrorType = iota
	// BadRequest error
	BadRequest
	// NotFound error
	NotFound
	// AccessDenied error
	AccessDenied
)

type customError struct {
	errorType     ErrorType
	originalError error
	context       errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// New creates a new customError
func (errorType ErrorType) New(msg string) error {
	return customError{errorType: errorType, originalError: pkgerrors.New(msg)}
}

// Newf creates a new customError with formatted message
func (errorType ErrorType) Newf(format string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: fmt.Errorf(format, args...)}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrapf creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, format string, args ...interface{}) error {
	return customError{errorType: errorType, originalError: pkgerrors.Wrapf(err, format, args...)}
}

// Error returns the mssage of a customError
func (error customError) Error() string {
	return error.originalError.Error()
}

// New creates a no type error
func New(msg string) error {
	return customError{errorType: NoType, originalError: pkgerrors.New(msg)}
}

// Newf creates a no type error with formatted message
func Newf(format string, args ...interface{}) error {
	return customError{errorType: NoType, originalError: pkgerrors.New(fmt.Sprintf(format, args...))}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return pkgerrors.Cause(err)
}

// Wrapf an error with format string
func Wrapf(err error, format string, args ...interface{}) error {
	wrappedError := pkgerrors.Wrapf(err, format, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}

	return customError{errorType: NoType, originalError: wrappedError}
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(customError); ok {
		return customError{errorType: customErr.errorType, originalError: customErr.originalError, context: context}
	}

	return customError{errorType: NoType, originalError: err, context: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(customError); ok || customErr.context != emptyContext {

		return map[string]string{"field": customErr.context.Field, "message": customErr.context.Message}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}

	return NoType
}
