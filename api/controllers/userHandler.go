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
	role := c.GetString("role")
	if !permissions.ValidatePermission(role, "profile", "edit") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Permission"})
		return
	}
	type Input struct {
		Email    *string `json:"email"`
		Password *string `json:"password"`
	}
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.UserModel
	if err := database.DB.First(&user, username).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
	}

}

func Logout(c *gin.Context) {
	// Remove the token from the browser storage
	c.SetCookie("token", "", -1, "/", "academix.com", false, true) // Expire token

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
	return
}
