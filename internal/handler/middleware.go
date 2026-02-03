package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/yourusername/resume-builder/internal/service"
)

type contextKey string

const userIDKey contextKey = "userID"

type AuthMiddleware struct {
	jwtService *service.JWTService
	logger     *slog.Logger
}

func NewAuthMiddleware(jwtService *service.JWTService, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn("missing authorization header", "path", r.URL.Path)
			respondError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.logger.Warn("invalid authorization header format", "path", r.URL.Path)
			respondError(w, http.StatusUnauthorized, "invalid authorization header")
			return
		}

		userID, err := m.jwtService.ValidateToken(parts[1])
		if err != nil {
			m.logger.Warn("invalid token", "error", err, "path", r.URL.Path)
			respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		m.logger.Info("authenticated request", "user_id", userID, "path", r.URL.Path)
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

func respondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondJSON(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
