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

func ClearDatabase(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE users")

	if err != nil {
		log.Printf("Error clearing users table: %v", err)
		return err
	}

	return nil
}

func DeleteUserFromDatabaseByEmail(db *sql.DB, email string) error {
	_, err := db.Exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		log.Printf("Error deleting user by email: %v", err)
		return err
	}

	return nil
}

func CreateUser(db *sql.DB, user *models.User) error {
	_, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Password)

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	return nil
}

func UpdateUser(db *sql.DB, email string, updatedUser models.User) error {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}

	if updatedUser.Email != "" {
		_, err = tx.Exec("UPDATE users SET email = ? WHERE email = ?", updatedUser.Email, email)
		if err != nil {
			log.Printf("Error updating email: %v", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	if updatedUser.Password != "" {
		_, err = tx.Exec("UPDATE users SET password = ? WHERE email = ?", updatedUser.Password, email)
		if err != nil {
			log.Printf("Error updating password: %v", err)
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}
