package veilmail

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// ErrorType classifies the kind of API error.
type ErrorType string

const (
	// ErrorTypeAuthentication indicates an invalid or missing API key (401).
	ErrorTypeAuthentication ErrorType = "authentication"
	// ErrorTypeForbidden indicates insufficient permissions (403).
	ErrorTypeForbidden ErrorType = "forbidden"
	// ErrorTypeNotFound indicates the requested resource was not found (404).
	ErrorTypeNotFound ErrorType = "not_found"
	// ErrorTypeValidation indicates a request validation failure (400).
	ErrorTypeValidation ErrorType = "validation"
	// ErrorTypePiiDetected indicates PII was detected in email content (422).
	ErrorTypePiiDetected ErrorType = "pii_detected"
	// ErrorTypeRateLimit indicates the rate limit has been exceeded (429).
	ErrorTypeRateLimit ErrorType = "rate_limit"
	// ErrorTypeServer indicates an internal server error (5xx).
	ErrorTypeServer ErrorType = "server"
	// ErrorTypeUnknown indicates an unrecognized error.
	ErrorTypeUnknown ErrorType = "unknown"
)

// Error represents an API error from the Veil Mail API.
type Error struct {
	// Type classifies the error.
	Type ErrorType

	// Code is the machine-readable error code from the API.
	Code string

	// Message is the human-readable error message.
	Message string

	// StatusCode is the HTTP status code.
	StatusCode int

	// Details contains additional error context.
	Details map[string]interface{}

	// RetryAfter is the number of seconds to wait before retrying (rate limit errors only).
	RetryAfter int

	// PiiTypes lists the types of PII detected (PII detection errors only).
	PiiTypes []string
}

func (e *Error) Error() string {
	return fmt.Sprintf("veilmail: %s (code: %s, status: %d)", e.Message, e.Code, e.StatusCode)
}

// IsAuthenticationError reports whether the error is an authentication error (401).
func IsAuthenticationError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypeAuthentication
}

// IsForbiddenError reports whether the error is a forbidden error (403).
func IsForbiddenError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypeForbidden
}

// IsNotFoundError reports whether the error is a not found error (404).
func IsNotFoundError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypeNotFound
}

// IsValidationError reports whether the error is a validation error (400).
func IsValidationError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypeValidation
}

// IsPiiDetectedError reports whether the error is a PII detection error (422).
func IsPiiDetectedError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypePiiDetected
}

// IsRateLimitError reports whether the error is a rate limit error (429).
func IsRateLimitError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypeRateLimit
}

// IsServerError reports whether the error is a server error (5xx).
func IsServerError(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Type == ErrorTypeServer
}

type apiErrorResponse struct {
	Error struct {
		Code    string                 `json:"code"`
		Message string                 `json:"message"`
		Details map[string]interface{} `json:"details,omitempty"`
	} `json:"error"`
}

func parseAPIError(statusCode int, header http.Header, body []byte) *Error {
	e := &Error{
		StatusCode: statusCode,
		Code:       "unknown",
		Message:    "An unknown error occurred",
	}

	var apiErr apiErrorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil {
		if apiErr.Error.Code != "" {
			e.Code = apiErr.Error.Code
		}
		if apiErr.Error.Message != "" {
			e.Message = apiErr.Error.Message
		}
		e.Details = apiErr.Error.Details
	}

	switch {
	case statusCode == 401:
		e.Type = ErrorTypeAuthentication
	case statusCode == 403:
		e.Type = ErrorTypeForbidden
	case statusCode == 404:
		e.Type = ErrorTypeNotFound
	case statusCode == 422:
		e.Type = ErrorTypePiiDetected
		if e.Details != nil {
			if types, ok := e.Details["piiTypes"]; ok {
				if arr, ok := types.([]interface{}); ok {
					for _, v := range arr {
						if s, ok := v.(string); ok {
							e.PiiTypes = append(e.PiiTypes, s)
						}
					}
				}
			}
		}
	case statusCode == 429:
		e.Type = ErrorTypeRateLimit
		if ra := header.Get("Retry-After"); ra != "" {
			if seconds, err := strconv.Atoi(ra); err == nil {
				e.RetryAfter = seconds
			}
		}
	case statusCode >= 400 && statusCode < 500:
		e.Type = ErrorTypeValidation
	case statusCode >= 500:
		e.Type = ErrorTypeServer
	default:
		e.Type = ErrorTypeUnknown
	}

	return e
}
