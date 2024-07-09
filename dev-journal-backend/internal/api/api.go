package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	projectRepository "github.com/hwaengfan/dev-journal-backend/internal/database/repositories/project"
	taskRepository "github.com/hwaengfan/dev-journal-backend/internal/database/repositories/task"
	userRepository "github.com/hwaengfan/dev-journal-backend/internal/database/repositories/user"
	projectService "github.com/hwaengfan/dev-journal-backend/internal/services/project"
	taskService "github.com/hwaengfan/dev-journal-backend/internal/services/task"
	userService "github.com/hwaengfan/dev-journal-backend/internal/services/user"
)

type Server struct {
	address  string
	database *sql.DB
}

func NewServer(address string, database *sql.DB) *Server {
	return &Server{address: address, database: database}
}

func (server *Server) Run() error {
	// Set up router
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Set up stores
	userStore := userRepository.NewStore(server.database)
	projectStore := projectRepository.NewStore(server.database)
	taskStore := taskRepository.NewStore(server.database)

	// Set up user routes
	userHandler := userService.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Set up project routes
	projectHandler := projectService.NewHandler(projectStore, userStore, taskStore)
	projectHandler.RegisterRoutes(subrouter)

	// Set up task routes
	taskHandler := taskService.NewHandler(taskStore, userStore, projectStore)
	taskHandler.RegisterRoutes(subrouter)

	// Start server
	log.Println("Starting HTTP server on address", server.address)
	return http.ListenAndServe(server.address, router)
}
