package models

import (
	"time"

	"gorm.io/gorm"
)

type PltaUnit struct {
	Id        string          `gorm:"column:id"`
	PltaId    string          `gorm:"column:plta_id"`
	Name      string          `gorm:"column:name"`
	Status    bool            `gorm:"column:status"`
	CreatedAt *time.Time      `gorm:"column:created_at"`
	CreatedBy string          `gorm:"column:created_by"`
	UpdatedAt *time.Time      `gorm:"column:updated_at;<-:update"`
	UpdatedBy string          `gorm:"column:updated_by"`
	DeletedAt *gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy string          `gorm:"column:deleted_by"`

	CreatedByUser *User `gorm:"foreignKey:CreatedBy;references:Id"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBy;references:Id"`
	DeletedByUser *User `gorm:"foreignKey:DeletedBy;references:Id"`

	Plta *Plta `gorm:"foreignKey:PltaId;references:Id"`
}

func (p *PltaUnit) IsEmpty() bool {
	return p == nil
}

type PltaUnitWriteRequest struct {
	Id        string `json:"Id" form:"Id"`
	Name      string `json:"Name" form:"Name" validate:"required"`
	Status    bool   `json:"Status" form:"status"`
	WrittenBy string
}

type PltaUnitCreateOrUpdate struct {
	PltaId    string
	Units     []PltaUnitWriteRequest `json:"Units" form:"Units" validate:"required,dive"`
	WrittenBy string
}

func (p *PltaUnitCreateOrUpdate) ToModelCreateOrUpdate() []*PltaUnit {
	units := []*PltaUnit{}

	for _, unit := range p.Units {
		u := &PltaUnit{
			Id:     unit.Id,
			PltaId: p.PltaId,
			Name:   unit.Name,
			Status: unit.Status,
		}

		if u.Id == "" {
			u.CreatedBy = p.WrittenBy
		} else {
			u.UpdatedBy = p.WrittenBy
		}

		units = append(units, u)
	}

	return units
}

type PltaUnitResponse struct {
	Id     string `json:"Id"`
	Name   string `json:"Name"`
	Status bool   `json:"Status"`
}

func (p *PltaUnit) ToResponse() *PltaUnitResponse {
	return &PltaUnitResponse{
		Id:     p.Id,
		Name:   p.Name,
		Status: p.Status,
	}
}

func (p *PltaUnit) IsActive() bool {
	return p.Status
}
