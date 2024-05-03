package config

import (
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("Error loading .env file")
	}
}

func FB_SECRET_CREDENTIAL() string {
	loadEnv()
	v, ok := os.LookupEnv("FB_SECRET_CREDENTIAL")
	if !ok {
		panic("FB_SECRET_CREDENTIAL is not set")
	}
	return v
}
