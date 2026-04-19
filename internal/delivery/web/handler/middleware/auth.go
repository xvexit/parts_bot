package middleware

import (
	"context"
	"net/http"
	"partsBot/internal/usecase/auth"
	"strings"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func Auth(tm auth.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(header, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			userID, err := tm.ParseAccessToken(parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}
