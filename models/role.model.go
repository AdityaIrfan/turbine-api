package models

import (
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
	"pln/AdityaIrfan/turbine-api/helpers"
)

type RoleType string

var RoleType_Admin = "admin"
var RoleType_User = "user"

var Roles []Role

type Role struct {
	Id        string     `gorm:"column:id"`
	Type      RoleType   `gorm:"column:type"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	// CreatedBy string          `gorm:"column:created_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at;<-:update"`
	// UpdatedBy *string         `gorm:"column:updated_by"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`
	// DeletedBy *string         `gorm:"column:deleted_by"`
}

func (r *Role) IsEmpty() bool {
	return r == nil
}

func (r *Role) ToResponse() *RoleResponse {
	res := &RoleResponse{
		Id:        r.Id,
		Type:      r.Type,
		CreatedAt: r.CreatedAt.Format(helpers.DefaultTimeFormat),
		// CreatedBy: r.CreatedBy,
		// UpdatedBy: r.UpdatedBy,
	}

	if r.UpdatedAt != nil {
		res.UpdatedAt = r.UpdatedAt.Format(helpers.DefaultTimeFormat)
	}

	return res
}

type RoleWriteRequest struct {
	Id   string
	Type RoleType `gorm:"column:type" validate:"required"`
}

func (r *RoleWriteRequest) ToModelCreate() *Role {
	return &Role{
		Id:   ulid.Make().String(),
		Type: r.Type,
	}
}

func (r *RoleWriteRequest) ToModelUpdate() *Role {
	return &Role{
		Id:   ulid.Make().String(),
		Type: r.Type,
	}
}

type RoleResponse struct {
	Id        string   `json:"Id"`
	Type      RoleType `json:"Type"`
	CreatedAt string   `json:"CreatedAt"`
	// CreatedBy string   `json:"CreatedBy"`
	UpdatedAt string `json:"UpdatedAt"`
	// UpdatedBy *string  `json:"UpdatedBy"`
}
