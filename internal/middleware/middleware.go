package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/y3933y3933/knowstro/internal/contexts"
	"github.com/y3933y3933/knowstro/internal/response"
	"github.com/y3933y3933/knowstro/internal/store"
	"github.com/y3933y3933/knowstro/internal/tokens"
)

type UserMiddleware struct {
	UserStore store.UserStore
}

func (um *UserMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Vary", "Authorization")

		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.Request = contexts.SetUser(c.Request, store.AnonymousUser)
			c.Next()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.InvalidAuthenticationToken(c)
			return
		}

		token := headerParts[1]

		user, err := um.UserStore.GetForToken(tokens.ScopeAuth, token)
		if err != nil {
			response.InvalidAuthenticationToken(c)
			return
		}

		c.Request = contexts.SetUser(c.Request, user)
		c.Next()
	}
}
