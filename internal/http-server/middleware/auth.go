package middleware

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthRequired() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, ok := r.Context().Value(userCtxKey).(jwt.MapClaims)
			if !ok || claims == nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			if exp, ok := claims["exp"].(float64); ok {
				if int64(exp) < time.Now().Unix() {
					http.Error(w, "token expired", http.StatusUnauthorized)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
