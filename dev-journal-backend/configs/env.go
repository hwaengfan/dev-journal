package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfigs struct {
	User string
	Password string
	Address string
	Name string
}

type ServerConfigs struct {
	PublicHost string
	Port string
}

var DatabaseEnvironmentVariables = initializeDatabaseConfigs()

var ServerEnvironmentVariables = initializeServerConfigs()

// return environment variables for MySQL
func initializeDatabaseConfigs() DatabaseConfigs {
	godotenv.Load()

	return DatabaseConfigs{
		User: getEnvironmentVariable("DB_USER", "root"),
		Password: getEnvironmentVariable("DB_PASSWORD", "mypassword"),
		Address: fmt.Sprintf("%s:%s", getEnvironmentVariable("DB_HOST", "127.0.0.1"), getEnvironmentVariable("DB_PORT", "3306")),
		Name: getEnvironmentVariable("DB_NAME", "dev-journal-database"),
	}
}

// return environment variables for server
func initializeServerConfigs() ServerConfigs {
	godotenv.Load()

	return ServerConfigs{
		PublicHost: getEnvironmentVariable("PUBLIC_HOST", "http://localhost"),
		Port: getEnvironmentVariable("PORT", "8080"),
	}
}

// return environment variable value if exists, otherwise return default value
func getEnvironmentVariable(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
