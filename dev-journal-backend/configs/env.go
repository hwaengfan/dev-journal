package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfigs struct {
	User     string
	Password string
	Address  string
	Name     string
}

type ServerConfigs struct {
	PublicHost string
	Port       string
}

type GlobalConfigs struct {
	JWTExpirationInSeconds int64
	JWTSecret              string
}

var DatabaseEnvironmentVariables = initializeDatabaseConfigs()

var ServerEnvironmentVariables = initializeServerConfigs()

var GlobalEnvironmentVariables = initializeGlobalConfigs()

// return environment variables for MySQL
func initializeDatabaseConfigs() DatabaseConfigs {
	godotenv.Load()

	return DatabaseConfigs{
		User:     getEnvironmentVariable("DB_USER", "root"),
		Password: getEnvironmentVariable("DB_PASSWORD", "mypassword"),
		Address:  fmt.Sprintf("%s:%s", getEnvironmentVariable("DB_HOST", "127.0.0.1"), getEnvironmentVariable("DB_PORT", "3306")),
		Name:     getEnvironmentVariable("DB_NAME", "dev-journal-database"),
	}
}

// return environment variables for server
func initializeServerConfigs() ServerConfigs {
	godotenv.Load()

	return ServerConfigs{
		PublicHost: getEnvironmentVariable("PUBLIC_HOST", "http://localhost"),
		Port:       getEnvironmentVariable("PORT", "8080"),
	}
}

// return global environment variables
func initializeGlobalConfigs() GlobalConfigs {
	godotenv.Load()
	return GlobalConfigs{
		JWTExpirationInSeconds: getEnvironmentVariableAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		JWTSecret:              getEnvironmentVariable("JWT_SECRET", "not-so-secret-anymore?"),
	}
}

// return environment variable value if exists, otherwise return default value
func getEnvironmentVariable(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

// return environment variable value as integer if exists, otherwise return default value
func getEnvironmentVariableAsInt(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		intValue, error := strconv.ParseInt(value, 10, 64)
		if error != nil {
			return defaultValue
		}

		return intValue
	}

	return defaultValue
}
