package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/y3933y3933/knowstro/internal/response"
)

// TODO: validator field error handling
func ReadJSON(c *gin.Context, dst any) error {
	err := c.ShouldBindJSON(dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	return nil
}

func ReadIDParam(c *gin.Context) (int64, error) {
	s := c.Param("id")
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil || i < 1 {
		return 0, errors.New("Invalid ID parameter")
	}

	return i, nil
}

func ValidateJSON(err error) validator.ValidationErrors {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return ve
	}
	return nil
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s is not valid", fe.Field())

	}
}

func ValidationErrors(err error) (details []response.FieldError, isValid bool) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details = make([]response.FieldError, len(ve))
		for i, fe := range ve {
			details[i] = response.FieldError{
				Field:   fe.Field(),
				Message: msgForTag(fe),
			}
		}
		return details, false
	}
	return nil, true
}
