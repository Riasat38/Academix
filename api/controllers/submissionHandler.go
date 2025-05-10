package controllers

import (
	database "academix/config"
	"academix/models"
	"academix/permissions"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func SubmitAssignment(c *gin.Context) {

	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "submission", "post") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	assignmentIDStr := c.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	courseCode := c.Param("courseCode")
	username := c.GetString("username")
	student := getUser(username)

	file, err := c.FormFile("submission")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	// allow only PDF and Word documents (.pdf, .doc, .docx)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" && ext != ".doc" && ext != ".docx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file format"})
		return
	}
	assignment, _ := GetAssignmentByID(uint(assignmentID))
	if assignment.Deadline != nil && time.Now().After(*assignment.Deadline) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Submission deadline has passed"})
		return
	}

	// Define the upload directory and create it if it doesn't exist.
	uploadDir := "uploads/assignments"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if errDir := os.MkdirAll(uploadDir, 0755); errDir != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create upload directory"})
			return
		}
	}

	filename := fmt.Sprintf("%s_%d_%d_%s_%s", courseCode, assignment.Serial, student.ID, student.Name, ext)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	submissionRecord := models.AssignmentSubmission{
		AssignmentID: uint(assignmentID),
		StudentID:    student.ID,
		Submission:   filePath,
	}

	// Insert the record into the database.
	if err := database.DB.Create(&submissionRecord).Error; err != nil {
		// Optionally, remove the file if DB insertion fails to avoid orphaned files.
		err := os.Remove(filePath)
		if err != nil {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save submission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Submission successful",
		"submission": submissionRecord,
	})
}
func GetAssignmentSubmissions(c *gin.Context) {

	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "submission", "viewAll") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	assignmentIDStr := c.Param("assignment_id")
	assignmentID, err := strconv.ParseUint(assignmentIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assignment ID"})
		return
	}

	var submissions []models.AssignmentSubmission
	if err := database.DB.Preload("Student").
		Where("assignment_id = ?", uint(assignmentID)).
		Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submissions": submissions})
}
func UpdateSubmissionFeedback(c *gin.Context) {

	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "submission", "postMarks:Feedback") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	submissionIDStr := c.Param("submission_id")
	submissionID, err := strconv.ParseUint(submissionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	var input struct {
		Marks    int     `json:"marks"`
		Feedback *string `json:"feedback"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var submission models.AssignmentSubmission
	if err := database.DB.First(&submission, uint(submissionID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	submission.Marks = &input.Marks
	submission.Feedback = input.Feedback
	if err := database.DB.Save(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post marks/feedback"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submission": submission})
}

func GetStudentSubmissions(c *gin.Context) {

	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "submission", "getMarks:Feedback") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	//username := c.GetString("username")
	submissionID, err := strconv.Atoi(c.Param("submission_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
	}
	var submission models.AssignmentSubmission
	if err := database.DB.Preload("Student").
		Where("id = ? ", submissionID).
		First(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch student data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submission": submission})
	return
}
