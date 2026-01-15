package config

import "os"

type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
}

func FromEnv() Config {
	port := getenv("PORT", "8080")
	mongoURI := getenv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getenv("MONGO_DB", "portal")
	return Config{Port: port, MongoURI: mongoURI, MongoDB: mongoDB}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
