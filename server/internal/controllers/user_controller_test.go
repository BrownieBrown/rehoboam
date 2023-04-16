package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"rehoboam/internal/models"
	"strings"
	"testing"
)

func setupTestRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	admin := r.Group("/api/v1/admin")
	admin.GET("/user", func(c *gin.Context) {
		GetAllUsers(db, c)
	})
	admin.GET("/user/:email", func(c *gin.Context) {
		GetUser(db, c)
	})
	admin.DELETE("/user/:email", func(c *gin.Context) {
		DeleteUserByEmail(db, c)
	})
	admin.DELETE("/user", func(c *gin.Context) {
		DeleteAllUsers(db, c)
	})
	admin.POST("/user", func(c *gin.Context) {
		CreateUser(db, c)
	})

	return r
}

func TestGetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("Successful retrieval of all users", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"email"}).
			AddRow("user1@example.com").
			AddRow("user2@example.com").
			AddRow("user3@example.com")

		mock.ExpectQuery("^SELECT email FROM users").
			WillReturnRows(rows)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/admin/user", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var users []models.UserResponse
		err = json.Unmarshal(resp.Body.Bytes(), &users)
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

		req, err := http.NewRequest(http.MethodGet, "/api/v1/admin/user", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), "Error retrieving users")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("Successful retrieval of user by email", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"email"}).
			AddRow("user1@example.com")

		mock.ExpectQuery("^SELECT email FROM users WHERE email = ?").
			WithArgs("user1@example.com").
			WillReturnRows(rows)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/admin/user/user1@example.com", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		var userResponse models.UserResponse
		err = json.Unmarshal(resp.Body.Bytes(), &userResponse)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "user1@example.com", userResponse.Email)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User not found", func(t *testing.T) {
		mock.ExpectQuery("^SELECT email FROM users WHERE email = ?").
			WithArgs("notfound@example.com").
			WillReturnError(sql.ErrNoRows)

		req, err := http.NewRequest(http.MethodGet, "/api/v1/admin/user/notfound@example.com", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Error retrieving user by email", func(t *testing.T) {
		mock.ExpectQuery("^SELECT email FROM users WHERE email = ?").
			WithArgs("error@example.com").
			WillReturnError(errors.New("Error retrieving user"))

		req, err := http.NewRequest(http.MethodGet, "/api/v1/admin/user/error@example.com", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("Successful deletion of user by email", func(t *testing.T) {
		mock.ExpectExec("^DELETE FROM users WHERE email = ?").
			WithArgs("user1@example.com").
			WillReturnResult(sqlmock.NewResult(1, 1))

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/admin/user/user1@example.com", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("Successfully clear database", func(t *testing.T) {
		mock.ExpectExec("^TRUNCATE TABLE users").WillReturnResult(sqlmock.NewResult(0, 0))

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/admin/user", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Unsuccessfully clear database", func(t *testing.T) {
		mock.ExpectExec("^TRUNCATE TABLE users").WillReturnError(errors.New("Error clearing users table"))

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/admin/user", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("Successfully create user", func(t *testing.T) {
		mock.ExpectExec("^INSERT INTO users").
			WithArgs("user1@example.com", "password123").
			WillReturnResult(sqlmock.NewResult(1, 1))

		reqBody := `{"email": "user1@example.com", "password": "password123"}`
		req, err := http.NewRequest(http.MethodPost, "/api/v1/admin/user", strings.NewReader(reqBody))
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Contains(t, resp.Body.String(), "User successfully created")

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Unsuccessfully create user", func(t *testing.T) {
		mock.ExpectExec("^INSERT INTO users").
			WithArgs("user2@example.com", "password123").
			WillReturnError(errors.New("Error creating user"))

		reqBody := `{"email": "user2@example.com", "password": "password123"}`
		req, err := http.NewRequest(http.MethodPost, "/api/v1/admin/user", strings.NewReader(reqBody))
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), "Error creating user")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
