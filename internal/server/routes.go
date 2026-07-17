package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"open-fermentations/internal/route"

	"fmt"
	"time"

	"github.com/coder/websocket"
)

const (
	ContentTypeJSON = "application/json"
)

const ServerPrefix = "/api"

func registerRoute(mux *http.ServeMux, h *route.Route) {
	slog.Info("Registering route", slog.String("method", h.Method), slog.String("route", h.Route))
	mux.Handle(fmt.Sprintf("%v %v", h.Method, h.Route), h.Handler)
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	defaultMiddleware := []route.Middleware{
		s.corsMiddleware,
	}

	// Register routes
	authenticatedRoutesWithPrefix := []*route.Route{
		route.New(http.MethodGet, "/logout", http.HandlerFunc(s.logoutHandler)),
	}

	routeHandlers := []*route.Route{
		route.New(http.MethodPost, "/login", http.HandlerFunc(s.loginHandler)).
			WithPrefix(ServerPrefix),
		route.New(http.MethodGet, "/health", http.HandlerFunc(s.healthHandler)),
	}

	for _, r := range authenticatedRoutesWithPrefix {
		routeHandlers = append(routeHandlers, r.WithMiddleware(s.authenticationMiddleware).
			WithPrefix(ServerPrefix))
	}

	for _, r := range routeHandlers {
		registerRoute(mux, r)
	}

	slog.Info("Registering route", slog.String("route", "/websocket"))
	mux.Handle("/websocket", http.HandlerFunc(s.websocketHandler))

	var handler http.Handler = mux
	for _, m := range defaultMiddleware {
		handler = m(handler)
	}
	return handler
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
