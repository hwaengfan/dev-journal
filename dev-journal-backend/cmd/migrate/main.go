package main

import (
	"log"
	"os"

	mysqlConfig "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hwaengfan/dev-journal-backend/configs"
	"github.com/hwaengfan/dev-journal-backend/internal/database"
)

func main() {
	// Connect to database
	mysqlStorage, error := database.InitializeMySQLStorage(mysqlConfig.Config{
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

	// Create MySQL driver for migration
	driver, error := mysql.WithInstance(mysqlStorage, &mysql.Config{})
	if error != nil {
		log.Fatalf("Error occured while creating MySQL driver for migration: %v", error)
	}

	// Create migrator
	migrator, error := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)
	if error != nil {
		log.Fatalf("Error occured while creating migrator: %v", error)
	}

	// Get current version
	version, dirty, _ := migrator.Version()
	log.Printf("Version: %d, dirty: %v", version, dirty)

	// Run migration
	command := os.Args[(len(os.Args) - 1)]
	if command == "up" {
		if error := migrator.Up(); error != nil && error != migrate.ErrNoChange {
			log.Fatalf("Error occured while migrating up: %v", error)
		}
	}

	if command == "down" {
		if error := migrator.Down(); error != nil && error != migrate.ErrNoChange {
			log.Fatalf("Error occured while migrating down: %v", error)
		}
	}
}
