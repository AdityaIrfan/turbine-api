package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Tower struct {
	Id        string          `gorm:"column:id"`
	Name      string          `gorm:"column:name"`
	CreatedAt *time.Time      `gorm:"column:created_at"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;<-:update"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (t *Tower) IsEmpty() bool {
	return t == nil
}

func (t *Tower) ToResponse() *TowerResponse {
	return &TowerResponse{
		Id:   t.Id,
		Name: t.Name,
	}
}

type TowerWriteRequest struct {
	Id   string
	Name string `json:"Name" validate:"required"`
}

func (t *TowerWriteRequest) ToModelCreate() *Tower {
	return &Tower{
		Id:   ulid.Make().String(),
		Name: t.Name,
	}
}

func (t *TowerWriteRequest) ToModelUpdate() *Tower {
	return &Tower{
		Id:   t.Id,
		Name: t.Name,
	}
}

type TowerResponse struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}
