package middleware

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/sainak/bitsb/domain"
)

// AccessAbove is a middleware that checks if the user has access level above the given level
func AccessAbove(level domain.AccessLevel) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get user form context
			user := r.Context().Value(UserCtxKey).(*domain.User)
			if user == nil {
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, render.M{"message": "user not found in context"})
				return
			}
			if user.Access < level {
				w.WriteHeader(http.StatusForbidden)
				render.JSON(w, r, render.M{"message": "you don't have permission to perform this action"})
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
