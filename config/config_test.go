package config

import (
	"testing"
)

func TestLoad_DefaultPort(t *testing.T) {
	t.Setenv("MONGO_URI", "mongodb://localhost:27017")
	t.Setenv("PORT", "")

	cfg, err := Load()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Port != "8080" {
		t.Fatalf("expected port 8080,got %s", cfg.Port)
	}
}

func TestLoad_MissingMonogoURI(t *testing.T) {
	t.Setenv("MONGO_URI", "")
	t.Setenv("PORT", "8080")

	_, err := Load()

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}
