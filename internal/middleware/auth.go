package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"service/internal/auth"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		token, err := auth.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "invalid user id", http.StatusUnauthorized)
			return
		}

		role, _ := claims["role"].(string)

		ctx := context.WithValue(r.Context(), UserIDKey, int64(userID))
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(RoleKey).(string)

		if !ok || role != "admin" {
			http.Error(w, "admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
