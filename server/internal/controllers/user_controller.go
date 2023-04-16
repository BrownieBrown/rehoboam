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
