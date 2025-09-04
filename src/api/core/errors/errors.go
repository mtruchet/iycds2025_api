package errors

import "fmt"

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.Code, e.Message)
}

func NewBadRequest(message string) *APIError {
	return &APIError{
		Code:    400,
		Message: message,
	}
}

func NewUnauthorized(message string) *APIError {
	return &APIError{
		Code:    401,
		Message: message,
	}
}

func NewNotFound(message string) *APIError {
	return &APIError{
		Code:    404,
		Message: message,
	}
}

func NewInternalServerError(message string) *APIError {
	return &APIError{
		Code:    500,
		Message: message,
	}
}
