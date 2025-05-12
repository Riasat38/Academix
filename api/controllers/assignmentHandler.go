package controllers

import (
	database "academix/config"
	"academix/models"
	"academix/permissions"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func parseTime(timeStr string) (time.Time, error) {
	dt, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time %q: %w", timeStr, err)
	}
	return dt, nil
}

func CreateAssignment(c *gin.Context) {
	//username := c.GetString("username")
	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "assignment", "create") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}
	var courseCode = c.Param("courseCode")
	var course models.CourseModel

	if err := database.DB.Where("Code = ?", courseCode).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	serialStr := c.PostForm("serial")
	instruction := c.PostForm("instruction")
	publishTimeStr := c.PostForm("publishTime")
	deadlineStr := c.PostForm("deadline")

	serial, err := strconv.Atoi(serialStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid serial value"})
		return
	}

	publishTime, err := parseTime(publishTimeStr)
	deadline, err := parseTime(deadlineStr)

	questionFile, err := c.FormFile("question")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Question file not provided"})
		return
	}

	// Only allow PDF and Word documents.
	ext := strings.ToLower(filepath.Ext(questionFile.Filename))
	if ext != ".pdf" && ext != ".doc" && ext != ".docx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format for question"})
		return
	}

	uploadDir := "uploads/assignment_questions"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if errDir := os.MkdirAll(uploadDir, 0755); errDir != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create upload directory"})
			return
		}
	}
	filename := fmt.Sprintf("%s_Assignment%d_question_%s", course.Code, serial, ext)
	filePath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(questionFile, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save question file"})
		return
	}

	newAssignment := models.Assignment{
		Serial:       serial,
		CourseCode:   course.Code,
		Instructions: &instruction,
		PublishTime:  &publishTime,
		Deadline:     &deadline,
		Question:     filePath,
	}
	if err := database.DB.Create(&newAssignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment created successfully", "assignment": newAssignment})
}

func GetAllAssignments(c *gin.Context) {

	role := c.GetString("role")
	courseCode := c.Param("courseCode")

	username := c.GetString("username")
	if !permissions.ValidatePermission(role, "assignment", "viewAll") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	currentTime := time.Now()
	fmt.Println("currentTime", currentTime)
	var assignments []models.Assignment
	var courses []models.CourseModel

	if role == "student" {
		courses = GetStudentCourses(username)
		for _, course := range courses {
			if course.Code == courseCode {
				if err := database.DB.Where("course_code =? AND publish_time <= ?", courseCode, currentTime).Find(&assignments).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
					return
				}
				break
			}

		}
	}
	if role == "teacher" {
		courses = GetInstructorCourses(username)
		for _, course := range courses {
			if course.Code == courseCode {
				if err := database.DB.Where("course_code = ?", courseCode).Find(&assignments).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
					return
				}
				break
			}
		}

	}
	if role == "admin" {
		if err := database.DB.Where("course_code = ?", courseCode).Find(&assignments).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignments"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"assignments": assignments})
	return

}
func GetAssignmentByID(assignmentID uint) (models.Assignment, error) {
	var assignment models.Assignment

	if err := database.DB.Where("id = ?", assignmentID).First(&assignment).Error; err != nil {
		return assignment, err
	}

	return assignment, nil
}
func GetAssignment(c *gin.Context) {
	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "assignment", "view") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}
	assignmentID, err := strconv.Atoi(c.Param("assignment_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment id"})
	}

	var assignment models.Assignment
	assignment, err = GetAssignmentByID(uint(assignmentID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assignment"})
	}

	var encodedQuestion string
	if assignment.Question != "" {
		// Read the file content from disk.
		fileContent, err := os.ReadFile(assignment.Question)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read question file"})
			return
		}
		encodedQuestion = base64.StdEncoding.EncodeToString(fileContent)
	}

	c.JSON(http.StatusOK, gin.H{
		"assignment":    assignment,
		"question_file": encodedQuestion,
	})

}

func UpdateAssignment(c *gin.Context) {

	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "assignment", "edit") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	assignmentIDStr := c.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var assignment models.Assignment
	if err := database.DB.Where("id = ?", assignmentID).First(&assignment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	var input struct {
		Serial      *int       `json:"serial"`      // Pointer so that if omitted, we don't update it.
		Instruction *string    `json:"instruction"` // Pointer so that we can accept null or missing values.
		PublishTime *time.Time `json:"publishTime"` // Pointer to update only when provided.
		Deadline    *time.Time `json:"deadline"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	if input.Serial != nil {
		assignment.Serial = *input.Serial
	}
	if input.Instruction != nil {
		assignment.Instructions = input.Instruction
	}
	if input.PublishTime != nil {
		assignment.PublishTime = input.PublishTime
	}
	if input.Deadline != nil {
		assignment.Deadline = input.Deadline
	}

	questionFile, fileErr := c.FormFile("question")
	if fileErr == nil {
		// Validate file extension.
		ext := strings.ToLower(filepath.Ext(questionFile.Filename))
		if ext != ".pdf" && ext != ".doc" && ext != ".docx" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format for question"})
			return
		}

		uploadDir := "uploads/assignment_questions"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			if errDir := os.MkdirAll(uploadDir, 0755); errDir != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create upload directory"})
				return
			}
		}

		filename := fmt.Sprintf("%s_Assignment%d_question_%s", assignment.CourseCode, assignment.Serial, ext)
		filePath := filepath.Join(uploadDir, filename)

		if err := c.SaveUploadedFile(questionFile, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated question file"})
			return
		}
		assignment.Question = filePath
	}

	if err := database.DB.Save(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Assignment updated successfully",
		"assignment": assignment,
	})
	return
}
func DeleteAssignment(c *gin.Context) {

	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "assignment", "delete") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	assignmentIDStr := c.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var assignment models.Assignment
	if err := database.DB.Where("id = ?", assignmentID).First(&assignment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	if err := database.DB.Unscoped().Delete(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete assignment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assignment deleted successfully"})
	return
}
