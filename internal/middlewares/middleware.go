package middlewares

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/vinayb91/task-manager-ai-backend.git/internal/services"
)

func AuthMiddleware(authService *services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := os.Getenv("ENV")

		allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
		if allowedOrigin == "" || env == "dev" {
			allowedOrigin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", strings.TrimRight(allowedOrigin, "/"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[HTTP] START %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("[HTTP] END %s %s completed in %v", r.Method, r.RequestURI, time.Since(start))
	})
}
