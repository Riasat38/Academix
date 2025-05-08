package models

import (
	"gorm.io/gorm"
	"time"
)

type Assignment struct {
	gorm.Model
	Serial       int
	CourseCode   string      `gorm:"column:course_code;not null"` // Foreign Key
	Course       CourseModel `gorm:"foreignKey:CourseCode;references:Code;constraint:OnDelete:CASCADE"`
	Instructions *string     `gorm:"type:text;default:null"`
	PublishTime  *time.Time  `gorm:"default:null"`
}
type AssignmentSubmission struct {
	gorm.Model              //createdAt is automically there for submission time
	AssignmentID uint       `gorm:"not null"` // Foreign Key
	Assignment   Assignment `gorm:"foreignKey:AssignmentID;references:ID;constraint:OnDelete:CASCADE"`
	StudentID    uint       `gorm:"not null"` // Foreign Key
	Student      UserModel  `gorm:"foreignKey:StudentID;references:ID;constraint:OnDelete:CASCADE"`
	Submission   string     `gorm:"type:varchar(255)"`      // File path or URL
	Marks        *int       `gorm:"default:null"`           // Nullable marks
	Feedback     *string    `gorm:"type:text;default:null"` // Optional

}
