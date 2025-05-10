package controllers

import (
	database "academix/config"
	"academix/models"
	"academix/permissions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStudentList(c *gin.Context) {
	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "user", "view") {
		c.JSON(http.StatusOK, gin.H{"message": "Permission Denied", "role": role})
	}
	var students []models.UserModel
	if err := database.DB.Where("role = ?", "student").Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"students": students})

}

func GetTeachersList(c *gin.Context) {
	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "user", "view") {
		c.JSON(http.StatusOK, gin.H{"message": "Permission Denied", "role": role})
	}
	var teachers []models.UserModel
	if err := database.DB.Where("role = ?", "teacher").Find(&teachers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teachers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"teachers": teachers})
}
func AssignUserToCourse(c *gin.Context) {

	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "course", "addUser") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	type AssignUserRequest struct {
		AssignableUsername string `json:"assignableUsername"`
	}

	var request AssignUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	assignableUsername := request.AssignableUsername

	courseCode := c.Param("courseCode")
	var course models.CourseModel

	if err := database.DB.Preload("Students").Preload("Instructors").Where("Code = ?", courseCode).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	user := getUser(assignableUsername)
	if user.Role == "student" {
		err := database.DB.Model(&course).Association("Students").Append(&user)
		if err != nil {
			return
		}
		err = database.DB.Model(&user).Association("Courses").Append(&course)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Course assigned to student", "course": course})

	} else if user.Role == "teacher" {
		err := database.DB.Model(&course).Association("Instructors").Append(&user)
		if err != nil {
			return
		}
		err = database.DB.Model(&user).Association("TaughtCourses").Append(&course)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Instructors assigned", "course": course})
	}

}

func RemoveUserFromCourse(c *gin.Context) {

	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "course", "addUser") { // assuming "removeUser" permission
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	type RemoveUserRequest struct {
		RemovableUsername string `json:"removableUsername"`
	}

	var request RemoveUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	removableUsername := request.RemovableUsername

	courseCode := c.Param("courseCode")
	var course models.CourseModel

	if err := database.DB.Preload("Students").Preload("Instructors").Where("Code = ?", courseCode).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	user := getUser(removableUsername)

	if user.Role == "student" {

		if err := database.DB.Model(&course).Association("Students").Delete(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove student from the course"})
			return
		}
		if err := database.DB.Model(&user).Association("Courses").Delete(&course); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove course from student"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Student removed successfully", "course": course})
	} else if user.Role == "teacher" {

		if err := database.DB.Model(&course).Association("Instructors").Delete(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove teacher from the course"})
			return
		}
		if err := database.DB.Model(&user).Association("TaughtCourses").Delete(&course); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not remove course from teacher"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Teacher removed successfully", "course": course})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requested user must be a student or teacher"})
	}
}
