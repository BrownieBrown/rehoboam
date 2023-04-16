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
	"testing"
)

func setupTestRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	admin := r.Group("/api/v1/admin")
	admin.GET("/user", func(c *gin.Context) {
		GetAllUsers(db, c)
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

		req, err := http.NewRequest("GET", "/api/v1/admin/user", nil)
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

		req, err := http.NewRequest("GET", "/api/v1/admin/user", nil)
		require.NoError(t, err)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		assert.Contains(t, resp.Body.String(), "Error retrieving users")

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
