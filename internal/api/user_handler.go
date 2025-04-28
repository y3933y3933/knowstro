package api

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/mailer"
	"github.com/y3933y3933/knowstro/internal/response"
	"github.com/y3933y3933/knowstro/internal/store"
	"github.com/y3933y3933/knowstro/internal/tokens"
	"github.com/y3933y3933/knowstro/internal/utils"
)

type registerUserRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,max=15,min=8"`
}

type UserHandler struct {
	userStore  store.UserStore
	tokenStore store.TokenStore
	logger     *slog.Logger
	mailer     *mailer.Mailer
}

func NewUserHandler(userStore store.UserStore, tokenStore store.TokenStore, logger *slog.Logger, mailer *mailer.Mailer) *UserHandler {
	return &UserHandler{
		userStore:  userStore,
		tokenStore: tokenStore,
		logger:     logger,
		mailer:     mailer,
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
			response.FailedValidationError(c, []response.FieldError{{Field: "email", Message: "duplicate email"}})
		case errors.Is(err, store.ErrDuplicateUserName):
			response.FailedValidationError(c, []response.FieldError{{Field: "name", Message: "duplicate username"}})

		default:
			response.InternalError(c)
		}

		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 3*24*time.Hour, tokens.ScopeActivation)
	if err != nil {
		h.logger.Error("create new token: %v", err)
		response.InternalError(c)
		return
	}
	go func() {
		data := struct {
			AppName       string
			UserName      string
			ActivationURL string
			Token         string
		}{
			AppName:       "Knowstro",
			UserName:      user.Name,
			ActivationURL: "test",
			Token:         token.Plaintext,
		}

		err = h.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			h.logger.Error(err.Error())

		}
	}()
	response.SuccessCreated(c, user)
}

func (h *UserHandler) HandlerActivateUser(c *gin.Context) {
	var input struct {
		TokenPlaintext string `json:"token" binding:"required"`
	}

	err := utils.ReadJSON(c, input)
	if err != nil {
		details, isValid := utils.ValidationErrors(err)
		if !isValid {
			response.FailedValidationError(c, details)
		} else {
			response.BadRequest(c, err.Error())
		}
		return
	}

	user, err := h.userStore.GetForToken(tokens.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		h.logger.Error("get for token: %v", err)
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.FailedValidationError(c, []response.FieldError{{
				Field:   "token",
				Message: "invalid or expired activation token",
			}})
		default:
			response.InternalError(c)
		}
		return
	}

	user.Activated = true

	err = h.userStore.UpdateUser(user)
	if err != nil {
		h.logger.Error("update user: %v", err)
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			response.FailedValidationError(c, []response.FieldError{{Field: "email", Message: "duplicate email"}})
		case errors.Is(err, store.ErrDuplicateUserName):
			response.FailedValidationError(c, []response.FieldError{{Field: "name", Message: "duplicate username"}})

		default:
			response.InternalError(c)
		}
	}

	err = h.tokenStore.DeleteAllTokensForUser(user.ID, tokens.ScopeActivation)
	if err != nil {
		h.logger.Error("delete all tokens for user: %v", err)
		response.InternalError(c)
		return
	}

	response.SuccessOK(c, user)

}
