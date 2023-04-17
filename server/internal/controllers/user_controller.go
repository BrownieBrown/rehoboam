package controllers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
	"rehoboam/internal/database"
	"rehoboam/internal/helper"
	"rehoboam/internal/models"
)

func GetAllUsers(db *sql.DB, c *gin.Context) {
	users, err := database.GetAllUsersFromDB(db)
	if err != nil {
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = models.UserResponse{
			Email: user.Email,
		}
	}

	c.JSON(http.StatusOK, userResponses)
}

func GetUser(db *sql.DB, c *gin.Context) {
	email := c.Param("email")

	user, err := database.GetUserByEmail(db, email)
	if err != nil {
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userResponses := models.UserResponse{Email: user.Email}
	c.JSON(http.StatusOK, userResponses)
}

func DeleteAllUsers(db *sql.DB, c *gin.Context) {
	err := database.ClearDatabase(db)
	if err != nil {
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users table successfully cleared"})
}

func DeleteUserByEmail(db *sql.DB, c *gin.Context) {
	email := c.Param("email")
	err := database.DeleteUserFromDatabaseByEmail(db, email)

	if err != nil {
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Users deleted successfully"})
}

func CreateUser(db *sql.DB, c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		helper.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	err := database.CreateUser(db, &newUser)
	if err != nil {
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User successfully created"})
}

func UpdateUserByEmail(db *sql.DB, c *gin.Context) {
	email := c.Param("email")

	var updatedUser models.User
	err := c.BindJSON(&updatedUser)
	if err != nil {
		helper.HandleError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = database.UpdateUser(db, email, updatedUser)
	if err != nil {
		helper.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}
