// internal/platform/config/config.go
package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	MongoURI       string
	MongoDB        string
	AllowedOrigins string
	JWTSecret      string
	JWTIssuer      string
	JWTTTLMinutes  int
}

func FromEnv() Config {
	port := getenv("PORT_API", "8080")
	mongoURI := getenv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getenv("MONGO_DB", "portal")
	allowedOrigins := getenv("CORS_ORIGINS", "http://localhost:3000")

	jwtSecret := getenv("JWT_SECRET", "")
	jwtIssuer := getenv("JWT_ISSUER", "djengua-api")

	jwtTTL := atoi(getenv("JWT_TTL_MINUTES", "60"), 60)

	return Config{
		Port:           port,
		MongoURI:       mongoURI,
		MongoDB:        mongoDB,
		AllowedOrigins: allowedOrigins,
		JWTSecret:      jwtSecret,
		JWTIssuer:      jwtIssuer,
		JWTTTLMinutes:  jwtTTL,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func atoi(v string, def int) int {
	i, err := strconv.Atoi(v)
	if err != nil || i <= 0 {
		return def
	}
	return i
}
