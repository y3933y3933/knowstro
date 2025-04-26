package api

import (
	"errors"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/response"
	"github.com/y3933y3933/knowstro/internal/store"
	"github.com/y3933y3933/knowstro/internal/utils"
)

type registerUserRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=15,min=8"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *slog.Logger
}

func NewUserHandler(userStore store.UserStore, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) HandleRegisterUser(c *gin.Context) {
	var req registerUserRequest

	err := utils.ReadJSON(c, &req)
	if err != nil {
		h.logger.Error("decoding register request: %v", err)
		if details, isValid := utils.ValidationErrors(err); !isValid {
			response.FailedValidationError(c, details)
		} else {
			response.BadRequest(c, err.Error())

		}
		return
	}

	user := &store.User{
		Name:      req.Name,
		Email:     req.Email,
		Activated: false,
	}

	err = user.Password.Set(req.Password)
	if err != nil {
		h.logger.Error("hashing password: ", err.Error())
		response.InternalError(c)
		return
	}

	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Error("register user: ", err.Error())
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			response.FailedValidationError(c, []response.FieldError{{Field: "Email", Message: "duplicate email"}})
		case errors.Is(err, store.ErrDuplicateUserName):
			response.FailedValidationError(c, []response.FieldError{{Field: "Name", Message: "duplicate username"}})

		default:
			response.InternalError(c)
		}

		return
	}

	// TODO: http.StatusCreated
	response.Success(c, user)
}
