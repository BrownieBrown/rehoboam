package database

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rehoboam/internal/config"
	"rehoboam/internal/models"
	"testing"
)

func TestGetAllUsersFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Successful retrieval of all users", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"email"}).
			AddRow("user1@example.com").
			AddRow("user2@example.com").
			AddRow("user3@example.com")

		mock.ExpectQuery("^SELECT email FROM users").
			WillReturnRows(rows)

		users, err := GetAllUsersFromDB(db)

		require.NoError(t, err)
		assert.Len(t, users, 3)
		assert.Equal(t, "user1@example.com", users[0].Email)
		assert.Equal(t, "user2@example.com", users[1].Email)
		assert.Equal(t, "user3@example.com", users[2].Email)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error retrieving users", func(t *testing.T) {
		mock.ExpectQuery("^SELECT email FROM users").
			WillReturnError(errors.New("Error retrieving users"))

		users, err := GetAllUsersFromDB(db)

		require.Error(t, err)
		assert.Nil(t, users)
		assert.Contains(t, err.Error(), "Error retrieving users")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Successful retrieval of user by email", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"email"}).
			AddRow("user1@example.com")

		mock.ExpectQuery("^SELECT email FROM users WHERE email = ?").
			WithArgs("user1@example.com").
			WillReturnRows(rows)

		user, err := GetUserByEmail(db, "user1@example.com")

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "user1@example.com", user.Email)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery("^SELECT email FROM users WHERE email = ?").
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		user, err := GetUserByEmail(db, "notfound@example.com")

		require.NoError(t, err)
		assert.Nil(t, user)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error retrieving user by email", func(t *testing.T) {
		mock.ExpectQuery("^SELECT email FROM users WHERE email = ?").
			WithArgs("error@example.com").
			WillReturnError(errors.New("Error retrieving user"))

		user, err := GetUserByEmail(db, "error@example.com")

		require.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "Error retrieving user")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestClearDatabase(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Successfully clear database", func(t *testing.T) {
		mock.ExpectExec("^TRUNCATE TABLE user").WillReturnResult(sqlmock.NewResult(0, 0))

		err := ClearDatabase(db)

		require.NoError(t, err)
		assert.Nil(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Unsuccessfully clear database", func(t *testing.T) {
		mock.ExpectExec("^TRUNCATE TABLE user").WillReturnError(errors.New("Error clearing users table"))

		err := ClearDatabase(db)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "Error clearing users table")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Successfully create user", func(t *testing.T) {
		mock.ExpectExec("^INSERT INTO users").
			WithArgs("user1@example.com", "password123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		user := models.User{
			Email:    "user1@example.com",
			Password: "password123",
		}

		err := CreateUser(db, &user)
		require.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Unsuccessfully create user", func(t *testing.T) {
		mock.ExpectExec("^INSERT INTO users").
			WithArgs("user2@example.com", "password123").
			WillReturnError(errors.New("Error creating user"))

		user := models.User{
			Email:    "user2@example.com",
			Password: "password123",
		}

		err := CreateUser(db, &user)
		require.Error(t, err)
		assert.Equal(t, "Error creating user", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDBManager(t *testing.T) {
	t.Run("Connect and close database", func(t *testing.T) {
		cfg := config.DBConfig{
			Username: "username",
			Password: "password",
			Hostname: "hostname",
			DBName:   "dbname",
		}

		dbm := &DBManager{}
		db, err := dbm.ConnectDB(cfg)
		require.NoError(t, err)
		require.NotNil(t, db)

		err = dbm.Close()
		require.NoError(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	t.Run("Successfully update user", func(t *testing.T) {
		email := "user1@example.com"
		updatedUser := models.User{
			Email:    "user1_new@example.com",
			Password: "new_password",
		}

		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE users SET email = (.+) WHERE email = (.+)$").
			WithArgs(updatedUser.Email, email).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectExec("^UPDATE users SET password = (.+) WHERE email = (.+)$").
			WithArgs(updatedUser.Password, email).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := UpdateUser(db, email, updatedUser)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error updating user", func(t *testing.T) {
		email := "user2@example.com"
		updatedUser := models.User{
			Email:    "user2_new@example.com",
			Password: "",
		}

		mock.ExpectBegin()
		mock.ExpectExec("^UPDATE users SET email = (.+) WHERE email = (.+)$").
			WithArgs(updatedUser.Email, email).
			WillReturnError(errors.New("Error updating email"))
		mock.ExpectRollback()

		err := UpdateUser(db, email, updatedUser)
		assert.Error(t, err)
		assert.Equal(t, "Error updating email", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
