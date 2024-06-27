package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hwaengfan/dev-journal-backend/internal/services/user"
	"github.com/hwaengfan/dev-journal-backend/internal/database/repositories/user"
)

type Server struct {
	address string
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

	// Start server
	log.Println("Starting HTTP server on address", server.address)
	return http.ListenAndServe(server.address, router)
}
