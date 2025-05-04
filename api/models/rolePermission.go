package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type Permission struct {
	gorm.Model
	Action   string `gorm:"not null"`
	Resource string `gorm:"not null"`
}

type RolePermission struct {
	RoleID       uint `gorm:"not null"`
	PermissionID uint `gorm:"not null"`
}

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	RoleID   uint   `gorm:"not null"`
}
