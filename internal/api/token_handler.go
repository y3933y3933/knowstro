package api

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/response"
	"github.com/y3933y3933/knowstro/internal/store"
	"github.com/y3933y3933/knowstro/internal/tokens"
	"github.com/y3933y3933/knowstro/internal/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *slog.Logger
}

type createTokenRequest struct {
	Name     string `json:"name" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=15,min=8"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *slog.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(c *gin.Context) {
	var req createTokenRequest

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

	user, err := h.userStore.GetUserByName(req.Name)
	if err != nil {
		h.logger.Error("GetUserByName: %v", err)

		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			response.InvalidCredential(c)
		default:
			response.InternalError(c)
		}
		return

	}

	passwordsDoMatch, err := user.Password.Matches(req.Password)
	if err != nil {
		h.logger.Error("PasswordHash.Matches: %v", err)
		response.InternalError(c)
		return
	}

	if !passwordsDoMatch {
		response.InvalidCredential(c)
		return
	}

	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeActivation)
	if err != nil {
		h.logger.Error("Creating Token: %v", err)
		response.InternalError(c)
		return
	}

	response.SuccessCreated(c, token)
}
