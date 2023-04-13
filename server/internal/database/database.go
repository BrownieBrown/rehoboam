package database

import (
	"database/sql"
	"fmt"
	"log"
	"rehoboam/internal/config"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DBManager struct {
	once       sync.Once
	dbInstance *sql.DB
}

func (dbm *DBManager) ConnectDB(config config.DBConfig) (*sql.DB, error) {
	var err error
	dbm.once.Do(func() {
		connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True",
			config.Username, config.Password, config.Hostname, config.DBName)
		dbm.dbInstance, err = sql.Open("mysql", connStr)
		if err != nil {
			log.Fatalf("Unable to connect to the database: %v", err)
			return
		}
		dbm.dbInstance.SetMaxOpenConns(10)
		dbm.dbInstance.SetMaxIdleConns(5)
		dbm.dbInstance.SetConnMaxLifetime(time.Minute * 5)
	})

	return dbm.dbInstance, err
}

func (dbm *DBManager) Close() error {
	if dbm.dbInstance != nil {
		return dbm.dbInstance.Close()
	}
	return nil
}
