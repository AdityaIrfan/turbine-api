package models

import (
	"time"

	"gorm.io/gorm"
)

type DivisionType string

const DivisionType_Engineer = "engineer"

type Division struct {
	Id        string          `json:"id"`
	Type      DivisionType    `json:"type"`
	CreatedAt *time.Time      `gorm:"column:created_at"`
	CreatedBy string          `gorm:"column:created_by"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;<-:update"`
	UpdatedBy *string         `gorm:"column:updated_by"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy *string         `gorm:"column:deleted_by"`
}
