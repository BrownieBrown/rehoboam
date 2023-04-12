package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type DBConfig struct {
	Username string
	Password string
	Hostname string
	DBName   string
}

func LoadDBConfig() DBConfig {
	err := godotenv.Load("../../local.env")
	if err != nil {
		log.Fatalf("Error loading local.env file: %v", err)
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOSTNAME")
	dbname := os.Getenv("DB_NAME")

	return DBConfig{
		Username: username,
		Password: password,
		Hostname: host,
		DBName:   dbname,
	}
}
