package config

import "os"

type Config struct {
	DBUrl     string
	JWTSecret string
	HTTPPort  string
}

func Load() Config {
	dbUrl := os.Getenv("DORM_DB_URL")
	if dbUrl == "" {
		dbUrl = "postgres://postgres:postgres@localhost:15432/dormdb?sslmode=disable"
	}
	secret := os.Getenv("DORM_JWT_SECRET")
	if secret == "" {
		secret = "change-this-secret"
	}
	port := os.Getenv("DORM_HTTP_PORT")
	if port == "" {
		port = ":8080"
	}
	return Config{
		DBUrl:     dbUrl,
		JWTSecret: secret,
		HTTPPort:  port,
	}
}

