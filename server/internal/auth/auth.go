package auth

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"rehoboam/internal/models"
)

func hashAndSalt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func comparePasswords(hashedPwd, plainPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
	return err == nil
}

func handleError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func prepareAndExecute(db *sql.DB, query string, args ...interface{}) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.Exec(args...)
}

func SignUp(db *sql.DB, c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := hashAndSalt(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	_, err = prepareAndExecute(db, "INSERT INTO users (email, password) VALUES (?, ?)", user.Email, hashedPassword)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			handleError(c, http.StatusConflict, "User already exists")
		} else {
			handleError(c, http.StatusInternalServerError, "Error creating user")
		}
		return
	}

	c.Status(http.StatusCreated)
}

func SignIn(db *sql.DB, c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	var hashedPassword string
	stmt, err := db.Prepare("SELECT password FROM users WHERE email=?")
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Error preparing query")
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(user.Email).Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			handleError(c, http.StatusUnauthorized, "Invalid email or password")
		} else {
			handleError(c, http.StatusInternalServerError, "Error retrieving user")
		}
		return
	}

	if !comparePasswords(hashedPassword, user.Password) {
		handleError(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	c.Status(http.StatusOK)
}
