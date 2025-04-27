package main

import (
	database "academix/config"
	"academix/controllers"
	"academix/middleware"
	"academix/models"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()

	database.ConnectDB()
	if database.DB == nil {
		log.Fatal("Database connection failed!")
	}
	fmt.Println("DB instance:", database.DB)
	//database.DB.Debug().AutoMigrate(&models.UserModel{}) 	//used when there were issues
	err := database.DB.AutoMigrate(&models.UserModel{})
	if err != nil {
		return
	}

	fmt.Println("Migration completed successfully!")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Specify allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // HTTP methods allowed
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Headers allowed
		ExposeHeaders:    []string{"Content-Length"},                          // Headers exposed to frontend
		AllowCredentials: true,                                                // Allow cookies/session credentials
		MaxAge:           12 * time.Hour,                                      // Cache preflight requests for 12 hours
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome"})
	})

	router.POST("/academix/signup", controllers.SignUP)
	router.POST("/academix/login", controllers.LogIn)

	authorized := router.Group(`/academix`)
	authorized.Use(middleware.AuthenticateMiddleware())
	{ //below will be the protected routes
		authorized.POST("/logout", controllers.Logout)
		authorized.GET("/profile", controllers.ShowUser)
		//authorized.POST("/course", controllers.PublishCourse)
		authorized.GET("/course", controllers.GetCourses)

	}

	Port := ":8080"
	errP := router.Run(Port)
	if errP != nil {
		return
	}

}
