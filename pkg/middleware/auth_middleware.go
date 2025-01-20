package middleware

import (
	"context"
	"net/http"
	"strings"

	"gostockly/pkg/logger"
	"gostockly/pkg/utils"
)

type contextKey struct{}

var userIDKey = contextKey{}

// AuthMiddleware validates the JWT token from the Authorization header.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetLogger()

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("Missing Authorization header")
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authorization header missing")
				return
			}

			// Extract token from the Authorization header
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Error("Invalid Authorization header format")
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}
			token := parts[1]

			// Validate the token
			userID, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				log.Error("Invalid or expired token: %s", err.Error())
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Log successful authentication
			log.Info("Authenticated user ID: %s", userID)

			// Attach the user ID to the request context
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts the user ID from the request context.
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
