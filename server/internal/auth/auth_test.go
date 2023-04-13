package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"rehoboam/helper"
	"rehoboam/internal/models"
	"testing"
)

func setupTestRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1/user")
	api.POST("/signUp", func(c *gin.Context) {
		SignUp(db, c)
	})
	api.POST("/signIn", func(c *gin.Context) {
		SignIn(db, c)
	})
	return r
}

func TestSignUp(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	router := setupTestRouter(db)

	t.Run("Successful registration", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO users").
			ExpectExec().
			WithArgs("test@example.com", sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		user := models.User{
			Email:    "test@example.com",
			Password: "test_password",
		}

		body, err := json.Marshal(user)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/user/signUp", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User already exists", func(t *testing.T) {
		mock.ExpectPrepare("INSERT INTO users").
			ExpectExec().
			WithArgs("duplicate@example.com", sqlmock.AnyArg()).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry"})

		user := models.User{
			Email:    "duplicate@example.com",
			Password: "test_password",
		}

		body, err := json.Marshal(user)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/user/signUp", bytes.NewReader(body))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusConflict, resp.Code)
		assert.Equal(t, "User already exists", helper.GetErrorMessage(resp))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
