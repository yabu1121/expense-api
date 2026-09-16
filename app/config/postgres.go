package config

import (
	"fmt"
	"os"
)

func DatabaseURL() (string, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		return "", fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	return url, nil
}
