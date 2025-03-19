package server

import (
	"context"
	"net/http"
	"time"

	"github.com/mehmetali10/task-planner/internal/task/config"
	"github.com/mehmetali10/task-planner/internal/task/handler"

	"github.com/mehmetali10/task-planner/pkg/log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	handler    handler.Handler
	logger     log.Logger
}

func NewServer(handler handler.Handler) *Server {
	router := mux.NewRouter()

	return &Server{
		router:  router,
		handler: handler,
		logger:  log.NewLogger("server", config.GetApp().HTTPServerLogLevel),
	}

}

func (s *Server) Start(addr string) {
	s.setUpRoutes()

	s.httpServer = &http.Server{
		Addr: addr,
		Handler: handlers.CORS(
			handlers.AllowedOrigins(config.GetApp().HTTPAllowedOrigins),
			handlers.AllowedMethods(config.GetApp().HTTPAllowedMethods),
			handlers.AllowedHeaders(config.GetApp().HTTPAllowedHeaders),
		)(s.router),
	}

	go func() {
		s.logger.Info("Server is starting on addr=%s", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("Server failed to start: error=%v", err)
		}
	}()
}

func (s *Server) Stop() {
	s.logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Fatal("Server forced to shutdown: error=%v", err)
	}

	s.logger.Info("Server has been stopped")
}

func (s *Server) setUpRoutes() {
	s.router.HandleFunc("/task", s.handler.CreateTask()).Methods(http.MethodPost)
	s.router.HandleFunc("/tasks", s.handler.ListTasks()).Methods(http.MethodGet)

	s.router.HandleFunc("/assignment", s.handler.ScheduleAssaignments()).Methods(http.MethodGet)

	s.router.HandleFunc("/developers", s.handler.ListDevelopers()).Methods(http.MethodGet)
}
