package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"rehoboam/internal/config"
)

func ConnectDB(config config.DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True",
		config.Username, config.Password, config.Hostname, config.DBName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the database: %v", err)
	}
	return db, nil
}
