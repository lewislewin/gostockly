package middleware

import (
	"context"
	"net/http"
	"strings"

	"gostockly/internal/services"
	"gostockly/pkg/logger"
	"gostockly/pkg/utils"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "user_id"
const companyIDKey contextKey = "company_id"

func AuthMiddleware(userService *services.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.GetLogger()

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Error("Missing Authorization header")
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Authorization header missing")
				return
			}

			// Parse the token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				log.Error("Invalid Authorization header format")
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid Authorization header format")
				return
			}
			token := parts[1]

			// Validate the JWT token to extract the user ID
			userID, err := utils.ValidateJWT(token, userService.JWTSecret)
			if err != nil {
				log.Error("Invalid or expired token: %v", err)
				utils.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Fetch the user and associated company from the database
			user, err := userService.UserRepo.GetUserByID(userID)
			if err != nil {
				log.Error("User not found: %v", userID)
				utils.WriteErrorResponse(w, http.StatusForbidden, "User not found")
				return
			}

			if user.CompanyID == uuid.Nil {
				log.Error("No company associated with user ID: %v", userID)
				utils.WriteErrorResponse(w, http.StatusForbidden, "No company associated with user")
				return
			}

			// Set user_id and company_id in the context
			ctx := context.WithValue(r.Context(), "user_id", userID)
			ctx = context.WithValue(ctx, "company_id", user.CompanyID.String())
			log.Info("Authenticated user ID: %s, Company ID: %s", userID, user.CompanyID.String())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
