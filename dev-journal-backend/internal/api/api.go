package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	projectRepository "github.com/hwaengfan/dev-journal-backend/internal/database/repositories/project"
	userRepository "github.com/hwaengfan/dev-journal-backend/internal/database/repositories/user"
	projectService "github.com/hwaengfan/dev-journal-backend/internal/services/project"
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

	// Set up user routes
	userStore := userRepository.NewStore(server.database)
	userHandler := userService.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	// Set up project routes
	projectStore := projectRepository.NewStore(server.database)
	projectHandler := projectService.NewHandler(projectStore, userStore)
	projectHandler.RegisterRoutes(subrouter)

	// Start server
	log.Println("Starting HTTP server on address", server.address)
	return http.ListenAndServe(server.address, router)
}
