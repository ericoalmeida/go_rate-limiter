package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warn! .env file not found. Using environment variables")
	}
}

func GetEnv(key string) string {
	envValue := os.Getenv(key)

	if envValue == "" {
		log.Fatal("Environment variable not defined")
	}

	return envValue
}

func GetEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}
