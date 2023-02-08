package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/render"

	"github.com/sainak/bitsb/domain"
	"github.com/sainak/bitsb/utils/jwt"
)

var UserCtxKey = &ContextKey{"user"}

// JWTAuth is a middleware that checks for a valid JWT in the Authorization header.
// If one is found, it will be parsed and the user will be added to the request context.
func JWTAuth(j *jwt.JWT, u domain.UserStorer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the JWT string from the auth header
			authHeader := r.Header.Get("Authorization")
			bearerToken := strings.Split(authHeader, " ")
			if authHeader == "" || len(bearerToken) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, render.M{"message": "auth token not provided"})
				return
			}

			id, err := j.GetUserID(bearerToken[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, render.M{"message": err.Error()})
				return
			}

			user, err := u.SelectByID(r.Context(), id)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				render.JSON(w, r, render.M{"message": err.Error()})
				return
			}

			ctx := context.WithValue(r.Context(), UserCtxKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
