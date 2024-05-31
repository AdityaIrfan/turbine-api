package models

import (
	"time"

	"gorm.io/gorm"
)

type RoleType string

var RoleType_Admin = "admin"
var RoleType_User = "user"

var Roles []Role

type Role struct {
	Id        string          `gorm:"column:id"`
	Type      RoleType        `gorm:"column:type"`
	CreatedAt *time.Time      `gorm:"column:created_at"`
	CreatedBy string          `gorm:"column:created_by"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;<-:update"`
	UpdatedBy *string         `gorm:"column:updated_by"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy *string         `gorm:"column:deleted_by"`
}
