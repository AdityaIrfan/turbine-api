package models

import (
	"time"
	"turbine-api/helpers"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type DivisionType string

const DivisionType_Engineer = "engineer"

type Division struct {
	Id        string       `gorm:"column:id"`
	Type      DivisionType `gorm:"column:type"`
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
		Type:      d.Type,
		CreatedAt: d.CreatedAt.Format(helpers.DefaultTimeFormat),
		// CreatedBy: d.CreatedBy,
		// UpdatedBy: d.UpdatedBy,
	}

	if d.UpdatedAt != nil {
		res.UpdatedAt = d.UpdatedAt.Format(helpers.DefaultTimeFormat)
	}

	return res
}

type DivisionWriteRequest struct {
	Id   string
	Type DivisionType `json:"Type"`
}

func (d *DivisionWriteRequest) ToModelCreate() *Division {
	return &Division{
		Id:   ulid.Make().String(),
		Type: d.Type,
	}
}

func (d *DivisionWriteRequest) ToModelUpdate() *Division {
	return &Division{
		Id:   d.Id,
		Type: d.Type,
	}
}

type DivisionResponse struct {
	Id        string       `json:"Id"`
	Type      DivisionType `json:"Type"`
	CreatedAt string       `json:"CreatedAt"`
	// CreatedBy string       `json:"CreatedBy"`
	UpdatedAt string `json:"UpdatedAt"`
	// UpdatedBy *string      `json:"UpdatedBy"`
}
