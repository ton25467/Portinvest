package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
	ServerPort  string
	CORSOrigins []string
	MaxDBConns  int32
}

func loadDotEnv() {
	files := []string{".env", "../.env", "../../.env"}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
					val = val[1 : len(val)-1]
				}
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
		break
	}
}

func Load() *Config {
	loadDotEnv()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/portinves?sslmode=disable"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkeychangeinproduction123"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	maxConnsStr := os.Getenv("DATABASE_MAX_CONNS")
	maxConns := int32(10)
	if maxConnsStr != "" {
		if val, err := strconv.Atoi(maxConnsStr); err == nil {
			maxConns = int32(val)
		}
	}

	return &Config{
		DatabaseURL: dbURL,
		JWTSecret:   secret,
		JWTExpiry:   24 * time.Hour,
		ServerPort:  port,
		CORSOrigins: []string{"*"}, // Allow all in dev, modify in prod
		MaxDBConns:  maxConns,
	}
}
