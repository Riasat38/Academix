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
	err := database.DB.AutoMigrate(&models.UserModel{}, &models.CourseModel{}, &models.Assignment{}, &models.AssignmentSubmission{})
	if err != nil {
		return
	}

	fmt.Println("Migration completed successfully!")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Specify allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // HTTP methods allowed
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
		authorized.PUT("/profile", controllers.EditProfile)

		//Course Routes
		authorized.GET("/course", controllers.ViewAllCourses) //browse
		authorized.GET("/own-course", controllers.ViewOwnCourses)
		authorized.GET("/course/:courseCode", controllers.ViewCourse) //viewing a course with details including instructor assignments and everything
		authorized.POST("/enroll-course/:courseCode", controllers.EnrollCourse)
		authorized.POST("/create-course", controllers.CreateCourse)
		authorized.PUT("/course/:courseCode", controllers.EditCourse)

		adminGroup := authorized.Group("admin")
		{
			adminGroup.GET("/student-list", controllers.GetStudentList)
			adminGroup.GET("/teacher-list", controllers.GetTeachersList)
			adminGroup.POST("/assign-user/:courseCode", controllers.AssignUserToCourse) //admin: assigning teacher or student
			adminGroup.DELETE("/remove-user/:courseCode", controllers.RemoveUserFromCourse)
		}

		//Assignment
		authorized.POST("/:courseCode/assignment", controllers.CreateAssignment)
		authorized.GET("/:courseCode/assignments", controllers.GetAllAssignments)            //get all assignments of a course [assignment1,assignment2,assignment3]
		authorized.GET("/:courseCode/assignments/:assignment_id", controllers.GetAssignment) //one assignment
		authorized.PUT("/:courseCode/assignments/:assignment_id", controllers.UpdateAssignment)
		authorized.DELETE("/:courseCode/assignments/:assignment_id", controllers.DeleteAssignment)
		//submission
		authorized.POST("/:courseCode/assignment/:assignment_id", controllers.SubmitAssignment)
		authorized.GET("/:courseCode/assignment/:assignment_id/submissions", controllers.GetAssignmentSubmissions) //teacher getting all the submission for that assignment

		authorized.GET("/submission/:submission_id", controllers.GetStudentSubmissions)     //student seeing his submission for that assignment
		authorized.PUT("/submissions/:submission_id", controllers.UpdateSubmissionFeedback) //marks update

	}

	Port := ":8080"
	errP := router.Run(Port)
	if errP != nil {
		return
	}

}
