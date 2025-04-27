package models

import "gorm.io/gorm"

type CourseModel struct {
	gorm.Model
	Code        string
	Title       string
	Instructor  []UserModel `gorm:"many2many:course_instructor"`
	Description string
}
