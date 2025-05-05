package controllers

import (
	database "academix/config"
	"academix/models"
	"academix/permissions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateAssignment(c *gin.Context) {
	//username := c.GetString("username")
	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "assignment", "create") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}
	courseCode := c.Param("courseCode")
	var course models.CourseModel
	if err := database.DB.Where("code = ?", courseCode).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var input struct {
		Serial      uint       `json:"serial"`
		Instruction *string    `json:"instruction"`
		PublishTime *time.Time `json:"publishTime"` // Time when assignment becomes visible
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	newAssignment := models.Assignment{
		Serial:       input.Serial,
		Course:       course,
		Instructions: input.Instruction,
		PublishTime:  input.PublishTime,
	}
	database.DB.Create(&newAssignment)
	c.JSON(http.StatusOK, gin.H{"message": "Assignment created successfully", "assignment": newAssignment})
}

func GetAssignments(c *gin.Context) {

	role := c.GetString("role")
	courseCode := c.Param("courseCode")
	username := c.GetString("username")
	if !permissions.ValidatePermission(role, "assignment", "view") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	var assignments []models.Assignment

	currentTime := time.Now()

	if role == "student" {
		courses := GetStudentCourses(username)
		for _, course := range courses {
			if courseCode == course.Code {
				database.DB.Where("CourseCode=  ? AND PublishTime <= ?", course.Code, currentTime).Find(&assignments)
				c.JSON(http.StatusOK, gin.H{"assignments": assignments})
			}
		}

	}
	if role == "teacher" {
		courses := GetInstructorCourses(username)
		for _, course := range courses {
			if courseCode == course.Code {
				database.DB.Where("CourseCode =? ", course.Code).Find(&assignments)
				c.JSON(http.StatusOK, gin.H{"assignments": assignments})
			}
		}

	}

}
