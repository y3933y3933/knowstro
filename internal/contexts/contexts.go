package contexts

import (
	"context"
	"net/http"

	"github.com/y3933y3933/knowstro/internal/store"
)

type contextKey = string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
