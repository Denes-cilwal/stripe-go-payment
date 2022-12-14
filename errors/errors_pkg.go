package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type responseError struct {
	httpErrorType HTTPErrorType
	originalError error
	contextInfo   errorContext
}

type errorContext struct {
	Field   string
	Message string
}

func (error responseError) Error() string {
	return error.originalError.Error()
}

func (errorType HTTPErrorType) New(msg string) error {
	return responseError{
		httpErrorType: errorType,
		originalError: errors.New(msg),
	}
}

func (errorType HTTPErrorType) Newf(msg string, args ...interface{}) error {
	return responseError{
		httpErrorType: errorType,
		originalError: fmt.Errorf(msg, args...),
	}
}

func (errorType HTTPErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

func (errorType HTTPErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	return responseError{
		httpErrorType: errorType,
		originalError: errors.Wrapf(err, msg, args...),
	}
}

func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if responseErr, ok := err.(responseError); ok {
		return responseError{
			httpErrorType: responseErr.httpErrorType,
			originalError: responseErr.originalError,
			contextInfo:   context,
		}
	}

	return responseError{
		httpErrorType: InternalError,
		originalError: err,
		contextInfo:   context,
	}
}

func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if responseErr, ok := err.(responseError); ok && responseErr.contextInfo != emptyContext {
		return map[string]string{
			"field":   responseErr.contextInfo.Field,
			"message": responseErr.contextInfo.Message,
		}
	}

	return nil
}

func GetErrorType(err error) HTTPErrorType {
	if responseErr, ok := err.(responseError); ok {
		return responseErr.httpErrorType
	}

	return InternalError
}
