package errors

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

// APIError returns an HTTP status code and an API-safe error message.
type APIError interface {
	APIError() (statusCode int, msg string)
}

type I18ner interface {
	LocalizeConfig() (lc *i18n.LocalizeConfig)
}

type Typer interface {
	HTTPStatusCode() int
}

type CustomType struct {
	Detail     string
	StatusCode int
}

func NewCustomType(detail string, statusCode int) CustomType {
	return CustomType{
		Detail:     detail,
		StatusCode: statusCode,
	}
}

func (c CustomType) HTTPStatusCode() int {
	return c.StatusCode
}

type errType int

// While adding a new Type, the respective helper functions should be added, also update the
// WriteHTTP method accordingly
const (
	// TypeInternal is error type for when there is an internal system error. e.g. Database errors
	TypeInternal errType = iota
	// TypeValidation is error type for when there is a validation error. e.g. invalid email address
	TypeValidation
	// TypeInput is error type for when an input data type error. e.g. invalid JSON
	TypeInput
	// TypeDuplicate is error type for when there's duplicate content
	TypeDuplicate
	// TypeUnauthenticated is error type when trying to access an authenticated API without authentication
	TypeUnauthenticated
	// TypeNoPermission is error type for when there's an unauthorized access attempt
	TypeNoPermission
	// TypeEmpty is error type for when an expected non-empty resource, is empty
	TypeEmpty
	// TypeNotFound is error type for an expected resource is not found e.g. user ID not found
	TypeNotFound
	// TypeLimitExceeded is error type for attempting the same action more than allowed
	TypeLimitExceeded
	// TypeSubscriptionExpired is error type for when a user's 'paid' account has expired
	TypeSubscriptionExpired

	// DefaultMessage is the default user friendly message
	DefaultMessage = "unknown error occurred"
)

var (
	defaultErrType = TypeInternal
)

// SetDefaultType will set the default error type, which is used in the 'New' function
func SetDefaultType(e errType) {
	defaultErrType = e
}

// HTTPStatusCode is a convenience method used to get the appropriate HTTP response status code for the respective error type
func (e errType) HTTPStatusCode() int {
	status := http.StatusInternalServerError
	switch e {
	case TypeValidation:
		{
			status = http.StatusUnprocessableEntity
		}
	case TypeInput:
		{
			status = http.StatusBadRequest
		}

	case TypeDuplicate:
		{
			status = http.StatusConflict
		}

	case TypeUnauthenticated:
		{
			status = http.StatusUnauthorized
		}
	case TypeNoPermission:
		{
			status = http.StatusForbidden
		}

	case TypeEmpty:
		{
			status = http.StatusGone
		}

	case TypeNotFound:
		{
			status = http.StatusNotFound

		}
	case TypeLimitExceeded:
		{
			status = http.StatusTooManyRequests
		}
	case TypeSubscriptionExpired:
		{
			status = http.StatusPaymentRequired
		}
	}

	return status
}
