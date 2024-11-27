//go:build !solution

package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

var ErrInvalidToken = errors.New("invalid token")

type User struct {
	Name  string
	Email string
}

func ContextUser(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserKey).(*User)
	return user, ok
}

type TokenChecker interface {
	CheckToken(ctx context.Context, token string) (*User, error)
}

func CheckAuth(checker TokenChecker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token = token[7:]
			user, err := checker.CheckToken(r.Context(), token)
			if err != nil {
				if errors.Is(err, ErrInvalidToken) {
					w.WriteHeader(http.StatusUnauthorized)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

var UserKey = &contextKey{"User"}

type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return "auth context key " + c.name
}
