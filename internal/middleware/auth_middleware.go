package middleware

import (
	"context"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/domain/token"
	"net/http"
	"strings"
)

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value("user").(*domain.User)
	return user, ok
}

func AuthMiddleware(tokenGen token.Generator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "not authorized", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims, err := tokenGen.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			user := &domain.User{
				Id:   claims.UserId,
				Role: claims.Role,
			}

			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
