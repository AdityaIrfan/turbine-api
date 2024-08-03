package models

import (
	"fmt"
	"pln/AdityaIrfan/turbine-api/helpers"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

var PltaDefaultSort = map[string]string{
	"Name":      "name",
	"CreatedAt": "createdat",
}

type PltaRadiusType string

const (
	PltaRadiusType_Meter     PltaRadiusType = "meter"
	PltaRadiusType_Kilometer PltaRadiusType = "kilometer"
)

type Plta struct {
	Id           string          `gorm:"column:id"`
	Name         string          `gorm:"column:name"`
	Status       bool            `gorm:"column:status"`
	Long         float64         `gorm:"column:long"`
	Lat          float64         `gorm:"column:lat"`
	RadiusStatus bool            `gorm:"column:radius_status"`
	Radius       float64         `gorm:"column:radius"`
	RadiusType   PltaRadiusType  `gorm:"column:radius_type"`
	CreatedAt    *time.Time      `gorm:"column:created_at"`
	CreatedBy    string          `gorm:"column:created_by"`
	UpdatedAt    *time.Time      `gorm:"column:updated_at;<-:update"`
	UpdatedBy    string          `gorm:"column:updated_by"`
	DeletedAt    *gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy    string          `gorm:"column:deleted_by"`

	CreatedByUser *User       `gorm:"foreignKey:CreatedBy;references:Id"`
	UpdatedByUser *User       `gorm:"foreignKey:UpdatedBy;references:Id"`
	DeletedByUser *User       `gorm:"foreignKey:DeletedBy;references:Id"`
	PltaUnits     []*PltaUnit `gorm:"foreignKey:PltaId"`
}

func (t *Plta) IsEmpty() bool {
	return t == nil
}

func (t *Plta) ToResponse() *PltaResponse {
	plta := &PltaResponse{
		Id:           t.Id,
		Name:         t.Name,
		Status:       t.Status,
		Long:         t.Long,
		Lat:          t.Lat,
		RadiusStatus: t.RadiusStatus,
		Radius:       t.Radius,
		RadiusType:   t.RadiusType,
		CreatedAt:    t.CreatedAt.Format(helpers.DefaultTimeFormat),
		PltaUnits:    []*PltaUnitResponse{},
	}

	if t.CreatedByUser != nil {
		plta.CreatedBy = t.CreatedByUser.Name
	}

	if t.UpdatedByUser != nil {
		plta.UpdatedAt = t.UpdatedAt.Format(helpers.DefaultTimeFormat)
		plta.UpdatedBy = t.UpdatedByUser.Name
	}

	for _, unit := range t.PltaUnits {
		plta.PltaUnits = append(plta.PltaUnits, unit.ToResponse())
	}

	return plta
}

func (t *Plta) ToResponseMaster(userRadiusStatus bool) []*PltaResponseMaster {
	plta := []*PltaResponseMaster{}

	for _, unit := range t.PltaUnits {
		if !userRadiusStatus {
			t.RadiusStatus = false
		}

		plta = append(plta, &PltaResponseMaster{
			Id:           unit.Id,
			Name:         fmt.Sprintf("%v - Unit %v", t.Name, unit.Name),
			Lat:          t.Lat,
			Long:         t.Long,
			RadiusStatus: t.RadiusStatus,
			Radius:       t.Radius,
			RadiusType:   t.RadiusType,
		})
	}

	return plta
}

type PltaCreateRequest struct {
	Id           string
	Name         string         `json:"Name" form:"Name" validate:"required"`
	Status       bool           `json:"Status" form:"Status"`
	TotalUnits   uint32         `json:"TotalUnits" form:"TotalUnits" validate:"required"`
	Long         float64        `json:"Long" form:"Long" validate:"required"`
	Lat          float64        `json:"Lat" form:"Lat" validate:"required"`
	RadiusStatus bool           `json:"RadiusStatus" form:"RadiusStatus"`
	Radius       float64        `json:"Radius" form:"Radius" validate:"required"`
	RadiusType   PltaRadiusType `json:"RadiusType" form:"RadiusType" validate:"required,eq=meter|eq=kilometer"`
	WrittenBy    string
}

func (t *PltaCreateRequest) ToModelCreate() *Plta {
	plta := &Plta{
		Id:           ulid.Make().String(),
		Name:         t.Name,
		Status:       t.Status,
		Long:         t.Long,
		Lat:          t.Lat,
		RadiusStatus: t.RadiusStatus,
		Radius:       t.Radius,
		RadiusType:   t.RadiusType,
		CreatedBy:    t.WrittenBy,
		PltaUnits:    []*PltaUnit{},
	}

	for unit := uint32(1); unit <= t.TotalUnits; unit++ {
		plta.PltaUnits = append(plta.PltaUnits, &PltaUnit{
			Id:        ulid.Make().String(),
			PltaId:    plta.Id,
			Name:      fmt.Sprintf("%v", unit),
			CreatedBy: t.WrittenBy,
		})
	}

	return plta
}

type PltaUpdateRequest struct {
	Id           string
	Name         string         `json:"Name" form:"Name" validate:"required"`
	Status       bool           `json:"Status" form:"Status"`
	Long         float64        `json:"Long" form:"Long" validate:"required"`
	Lat          float64        `json:"Lat" form:"Lat" validate:"required"`
	RadiusStatus bool           `json:"RadiusStatus" form:"RadiusStatus"`
	Radius       float64        `json:"Radius" form:"Radius" validate:"required"`
	RadiusType   PltaRadiusType `json:"RadiusType" form:"RadiusType" validate:"required"`
	WrittenBy    string
}

func (t *PltaUpdateRequest) ToModelUpdate() *Plta {
	return &Plta{
		Id:           t.Id,
		Name:         t.Name,
		Status:       t.Status,
		Long:         t.Long,
		Lat:          t.Lat,
		RadiusStatus: t.RadiusStatus,
		Radius:       t.Radius,
		RadiusType:   t.RadiusType,
		UpdatedBy:    t.WrittenBy,
	}
}

type PltaResponse struct {
	Id           string              `json:"Id"`
	Name         string              `json:"Name"`
	Status       bool                `json:"Status"`
	Long         float64             `json:"Long"`
	Lat          float64             `json:"Lat"`
	RadiusStatus bool                `json:"RadiusStatus"`
	Radius       float64             `json:"Radius"`
	RadiusType   PltaRadiusType      `json:"RadiusType"`
	CreatedAt    string              `json:"CreatedAt"`
	CreatedBy    string              `json:"CreatedBy"`
	UpdatedAt    string              `json:"UpdatedAt"`
	UpdatedBy    string              `json:"UpdatedBy"`
	PltaUnits    []*PltaUnitResponse `json:"Units"`
}

type PltaResponseMaster struct {
	Id           string         `json:"Id"`
	Name         string         `json:"Name"`
	Long         float64        `json:"Long"`
	Lat          float64        `json:"Lat"`
	RadiusStatus bool           `json:"RadiusStatus"`
	Radius       float64        `json:"Radius"`
	RadiusType   PltaRadiusType `json:"RadiusType"`
}

type PltaListResponse struct {
	Id        string `json:"Id"`
	Name      string `json:"Name"`
	Status    bool   `json:"Status"`
	CreatedAt string `json:"CreatedAt"`
	CreatedBy string `json:"CreatedBy"`
}

func (t *Plta) ToResponseList() *PltaListResponse {
	plta := &PltaListResponse{
		Id:        t.Id,
		Name:      t.Name,
		Status:    t.Status,
		CreatedAt: t.CreatedAt.Format(helpers.DefaultTimeFormat),
	}

	if t.CreatedByUser != nil {
		plta.CreatedBy = t.CreatedByUser.Name
	}

	return plta
}

type PltaDeleteRequest struct {
	Id        string
	DeletedBy string
}

type PltaGetListMasterRequest struct {
	UserId string
	Search string
}

func (p *Plta) IsRadiusStatusActive() bool {
	return p.RadiusStatus
}

func (p *Plta) GetRadiusInKilometer() float64 {
	if p.RadiusType == PltaRadiusType_Meter {
		return p.Radius * 1000
	}

	return p.Radius
}

func (p *Plta) IsActive() bool {
	return p.Status
}
