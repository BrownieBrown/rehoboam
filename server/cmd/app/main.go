package main

import (
	"github.com/gin-contrib/cors"
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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true                                             // Allow all origins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}            // Specify allowed methods
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"} // Specify allowed headers

	// Apply the CORS middleware to the router
	r.Use(cors.New(corsConfig))

	public := r.Group("/api/v1/auth")
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
		admin.POST("/user", func(c *gin.Context) { controllers.CreateUser(db, c) })
		admin.PUT("/user/:email", func(c *gin.Context) { controllers.UpdateUserByEmail(db, c) })
	}

	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
