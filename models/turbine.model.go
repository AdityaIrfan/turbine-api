package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/phuslu/log"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Turbine struct {
	Id                string          `gorm:"column:id"`
	TowerId           string          `gorm:"tower_id"`
	GenBearingKopling float32         `gorm:"column:gen_bearing_kopling"`
	Koplingturbine    float32         `gorm:"column:kopling_turbine"`
	Total             float32         `gorm:"column:total"`
	Ratio             float32         `gorm:"column:ratio"`
	Data              datatypes.JSON  `gorm:"data"`
	CreatedAt         *time.Time      `gorm:"column:created_at"`
	UpdatedAt         *time.Time      `gorm:"column:updated_at;<-:update"`
	DeletedAt         *gorm.DeletedAt `gorm:"column:deleted_at"`

	Tower *Tower `gorm:"foreignKey:TowerId;references:Id"`
}

type TurbineWriteRequest struct {
	Id                string
	TowerId           string                 `json:"TowerId" validate:"required"`
	GenBearingKopling float32                `json:"GenBearingKopling" validate:"required"`
	Koplingturbine    float32                `json:"Koplingturbine" validate:"required"`
	Total             float32                `json:"Total" validate:"required"`
	Ratio             float32                `json:"Ratio" validate:"required"`
	Data              map[string]interface{} `json:"Data" validate:"required"`
}

func (t *TurbineWriteRequest) ValidateData() error {
	if len(t.Data) < 1 {
		return errors.New("invalid data, data must be a json of the axes and their test value")
	}

	var availableAxes = map[string]bool{
		"A": false,
		"B": false,
		"C": false,
		"D": false,
	}

	var totalTest int
	var index int

	for axis, test := range t.Data {
		testValue, ok := test.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid value on axis %v, value must be a json", axis)
		}

		if len(testValue) == 0 {
			return fmt.Errorf("invalid value on axis %v, empty value", axis)
		}

		if index != 0 && totalTest != len(testValue) {
			return fmt.Errorf("invalid total data of each axis")
		}

		var missingTestvalue []string
		var indexTestValue = 1
		for test, value := range testValue {
			_, err := strconv.Atoi(test)
			if err != nil {
				return fmt.Errorf(`invalid key on axis %s, key '%s' is not accepted, key must be number of string`, axis, test)
			}
			if _, ok := value.(float64); !ok {
				return fmt.Errorf(`invalid value on axis %s, value of key '%v' must be a double`, axis, test)
			}
			v := strconv.Itoa(indexTestValue)
			if _, ok := testValue[v]; !ok {
				missingTestvalue = append(missingTestvalue, v)
			}
			indexTestValue++
		}

		if len(missingTestvalue) != 0 {
			return fmt.Errorf("missing [%v] test value on axis %v", strings.Join(missingTestvalue, ","), axis)
		}

		totalTest = len(testValue)

		if _, ok := availableAxes[axis]; ok {
			availableAxes[axis] = true
		}

		index++
	}

	var missingAxes []string
	for axis, isAvailable := range availableAxes {
		if !isAvailable {
			missingAxes = append(missingAxes, axis)
		}
	}

	if len(missingAxes) != 0 {
		return fmt.Errorf("missing axes [%v]", strings.Join(missingAxes, ", "))
	}

	return nil
}

func (t *TurbineWriteRequest) ToModelCreate() *Turbine {
	data, err := json.Marshal(t.Data)
	if err != nil {
		log.Error().Err(errors.New("ERROR MARSHAL TURBINE DATA : " + err.Error())).Msg("")
		data = datatypes.JSON{}
	}

	return &Turbine{
		GenBearingKopling: t.GenBearingKopling,
		Koplingturbine:    t.Koplingturbine,
		Total:             t.Total,
		Ratio:             t.Ratio,
		Data:              data,
	}
}
