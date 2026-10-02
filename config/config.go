package config

import (
	"errors"
	"os"
)

type Config struct {
	MongoURI string
	Port     string
}

func Load() (Config, error) {

	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")

	if mongoURI == "" {
		return Config{}, errors.New("MONGO_URI is required ")
	}

	if port == "" {
		port = "8080"
	}

	return Config{
		MongoURI: mongoURI,
		Port:     port,
	}, nil
}
