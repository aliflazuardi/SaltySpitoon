package server

import (
	"SaltySpitoon/internal/constants"
	"context"
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("POST /v1/register", s.registerHandler)
	mux.HandleFunc("POST /v1/login", s.loginHandler)
	mux.HandleFunc("POST /v1/activity", s.createActivityHandler)
	mux.HandleFunc("DELETE /v1/activity/", s.deleteActivityHandler)

	return s.authMiddleware(mux)
}

//func (s *Server) corsMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Set CORS headers
//		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
//		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
//		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required
//
//		// Handle preflight OPTIONS requests
//		if r.Method == http.MethodOptions {
//			w.WriteHeader(http.StatusNoContent)
//			return
//		}
//
//		// Proceed with the next handler
//		next.ServeHTTP(w, r)
//	})
//}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/health" || path == "/v1/login" || path == "/v1/register" {
			next.ServeHTTP(w, r)
			return
		}

		userID := 1234

		ctx := context.WithValue(r.Context(), constants.UserIDCtxKey, int64(userID))

		// Proceed with the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
