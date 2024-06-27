package database

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func InitializeMySQLStorage(config mysql.Config) (*sql.DB, error) {
	database, error := sql.Open("mysql", config.FormatDSN())
	if error != nil {
		log.Fatalf("Error occured while opening MySQL database: %v", error)
		return database, error
	}

	// Check connections
	if error := database.Ping(); error != nil {
		log.Fatalf("Error occured while pinging MySQL database: %v", error)
		return database, error
	}
	
	log.Println("Successfully connected to MySQL database")
	return database, nil
}
