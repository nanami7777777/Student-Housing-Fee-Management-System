package config

import "os"

type Config struct {
	DBUrl     string
	JWTSecret string
	HTTPPort  string
	MQUrl     string
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
	mqUrl := os.Getenv("DORM_MQ_URL")
	if mqUrl == "" {
		mqUrl = "amqp://guest:guest@localhost:5672/"
	}
	return Config{
		DBUrl:     dbUrl,
		JWTSecret: secret,
		HTTPPort:  port,
		MQUrl:     mqUrl,
	}
}
