package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/phuslu/log"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Turbine struct {
	Id                string          `gorm:"column:id"`
	GenBearingKopling float32         `gorm:"column:gen_bearing_kopling"`
	Koplingturbine    float32         `gorm:"column:kopling_turbine"`
	Total             float32         `gorm:"column:total"`
	Ratio             float32         `gorm:"column:ratio"`
	Data              datatypes.JSON  `gorm:"data"`
	CreatedAt         *time.Time      `gorm:"column:created_at"`
	UpdatedAt         *time.Time      `gorm:"column:updated_at;<-:update"`
	DeletedAt         *gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (t *Turbine) GetTurbineData() []TurbineData {
	var turbineData []TurbineData

	if err := json.Unmarshal(t.Data, &turbineData); err != nil {
		log.Error().Err(errors.New("ERROR UNMARSHAL TURBINE DATA : " + err.Error()))
		return []TurbineData{}
	}

	return turbineData
}

type Axis string

const (
	Axis_A       = "A"
	Axis_B       = "B"
	UppserAxis_C = "C"
	Axis_D       = "D"
)

type TurbineData struct {
	Axis Axis  `json:"Axis"`
	Test int32 `json:"Test"`
}

[
	{
		"Axis": "A",
		"Test": 1,
		"Value": 123
	},
	{
		"Axis": "A",
		"Test": 1,
		"Value": 123
	}
]

{
	"A": {
		"1": 23,
		"2": 32
	},
	"B": {
		"1": 23,
		"2": 32
	}
}

{
	"A1": 123,
	"A2": 123,
	"B1": 3123,
	"B2": 142
}

type TurbineWriteRequest struct {
	Id                string
	GenBearingKopling float32 `json:"GenBearingKopling"`
	Koplingturbine    float32 `json:"Koplingturbine"`
	Total             float32 `json:"Total"`
	Ratio             float32 `json:"Ratio"`

}
