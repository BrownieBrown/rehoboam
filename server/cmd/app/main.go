package main

import (
	"database/sql"
	"log"
	"rehoboam/internal/config"
	"rehoboam/internal/database"
)

func main() {
	dbConfig := config.LoadDBConfig()
	db, err := database.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing the database: %v", err)
		}
	}(db)
}
