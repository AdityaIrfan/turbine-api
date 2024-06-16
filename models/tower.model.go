package models

import (
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Tower struct {
	Id         string          `gorm:"column:id"`
	Name       string          `gorm:"column:name"`
	UnitNumber string          `gorm:"unit_number"`
	CreatedAt  *time.Time      `gorm:"column:created_at"`
	UpdatedAt  *time.Time      `gorm:"column:updated_at;<-:update"`
	DeletedAt  *gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (t *Tower) IsEmpty() bool {
	return t == nil
}

func (t *Tower) ToResponse() *TowerResponse {
	return &TowerResponse{
		Id:         t.Id,
		Name:       t.Name,
		UnitNumber: t.UnitNumber,
	}
}

func (t *Tower) ToResponseMaster() *TowerResponseMaster {
	return &TowerResponseMaster{
		Id:   t.Id,
		Name: fmt.Sprintf("%v - %v", t.Name, t.UnitNumber),
	}
}

type TowerWriteRequest struct {
	Id         string
	Name       string `json:"Name" form:"Name" validate:"required"`
	UnitNumber string `json:"UnitNumber" form:"UnitNumber" validate:"required,max=20"`
}

func (t *TowerWriteRequest) ToModelCreate() *Tower {
	return &Tower{
		Id:         ulid.Make().String(),
		Name:       t.Name,
		UnitNumber: t.UnitNumber,
	}
}

func (t *TowerWriteRequest) ToModelUpdate() *Tower {
	return &Tower{
		Id:         t.Id,
		Name:       t.Name,
		UnitNumber: t.UnitNumber,
	}
}

type TowerResponse struct {
	Id         string `json:"Id"`
	Name       string `json:"Name"`
	UnitNumber string `json:"UnitNumber"`
}

type TowerResponseMaster struct {
	Id   string `json:"Id"`
	Name string `json:"Name"`
}
