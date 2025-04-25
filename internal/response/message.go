package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MsgRecordNotFound      = "record not found"
	MsgInternalServerError = "internal server error"
	MsgFailedValidation    = "validation fail"
)

func Success(c *gin.Context, data any) {
	status, res := NewSuccess(data)
	c.JSON(status, res)
}

func RecordNotFound(c *gin.Context) {
	status, res := NewError(http.StatusNotFound, MsgRecordNotFound)
	c.JSON(status, res)
}

func InternalError(c *gin.Context) {
	status, res := NewError(http.StatusInternalServerError, MsgInternalServerError)
	c.JSON(status, res)
}

func BadRequest(c *gin.Context, msg string) {
	status, res := NewError(http.StatusBadRequest, msg)
	c.JSON(status, res)
}

func UnprocessableError(c *gin.Context, msg string) {
	status, res := NewError(http.StatusUnprocessableEntity, msg)
	c.JSON(status, res)
}

func FailedValidationError(c *gin.Context, details []FieldError) {
	status, res := NewError(http.StatusUnprocessableEntity, MsgFailedValidation, details...)
	c.JSON(status, res)
}
