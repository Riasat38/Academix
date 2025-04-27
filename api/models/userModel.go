package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model        //using for auto Id, createdAt, UpdatedAt fields
	Name       string `gorm:"not null"`
	Username   string `gorm:"unique;required"`
	Email      string `gorm:"unique;required"`
	//Courselist []CourseModel `gorm:"many2many:user_courses"`
	Password string `gorm:"not null;required"`
	Role     string `gorm:"not null;required"`
}

func (user *UserModel) CheckPassword(inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword))
	return err == nil
}
