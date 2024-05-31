package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id           string          `gorm:"column:id"`
	Name         string          `gorm:"column:name"`
	Username     string          `gorm:"column:username"`
	DivisionId   string          `gorm:"column:division_id"`
	RoleId       string          `gorm:"column:role_id"`
	PasswordHash string          `gorm:"column:password_hash"`
	PasswordSalt string          `gorm:"column:password_salt"`
	CreatedAt    *time.Time      `gorm:"column:created_at"`
	CreatedBy    string          `gorm:"column:created_by"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at;<-:update"`
	UpdatedBy    *string         `gorm:"column:updated_by"`
	DeletedAt    *gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy    *string         `gorm:"column:deleted_by"`

	Role     *Role     `gorm:"foreignKey:RoleId;references:Id"`
	Division *Division `gorm:"foreignKey:DivisionId;references:Id"`
}
