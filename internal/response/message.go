package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	MsgRecordNotFound             = "record not found"
	MsgInternalServerError        = "internal server error"
	MsgFailedValidation           = "validation fail"
	MsgInvalidCredentials         = "invalid authentication credentials"
	MsgInvalidAuthenticationToken = "invalid or missing authentication token"
)

func SuccessOK(c *gin.Context, data any) {
	status, res := NewSuccess(http.StatusOK, data)
	c.JSON(status, res)
}

func SuccessCreated(c *gin.Context, data any) {
	status, res := NewSuccess(http.StatusCreated, data)
	c.JSON(status, res)
}

func RecordNotFound(c *gin.Context) {
	status, res := NewError(http.StatusNotFound, MsgRecordNotFound)
	c.AbortWithStatusJSON(status, res)
}

func InternalError(c *gin.Context) {
	status, res := NewError(http.StatusInternalServerError, MsgInternalServerError)
	c.AbortWithStatusJSON(status, res)
}

func BadRequest(c *gin.Context, msg string) {
	status, res := NewError(http.StatusBadRequest, msg)
	c.AbortWithStatusJSON(status, res)
}

func UnprocessableError(c *gin.Context, msg string) {
	status, res := NewError(http.StatusUnprocessableEntity, msg)
	c.AbortWithStatusJSON(status, res)
}

func FailedValidationError(c *gin.Context, details []FieldError) {
	status, res := NewError(http.StatusUnprocessableEntity, MsgFailedValidation, details...)
	c.AbortWithStatusJSON(status, res)
}

func InvalidCredential(c *gin.Context) {
	status, res := NewError(http.StatusUnauthorized, MsgInvalidCredentials)
	c.AbortWithStatusJSON(status, res)
}

func InvalidAuthenticationToken(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer")
	status, res := NewError(http.StatusUnauthorized, MsgInvalidAuthenticationToken)
	c.AbortWithStatusJSON(status, res)
}
