package errors

import (
	"errors"
	"fmt"
)

func newErr(e error, message string, eType Typer) error {
	if e == nil {
		return &fundamental{
			msg:   message,
			eType: eType,
			stack: callers(),
		}
	}
	return &withMessage{
		cause: e,
		msg:   message,
		eType: eType,
	}
}

func newErrf(e error, eType Typer, format string, args ...interface{}) error {
	message := fmt.Sprintf(format, args...)
	return newErr(e, message, eType)
}

func getErrType(err error) Typer {
	e, _ := err.(interface {
		Type() Typer
	})
	if e == nil {
		return defaultErrType
	}
	return e.Type()
}

// Internal helper method for creating internal errors
func Internal(message string) error {
	return newErr(nil, message, TypeInternal)
}

// Internalf helper method for creating internal errors with formatted message
func Internalf(format string, args ...interface{}) error {
	return newErrf(nil, TypeInternal, format, args...)
}

// Validation is a helper function to create a new error of type TypeValidation
func Validation(message string) error {
	return newErr(nil, message, TypeValidation)
}

// Validationf is a helper function to create a new error of type TypeValidation, with formatted message
func Validationf(format string, args ...interface{}) error {
	return newErrf(nil, TypeValidation, format, args...)
}

// Input is a helper function to create a new error of type TypeInput
func Input(message string) error {
	return newErr(nil, message, TypeInput)
}

// Inputf is a helper function to create a new error of type TypeInput, with formatted message
func Inputf(format string, args ...interface{}) error {
	return newErrf(nil, TypeInput, format, args...)
}

// Duplicate is a helper function to create a new error of type TypeDuplicate
func Duplicate(message string) error {
	return newErr(nil, message, TypeDuplicate)
}

// Duplicatef is a helper function to create a new error of type TypeDuplicate, with formatted message
func Duplicatef(format string, args ...interface{}) error {
	return newErrf(nil, TypeDuplicate, format, args...)
}

// Unauthenticated is a helper function to create a new error of type TypeUnauthenticated
func Unauthenticated(message string) error {
	return newErr(nil, message, TypeUnauthenticated)
}

// Unauthenticatedf is a helper function to create a new error of type TypeUnauthenticated, with formatted message
func Unauthenticatedf(format string, args ...interface{}) error {
	return newErrf(nil, TypeUnauthenticated, format, args...)

}

// Unauthorized is a helper function to create a new error of type TypeUnauthorized
func Unauthorized(message string) error {
	return newErr(nil, message, TypeUnauthorized)
}

// Unauthorizedf is a helper function to create a new error of type TypeUnauthorized, with formatted message
func Unauthorizedf(format string, args ...interface{}) error {
	return newErrf(nil, TypeUnauthorized, format, args...)
}

// Empty is a helper function to create a new error of type TypeEmpty
func Empty(message string) error {
	return newErr(nil, message, TypeEmpty)
}

// Emptyf is a helper function to create a new error of type TypeEmpty, with formatted message
func Emptyf(format string, args ...interface{}) error {
	return newErrf(nil, TypeEmpty, format, args...)
}

// NotFound is a helper function to create a new error of type TypeNotFound
func NotFound(message string) error {
	return newErr(nil, message, TypeNotFound)
}

// NotFoundf is a helper function to create a new error of type TypeNotFound, with formatted message
func NotFoundf(format string, args ...interface{}) error {
	return newErrf(nil, TypeNotFound, format, args...)
}

// LimitExceeded is a helper function to create a new error of type TypeLimitExceeded
func LimitExceeded(message string) error {
	return newErr(nil, message, TypeLimitExceeded)
}

// LimitExceededf is a helper function to create a new error of type TypeLimitExceeded, with formatted message
func LimitExceededf(format string, args ...interface{}) error {
	return newErrf(nil, TypeLimitExceeded, format, args...)
}

// SubscriptionExpired is a helper function to create a new error of type TypeSubscriptionExpired
func SubscriptionExpired(message string) error {
	return newErr(nil, message, TypeSubscriptionExpired)
}

// SubscriptionExpiredf is a helper function to create a new error of type TypeSubscriptionExpired, with formatted message
func SubscriptionExpiredf(format string, args ...interface{}) error {
	return newErrf(nil, TypeSubscriptionExpired, format, args...)
}

// HasType will check if the provided err type is available anywhere nested in the error
func HasType(err error, et Typer) bool {
	if err == nil {
		return false
	}
	e, _ := err.(interface {
		Type() Typer
	})
	if e == nil {
		return HasType(errors.Unwrap(err), et)
	}
	if e.Type() == et {
		return true
	}
	return HasType(errors.Unwrap(err), et)
}

// GetAPIError tries to get the code and message from any error.
func GetAPIError(err error) (code int, msg string) {
	if err == nil {
		return defaultErrType.HTTPStatusCode(), DefaultMessage
	}
	msg = err.Error()
	for err != nil {
		if apiErr, _ := err.(APIError); apiErr != nil {
			code, _ = apiErr.APIError()
			return code, msg
		}
		err = Unwrap(err)
	}
	return defaultErrType.HTTPStatusCode(), msg
}
