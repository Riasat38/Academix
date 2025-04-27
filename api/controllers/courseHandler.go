package controllers

import (
	database "academix/config"
	"academix/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PublishCourse(c *gin.Context) {

	var course models.CourseModel

	if err := c.BindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	fmt.Println(course)
	/*
		new_course := models.CourseModel{
			Code:        courseCode,
			Title:       courseTitle,
			Instructor:  courseInstructor,
			Description: description,
		}
	*/
}

func GetCourses(c *gin.Context) {

	var courses []models.CourseModel

	if err := database.DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Error Fetching Data"})
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses})
}
