package controllers

import (
	"academix/auth"
	database "academix/config"
	"academix/models"
	"academix/permissions"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

func ShowUser(c *gin.Context) {

	username := c.GetString("username")

	var user models.UserModel
	if err := database.DB.Omit("Password").Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profile": user})
}

func getUser(username string) models.UserModel {
	var user models.UserModel
	err := database.DB.Omit("Password").Where("username = ?", username).First(&user).Error

	if err != nil {
		log.Fatal(fmt.Sprintf("User '%s' not found: %v", username, err))
	}
	return user
}

func SignUP(c *gin.Context) {

	name := c.PostForm("name")
	password := c.PostForm("password")
	username := c.PostForm("username")
	email := c.PostForm("email")
	role := c.PostForm("role")

	var existingUser models.UserModel
	result := database.DB.Where("username = ?", username).Or("email = ?", email).First(&existingUser)

	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	} else {
		hashedPassword, passErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if passErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error hashing password"})
		}
		newUser := models.UserModel{
			Name:     name,
			Password: string(hashedPassword),
			Username: username,
			Email:    email,
			Role:     role,
		}
		if err := database.DB.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		} else {
			token, err := auth.GenerateToken(newUser.Username, newUser.Role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"msg": "User created", "User": newUser, "Token": token, "redirect": "/academix/profile"})
		}
	}
}

func LogIn(c *gin.Context) {

	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.UserModel

	if err := database.DB.Where("username = ?", request.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}

	if !user.CheckPassword(request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}

	token, err := auth.GenerateToken(user.Username, user.Role)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "message": "Login successful", "redirect": "/academix/profile"})

}

func EditProfile(c *gin.Context) {
	username := c.GetString("username")

	var input struct {
		Email *string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.UserModel
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Email != nil && *input.Email != "" {
		// Check if email is already in use
		var existingUser models.UserModel
		if err := database.DB.Where("email = ?", *input.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email is already taken"})
			return
		}
		user.Email = *input.Email
		if err := database.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
			return
		}
	}
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "profile": user})
	return
}

func Logout(c *gin.Context) {
	// Remove the token from the browser storage
	c.SetCookie("token", "", -1, "/", "academix.com", false, true) // Expire token

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
	return
}

func UpdateUserPassword(c *gin.Context) {

	username := c.GetString("username")
	role := c.GetString("role")

	if !permissions.ValidatePermission(role, "profile", "edit") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}

	var input struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword"`     // optional: change password request
		ConfirmPassword string `json:"confirmPassword"` // optional: must match newPassword if provided
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input. Make sure you provide your current password and any new values you want to change."})
		return
	}
	var user models.UserModel
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	if strings.TrimSpace(input.NewPassword) != "" {
		if input.NewPassword != input.ConfirmPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirmation do not match"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process the new password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "user": user})
}
