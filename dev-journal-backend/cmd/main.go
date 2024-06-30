package main

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/hwaengfan/dev-journal-backend/configs"
	"github.com/hwaengfan/dev-journal-backend/internal/api"
	"github.com/hwaengfan/dev-journal-backend/internal/database"
)

func main() {
	// Connect to database
	mysqlStorage, error := database.InitializeMySQLStorage(mysql.Config{
		User: configs.DatabaseEnvironmentVariables.User,
		Passwd: configs.DatabaseEnvironmentVariables.Password,
		Addr: configs.DatabaseEnvironmentVariables.Address,
		DBName: configs.DatabaseEnvironmentVariables.Name,
		Net: "tcp",
		AllowNativePasswords: true,
		ParseTime: true,
	})
	if error != nil {
		log.Fatalf("Error occured while connecting to MySQL database: %v", error)
	}
	
	defer mysqlStorage.Close()

	// Setting up server
	address := fmt.Sprintf(":%s", configs.ServerEnvironmentVariables.Port)
	server := api.NewServer(address, mysqlStorage)
	if error := server.Run(); error != nil {
		log.Fatalf("Error occured while running HTTP server: %v", error)
	}
}
