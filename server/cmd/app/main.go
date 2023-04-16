package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"rehoboam/internal/auth"
	"rehoboam/internal/config"
	"rehoboam/internal/controllers"
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

	public := r.Group("/api/v1/user")
	{
		public.POST("/signUp", func(c *gin.Context) { auth.SignUp(db, c) })
		public.POST("/signIn", func(c *gin.Context) { auth.SignIn(db, c) })
	}

	admin := r.Group("/api/v1/admin")
	{
		admin.GET("/user", func(c *gin.Context) { controllers.GetAllUsers(db, c) })
		admin.GET("/user/:email", func(c *gin.Context) { controllers.GetUser(db, c) })
		admin.DELETE("/user", func(c *gin.Context) { controllers.DeleteAllUsers(db, c) })
		admin.DELETE("/user/:email", func(c *gin.Context) { controllers.DeleteUserByEmail(db, c) })
	}

	r.Run(":8080")
}
