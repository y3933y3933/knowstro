package response

import "net/http"

type Response[T any] struct {
	Success bool      `json:"success"`
	Data    *T        `json:"data,omitzero"`
	Error   *APIError `json:"error,omitzero"`
}

type APIError struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitzero"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewSuccess[T any](data T) (int, Response[T]) {
	return http.StatusOK, Response[T]{
		Success: true,
		Data:    &data,
	}
}

func NewError(code int, msg string, details ...FieldError) (int, Response[any]) {
	return code, Response[any]{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: msg,
			Details: details,
		},
	}
}
