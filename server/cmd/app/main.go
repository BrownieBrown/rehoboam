package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"rehoboam/internal/auth"
	"rehoboam/internal/config"
	"rehoboam/internal/database"
)

var dbManager database.DBManager

func main() {
	dbConfig := config.LoadDBConfig()
	db, err := dbManager.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer dbManager.Close()

	r := gin.Default()

	apiV1 := r.Group("/api/v1/user")
	{
		apiV1.POST("/signUp", func(c *gin.Context) { auth.SignUp(db, c) })
		apiV1.POST("/signIn", func(c *gin.Context) { auth.SignIn(db, c) })
	}

	r.Run(":8080")
}
