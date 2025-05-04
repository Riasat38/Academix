package models

import "gorm.io/gorm"

type CourseModel struct {
	gorm.Model
	Code        string `gorm:"unique;not null"`
	Title       string `gorm:"not null"`
	Description string
	Students    []UserModel `gorm:"many2many:user_courses;"`
	Instructors []UserModel `gorm:"many2many:instructor_courses;"`
}
