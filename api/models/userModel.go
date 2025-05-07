package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model                           //using for auto Id, createdAt, UpdatedAt fields
	Name          string                 `gorm:"not null"`
	Username      string                 `gorm:"unique;required"`
	Email         string                 `gorm:"unique;required"`
	Password      string                 `gorm:"not null;required" json:"-"`
	Role          string                 `gorm:"not null;required"`
	Courses       []CourseModel          `gorm:"many2many:user_courses;"`
	TaughtCourses []CourseModel          `gorm:"many2many:instructor_courses;"`
	Submissions   []AssignmentSubmission `gorm:"foreignKey:StudentID;references:ID"`
}

func (user *UserModel) CheckPassword(inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword))
	return err == nil
}
