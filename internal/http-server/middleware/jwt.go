package middleware

import (
	"context"
	"net/http"
	"service/internal/config"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userCtxKey = contextKey("user")

func JWTAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}
			cfg := config.MustLoad()

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, http.ErrAbortHandler
				}
				return []byte(cfg.JwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userCtxKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserClaims(r *http.Request) jwt.MapClaims {
	claims, _ := r.Context().Value(userCtxKey).(jwt.MapClaims)
	return claims
}
