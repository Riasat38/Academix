package models

import (
	"gorm.io/gorm"
	"time"
)

type Assignment struct {
	gorm.Model
	Serial       uint        `gorm:"increment"`
	CourseCode   uint        `gorm:"not null"` // Foreign Key
	Course       CourseModel `gorm:"foreignKey:CourseCode;constraint:OnDelete:CASCADE"`
	Instructions *string     `gorm:"type:text;default:null"`
	PublishTime  *time.Time  `gorm:"default:null"`
}
type AssignmentSubmission struct {
	gorm.Model
	AssignmentID uint       `gorm:"not null"` // Foreign Key
	Assignment   Assignment `gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE"`
	StudentID    uint       `gorm:"not null"` // Foreign Key
	Student      UserModel  `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE"`
	Submission   string     `gorm:"type:varchar(255)"`      // File path or URL
	Marks        *int       `gorm:"default:null"`           // Nullable marks
	Feedback     *string    `gorm:"type:text;default:null"` // Optional feedback
}
