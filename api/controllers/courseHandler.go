package controllers

import (
	database "academix/config"
	"academix/models"
	"academix/permissions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStudentCourses(username string) []models.CourseModel {
	// Assuming authentication provides username
	var user models.UserModel

	// Find the student and preload their enrolled courses
	if err := database.DB.Model(&user).Preload("Courses").
		Where("username = ?", username).Find(&user).Error; err != nil {
		return []models.CourseModel{}
	}

	return user.Courses
}

func GetInstructorCourses(username string) []models.CourseModel {
	// Assuming authentication provides username
	var user models.UserModel

	// Find the instructor and preload their taught courses
	if err := database.DB.Model(&user).Preload("TaughtCourses").
		Where("username = ?", username).Find(&user).Error; err != nil {
		return []models.CourseModel{}
	}

	return user.TaughtCourses
}

func ViewAllCourses(c *gin.Context) { //view all courses

	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "course", "viewAll") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	var courses []models.CourseModel

	if err := database.DB.Preload("Instructors").Find(&courses).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Error Fetching Data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func ViewOwnCourses(c *gin.Context) {
	role := c.GetString("role")
	username := c.GetString("username")

	if !permissions.ValidatePermission(role, "course", "viewOwn") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}
	var courses []models.CourseModel
	/*	alternative
		database.DB.Where("username = ?", user.Username).Find(&courses)
	*/
	if role == "student" {
		courses = GetStudentCourses(username)
	}
	if role == "teacher" {
		courses = GetInstructorCourses(username)
	}

	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func EnrollCourse(c *gin.Context) {
	courseCode := c.Param("courseCode")

	var course models.CourseModel

	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "course", "enroll") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	//finding user
	username := c.GetString("username")
	user := getUser(username)
	//finding course

	if err := database.DB.Where("Code = ?", courseCode).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	courses := GetStudentCourses(username)
	for _, course := range courses {
		if course.Code == courseCode {
			c.JSON(http.StatusConflict, gin.H{"error": "Already enrolled in the course"})
			return
		}
	}
	//adding student to the course
	if err := database.DB.Model(&course).Association("Students").Append(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not added in course"})
		return
	}
	//adding the course to students course
	if err := database.DB.Model(&user).Association("Courses").Append(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll in course"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student enrolled successful", "Course": course})

}
func CreateCourse(c *gin.Context) {
	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "course", "create") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	var course models.CourseModel
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var existingCourse models.CourseModel
	if err := database.DB.Where("Code = ?", course.Code).First(&existingCourse).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Course already exists"})
		return
	}

	if err := database.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course created", "course": course})
}

func AssignUser(c *gin.Context) {

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

func ViewCourse(c *gin.Context) {

	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "course", "view") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}
	courseCode := c.Param("courseCode")
	username := c.GetString("username")

	if role == "student" {
		var user models.UserModel

		// Find the student and preload their enrolled courses
		if err := database.DB.Model(&user).Preload("Courses").Preload("Courses.Assignments").
			Where("username = ?", username).Find(&user).Error; err != nil {
			return
		}
		for _, course := range user.Courses {
			if course.Code == courseCode {
				c.JSON(http.StatusOK, gin.H{"course": course})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Course not found", "course": courseCode})
		return
	}
	if role == "teacher" {

		var user models.UserModel
		if err := database.DB.Model(&user).Preload("TaughtCourses").Preload("Courses.Assignments").
			Where("username = ?", username).Find(&user).Error; err != nil {
			return
		}
		for _, course := range user.TaughtCourses {
			if course.Code == courseCode {
				c.JSON(http.StatusOK, gin.H{"course": course})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "Course not found", "course": courseCode})
		return
	}

	/*
		if err := database.DB.Model(&models.CourseModel{}).
			Preload("Assignments").
			Joins("JOIN user_courses ON user_courses.course_id = courses.id").
			Joins("JOIN users ON users.id = user_courses.user_id").
			Where("users.Username = ? AND courses.Code = ?", username, courseCode).
			First(&course).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Course viewed", "course": course})
		return*/

	if role == "admin" {
		var course models.CourseModel
		database.DB.Preload("Students").Preload("Instructors").Preload("Assignments").Where("Code =?", courseCode).Find(&course)
		c.JSON(http.StatusOK, gin.H{"course": course})
		return
	}
}
