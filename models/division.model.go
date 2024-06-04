package models

import (
	"time"
	"turbine-api/helpers"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type DivisionName string

const DivisionName_Engineer = "engineer"

type Division struct {
	Id        string       `gorm:"column:id"`
	Name      DivisionName `gorm:"column:name"`
	CreatedAt *time.Time   `gorm:"column:created_at"`
	// CreatedBy string          `gorm:"column:created_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at;<-:update"`
	// UpdatedBy *string         `gorm:"column:updated_by"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`
	// DeletedBy *string         `gorm:"column:deleted_by"`
}

func (d *Division) IsEmpty() bool {
	return d == nil
}

func (d *Division) ToResponse() *DivisionResponse {
	res := &DivisionResponse{
		Id:        d.Id,
		Name:      d.Name,
		CreatedAt: d.CreatedAt.Format(helpers.DefaultTimeFormat),
		// CreatedBy: d.CreatedBy,
		// UpdatedBy: d.UpdatedBy,
	}

	if d.UpdatedAt != nil {
		res.UpdatedAt = d.UpdatedAt.Format(helpers.DefaultTimeFormat)
	}

	return res
}

func (d *Division) ToMasterResponse() *DivisionMasterResponse {
	return &DivisionMasterResponse{
		Id:   d.Id,
		Name: d.Name,
	}
}

type DivisionWriteRequest struct {
	Id   string
	Name DivisionName `json:"Name" validate:"required"`
}

func (d *DivisionWriteRequest) ToModelCreate() *Division {
	return &Division{
		Id:   ulid.Make().String(),
		Name: d.Name,
	}
}

func (d *DivisionWriteRequest) ToModelUpdate() *Division {
	return &Division{
		Id:   d.Id,
		Name: d.Name,
	}
}

type DivisionResponse struct {
	Id        string       `json:"Id"`
	Name      DivisionName `json:"Name"`
	CreatedAt string       `json:"CreatedAt"`
	// CreatedBy string       `json:"CreatedBy"`
	UpdatedAt string `json:"UpdatedAt"`
	// UpdatedBy *string      `json:"UpdatedBy"`
}

type DivisionMasterResponse struct {
	Id   string       `json:"Id"`
	Name DivisionName `json:"Name"`
}
