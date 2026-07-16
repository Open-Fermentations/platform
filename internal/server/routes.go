package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"fmt"
	"time"

	"github.com/coder/websocket"
)

const (
	ContentTypeJSON = "application/json"
)

const ServerPrefix = "/api"

func registerServerPrefixedRoute(mux *http.ServeMux, method, route string, handler http.HandlerFunc) {
	registerRoute(mux, method, fmt.Sprintf("%v%v", ServerPrefix, route), handler)
}

func registerRoute(mux *http.ServeMux, method, route string, handler http.HandlerFunc) {
	slog.Info("Registering route", slog.String("method", method), slog.String("route", route))
	mux.HandleFunc(fmt.Sprintf("%v %v", method, route), handler)
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	registerServerPrefixedRoute(mux, http.MethodPost, "/login", s.loginHandler)
	registerServerPrefixedRoute(mux, http.MethodGet, "/logout", s.logoutHandler)

	registerRoute(mux, http.MethodGet, "/health", s.healthHandler)

	slog.Info("Registering route", slog.String("route", "/websocket"))
	mux.HandleFunc("/websocket", s.websocketHandler)

	// Wrap the mux with CORS middleware
	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Failed to marshal health check response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		slog.Error("failed to write response", slog.String("error", err.Error()))
	}
}

func (s *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to open websocket", http.StatusInternalServerError)
		return
	}
	defer socket.Close(websocket.StatusGoingAway, "Server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		if err := socket.Write(socketCtx, websocket.MessageText, []byte(payload)); err != nil {
			slog.Error("failed to write to socket", slog.String("error", err.Error()))
			break
		}
		time.Sleep(2 * time.Second)
	}
}
