package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func FB_SECRET_CREDENTIAL() string {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	v, ok := os.LookupEnv("FB_SECRET_CREDENTIAL")
	if !ok {
		panic("FB_SECRET_CREDENTIAL is not set")
	}
	return v
}
