package response

type Response[T any] struct {
	Success bool      `json:"success"`
	Data    *T        `json:"data,omitzero"`
	Error   *APIError `json:"error,omitzero"`
}

type APIError struct {
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitzero"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewSuccess[T any](status int, data T) (int, Response[T]) {
	return status, Response[T]{
		Success: true,
		Data:    &data,
	}
}

func NewError(code int, msg string, details ...FieldError) (int, Response[any]) {
	return code, Response[any]{
		Success: false,
		Error: &APIError{
			Message: msg,
			Details: details,
		},
	}
}
