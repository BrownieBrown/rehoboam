package database

import (
	"database/sql"
	"fmt"
	"log"
	"rehoboam/internal/config"
	"rehoboam/internal/models"
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

func GetAllUsersFromDB(db *sql.DB) ([]models.UserResponse, error) {
	rows, err := db.Query("SELECT email FROM users")
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.UserResponse
		err := rows.Scan(&user.Email)
		if err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Error iterating users: %v", err)
		return nil, err
	}

	return users, nil
}

func GetUserByEmail(db *sql.DB, email string) (*models.UserResponse, error) {
	row := db.QueryRow("SELECT email FROM users WHERE email = ?", email)

	var user models.UserResponse
	err := row.Scan(&user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Error scanning user: %v", err)
		return nil, err
	}

	return &user, nil
}
