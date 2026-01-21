package config

import "os"

type Config struct {
	Port           string
	MongoURI       string
	MongoDB        string
	AllowedOrigins string
}

func FromEnv() Config {
	port := getenv("PORT_API", "8080")
	mongoURI := getenv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getenv("MONGO_DB", "portal")
	allowedOrigins := getenv("CORS_ORIGINS", "http://localhost:3000")
	return Config{Port: port, MongoURI: mongoURI, MongoDB: mongoDB, AllowedOrigins: allowedOrigins}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
