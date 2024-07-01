package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"pln/AdityaIrfan/turbine-api/helpers"

	"github.com/oklog/ulid/v2"
	"github.com/phuslu/log"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var TurbineDefaultMap = map[string]string{
	"TowerName": "towername",
	"CreatedAt": "createdat",
}

type Turbine struct {
	Id                   string          `gorm:"column:id"`
	TowerId              string          `gorm:"tower_id"`
	GenBearingToCoupling float64         `gorm:"column:gen_bearing_to_coupling"`
	CouplingToTurbine    float64         `gorm:"column:coupling_to_turbine"`
	Data                 datatypes.JSON  `gorm:"data"`
	CreatedAt            *time.Time      `gorm:"column:created_at"`
	UpdatedAt            *time.Time      `gorm:"column:updated_at;<-:update"`
	DeletedAt            *gorm.DeletedAt `gorm:"column:deleted_at"`
	CreatedBy            string          `json:"column:created_by"`

	Tower *Tower `gorm:"foreignKey:TowerId;references:Id"`
	User  *User  `gorm:"foreignKey:CreatedBy;references:Id"`
}

func (t *Turbine) IsEmpty() bool {
	return t == nil
}

// type TurbineWriteRequest struct {
// 	Id                   string
// 	TowerId              string                 `json:"TowerId" form:"TowerId" validate:"required"`
// 	GenBearingToCoupling float64                `json:"GenBearingToCoupling" form:"GenBearingToCoupling" validate:"required"`
// 	CouplingToTurbine    float64                `json:"CouplingToTurbine" form:"CouplingToTurbine" validate:"required"`
// 	Data                 map[string]interface{} `json:"Data" form:"Data" validate:"required"`
// 	CreatedBy            string
// }

// func (t *TurbineWriteRequest) ValidateData() error {
// 	if len(t.Data) < 1 {
// 		return errors.New("data tidak valid, data harus berupa sumbu dan hasil percobaan dalam tipe data json")
// 	}

// 	var availableAxes = map[string]bool{
// 		"A": false,
// 		"B": false,
// 		"C": false,
// 		"D": false,
// 	}

// 	var availableParts = map[string]bool{
// 		"Upper":   false,
// 		"Clutch":  false,
// 		"Turbine": false,
// 	}

// 	for part, partValue := range t.Data {

// 		dataPerPart, ok := partValue.(map[string]interface{})
// 		if !ok {
// 			return fmt.Errorf("data dibagian %v harus betipe data json", part)
// 		}
// 		if _, ok := availableParts[part]; ok {
// 			availableParts[part] = true
// 		}

// 		var index int
// 		var totalTest int
// 		for axis, test := range dataPerPart {
// 			testValue, ok := test.(map[string]interface{})
// 			if !ok {
// 				return fmt.Errorf("data sumbu %v pada bagian %v tidak valid, data harus bertipa data json", axis, part)
// 			}

// 			if len(testValue) == 0 {
// 				return fmt.Errorf("data sumbu %v pada bagian %v tidak valid, nilai kosong", axis, part)
// 			}

// 			if index != 0 && totalTest != len(testValue) {
// 				return fmt.Errorf("data tidak valid, total percobaan pada masing masing sumbu bagian %v tidak sama", part)
// 			}

// 			var missingTestvalue []string
// 			var indexTestValue = 1
// 			for test, value := range testValue {
// 				_, err := strconv.Atoi(test)
// 				if err != nil {
// 					return fmt.Errorf(`percobaan pada sumbu %s bagian %v tidak valid, percobaan '%s' tidak diterima, percobaan harus berupa nomor string`, axis, part, test)
// 				}
// 				if _, ok := value.(float64); !ok {
// 					return fmt.Errorf(`hasil percobaan pada sumbu %s bagian %v tidak valid, hasil pada percobaan '%v' harus berupa desimal`, axis, part, test)
// 				}
// 				v := strconv.Itoa(indexTestValue)
// 				if _, ok := testValue[v]; !ok {
// 					missingTestvalue = append(missingTestvalue, v)
// 				}
// 				indexTestValue++
// 			}

// 			if len(missingTestvalue) != 0 {
// 				return fmt.Errorf("percobaan ke-%v tidak ditemukan pada sumbu %v bagian %v", strings.Join(missingTestvalue, ","), axis, part)
// 			}

// 			totalTest = len(testValue)

// 			if _, ok := availableAxes[axis]; ok {
// 				availableAxes[axis] = true
// 			}

// 			index++
// 		}
// 	}

// 	var missingParts []string
// 	for axis, isAvailable := range availableParts {
// 		if !isAvailable {
// 			missingParts = append(missingParts, axis)
// 		}
// 	}
// 	if len(missingParts) != 0 {
// 		return fmt.Errorf("bagian %v tidak ditemukan", strings.Join(missingParts, ", "))
// 	}

// 	var missingAxes []string
// 	for axis, isAvailable := range availableAxes {
// 		if !isAvailable {
// 			missingAxes = append(missingAxes, axis)
// 		}
// 	}
// 	if len(missingAxes) != 0 {
// 		return fmt.Errorf("sumbu %v tidak ditemukan", strings.Join(missingAxes, ", "))
// 	}

// 	return nil
// }

type TurbineWriteRequest struct {
	Id                   string
	TowerId              string                 `json:"TowerId" form:"TowerId" validate:"required"`
	GenBearingToCoupling float64                `json:"GenBearingToCoupling" form:"GenBearingToCoupling" validate:"required"`
	CouplingToTurbine    float64                `json:"CouplingToTurbine" form:"CouplingToTurbine" validate:"required"`
	Data                 map[string]interface{} `json:"Data" form:"Data" validate:"required"`
	CreatedBy            string
}

func (t *TurbineWriteRequest) ValidateData() error {
	if len(t.Data) < 1 {
		return errors.New("data tidak valid, data harus berupa sumbu dan hasil percobaan dalam tipe data json")
	}

	var availableAxes = map[string]bool{
		"A": false,
		"B": false,
		"C": false,
		"D": false,
	}

	var availableParts = map[string]bool{
		"Upper":   false,
		"Clutch":  false,
		"Turbine": false,
	}

	for part, partValue := range t.Data {

		dataPerPart, ok := partValue.(map[string]interface{})
		if !ok {
			return fmt.Errorf("data dibagian %v harus betipe json", part)
		}
		if _, ok := availableParts[part]; ok {
			availableParts[part] = true
		}

		var index int
		var totalTest int
		for axis, test := range dataPerPart {
			testValue, ok := test.([]interface{})
			if !ok {
				return fmt.Errorf("data sumbu %v pada bagian %v tidak valid, data harus bertipa array number", axis, part)
			}

			if len(testValue) == 0 {
				return fmt.Errorf("data sumbu %v pada bagian %v tidak valid, nilai kosong", axis, part)
			}

			if index != 0 && totalTest != len(testValue) {
				return fmt.Errorf("data tidak valid, total percobaan pada masing masing sumbu bagian %v tidak sama", part)
			}

			totalTest = len(testValue)

			if _, ok := availableAxes[axis]; ok {
				availableAxes[axis] = true
			}

			index++

			testValueTemp := make(map[string]interface{})
			for i := 0; i < len(testValue); i++ {
				value, ok := testValue[i].(float64)
				if !ok {
					return fmt.Errorf(`hasil percobaan pada sumbu %s bagian %v tidak valid, hasil pada percobaan '%v' harus berupa desimal`, axis, part, i+1)
				}
				testValueTemp[fmt.Sprintf("%v", i+1)] = value
			}
			t.Data[part].(map[string]interface{})[axis] = testValueTemp
		}
	}

	var missingParts []string
	for axis, isAvailable := range availableParts {
		if !isAvailable {
			missingParts = append(missingParts, axis)
		}
	}
	if len(missingParts) != 0 {
		return fmt.Errorf("bagian %v tidak ditemukan", strings.Join(missingParts, ", "))
	}

	var missingAxes []string
	for axis, isAvailable := range availableAxes {
		if !isAvailable {
			missingAxes = append(missingAxes, axis)
		}
	}
	if len(missingAxes) != 0 {
		return fmt.Errorf("sumbu %v tidak ditemukan", strings.Join(missingAxes, ", "))
	}

	return nil
}

// func (t *TurbineWriteRequest) ToModelCreate() *Turbine {
// 	data, err := json.Marshal(t.Data)
// 	if err != nil {
// 		log.Error().Err(errors.New("ERROR MARSHAL TURBINE DATA : " + err.Error())).Msg("")
// 		data = datatypes.JSON{}
// 	}

// 	return &Turbine{
// 		Id:                   ulid.Make().String(),
// 		TowerId:              t.TowerId,
// 		GenBearingToCoupling: t.GenBearingToCoupling,
// 		CouplingToTurbine:    t.CouplingToTurbine,
// 		Data:                 data,
// 		CreatedBy:            t.CreatedBy,
// 	}
// }

func (t *TurbineWriteRequest) ToModelCreate() *Turbine {

	data, err := json.Marshal(t.Data)
	if err != nil {
		log.Error().Err(errors.New("ERROR MARSHAL TURBINE DATA : " + err.Error())).Msg("")
		data = datatypes.JSON{}
	}

	return &Turbine{
		Id:                   ulid.Make().String(),
		TowerId:              t.TowerId,
		GenBearingToCoupling: t.GenBearingToCoupling,
		CouplingToTurbine:    t.CouplingToTurbine,
		Data:                 data,
		CreatedBy:            t.CreatedBy,
	}
}

type TurbineResponse struct {
	Id               string                 `json:"Id"`
	TowerName        string                 `json:"TowerName"`
	Shaft            TurbineShaft           `json:"Shaft"`
	Chart            map[string]interface{} `json:"Chart"`
	DetailData       map[string]interface{} `json:"DetailData"`
	ACCrockedness    float64                `json:"ACCrockedness"`
	BDCrockedness    float64                `json:"BDCrockedness"`
	TotalCrockedness float64                `json:"TotalCrockedness"`
	CreatedAt        string                 `json:"CreatedAt"`
	CreatedBy        string                 `json:"CreatedBy"`
}

type TurbineShaft struct {
	GenBearingToKopling float64 `json:"GenBearingToCoupling"`
	CouplingToTurbine   float64 `json:"CouplingToTurbine"`
	Total               float64 `json:"Total"`
	Ratio               float64 `json:"Ratio"`
}

func (t *Turbine) ToResponse() *TurbineResponse {
	total := t.GenBearingToCoupling + t.CouplingToTurbine
	ratio := t.GenBearingToCoupling / total
	chart := make(map[string]interface{})
	chart["AC"] = map[string]interface{}{
		"Upper":   "",
		"Clutch":  "",
		"Turbine": "",
	}
	chart["BD"] = map[string]interface{}{
		"Upper":   "",
		"Clutch":  "",
		"Turbine": "",
	}
	chart["Upper"] = ""

	var detailData map[string]interface{}
	json.Unmarshal(t.Data, &detailData)

	averageData := map[string]map[string]float64{
		"Upper": {
			"A": 0,
			"B": 0,
			"C": 0,
			"D": 0,
		},
		"Clutch": {
			"A": 0,
			"B": 0,
			"C": 0,
			"D": 0,
		},
		"Turbine": {
			"A": 0,
			"B": 0,
			"C": 0,
			"D": 0,
		},
	}

	// Get upper average
	upper := detailData["Upper"].(map[string]interface{})
	upperTotaltest := len(upper["A"].(map[string]interface{}))
	var upperTotalA, upperTotalB, upperTotalC, upperTotalD float64
	for i := 1; i <= upperTotaltest; i++ {
		upperTotalA += upper["A"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		upperTotalB += upper["B"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		upperTotalC += upper["C"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		upperTotalD += upper["D"].(map[string]interface{})[strconv.Itoa(i)].(float64)
	}
	averageData["Upper"]["A"] = upperTotalA / float64(upperTotaltest)
	averageData["Upper"]["B"] = upperTotalB / float64(upperTotaltest)
	averageData["Upper"]["C"] = upperTotalC / float64(upperTotaltest)
	averageData["Upper"]["D"] = upperTotalD / float64(upperTotaltest)

	// Get clutch average
	clutch := detailData["Clutch"].(map[string]interface{})
	clutchTotaltest := len(clutch["A"].(map[string]interface{}))
	var clutchTotalA, clutchTotalB, clutchTotalC, clutchTotalD float64
	for i := 1; i <= clutchTotaltest; i++ {
		clutchTotalA += clutch["A"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		clutchTotalB += clutch["B"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		clutchTotalC += clutch["C"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		clutchTotalD += clutch["D"].(map[string]interface{})[strconv.Itoa(i)].(float64)
	}
	averageData["Clutch"]["A"] = clutchTotalA / float64(clutchTotaltest)
	averageData["Clutch"]["B"] = clutchTotalB / float64(clutchTotaltest)
	averageData["Clutch"]["C"] = clutchTotalC / float64(clutchTotaltest)
	averageData["Clutch"]["D"] = clutchTotalD / float64(clutchTotaltest)

	// Get turbine average
	turbine := detailData["Turbine"].(map[string]interface{})
	turbineTotaltest := len(turbine["A"].(map[string]interface{}))
	var turbineTotalA, turbineTotalB, turbineTotalC, turbineTotalD float64
	for i := 1; i <= turbineTotaltest; i++ {
		turbineTotalA += turbine["A"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		turbineTotalB += turbine["B"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		turbineTotalC += turbine["C"].(map[string]interface{})[strconv.Itoa(i)].(float64)
		turbineTotalD += turbine["D"].(map[string]interface{})[strconv.Itoa(i)].(float64)
	}
	averageData["Turbine"]["A"] = turbineTotalA / float64(turbineTotaltest)
	averageData["Turbine"]["B"] = turbineTotalB / float64(turbineTotaltest)
	averageData["Turbine"]["C"] = turbineTotalC / float64(turbineTotaltest)
	averageData["Turbine"]["D"] = turbineTotalD / float64(turbineTotaltest)

	averageUpperAC := averageData["Upper"]["C"] - (averageData["Upper"]["A"])
	averageClutchAC := averageData["Clutch"]["C"] - (averageData["Clutch"]["A"])
	averageTurbineAC := averageData["Turbine"]["C"] - (averageData["Turbine"]["A"])

	averageUpperBD := averageData["Upper"]["D"] - (averageData["Upper"]["B"])
	averageClutchBD := averageData["Clutch"]["D"] - (averageData["Clutch"]["B"])
	averageTurbineBD := averageData["Turbine"]["D"] - (averageData["Turbine"]["B"])

	crockednessAC := math.Pow(float64(ratio*(averageTurbineAC-averageUpperAC))-(averageClutchAC-averageUpperAC), 2)
	crockednessBD := math.Pow(float64(ratio*(averageTurbineBD-averageUpperBD))-(averageClutchBD-averageUpperBD), 2)

	defaultUpperAC := averageUpperAC
	if defaultUpperAC > 0 {
		defaultUpperAC = defaultUpperAC * -1
	}
	chartClutchAC := averageClutchAC + defaultUpperAC
	chartTurbineAC := averageTurbineAC + defaultUpperAC
	chart["AC"].(map[string]interface{})["Upper"] = fmt.Sprintf("0|%f", t.GenBearingToCoupling)
	chart["AC"].(map[string]interface{})["Clutch"] = fmt.Sprintf("%f|0", chartClutchAC)
	chart["AC"].(map[string]interface{})["Turbine"] = fmt.Sprintf("%f|%f", chartTurbineAC, t.CouplingToTurbine)

	defaultUpperBD := averageUpperBD
	if defaultUpperBD > 0 {
		defaultUpperBD = defaultUpperBD * -1
	}
	chartClutchBD := averageClutchBD + defaultUpperBD
	chartTurbineBD := averageTurbineBD + defaultUpperBD
	chart["BD"].(map[string]interface{})["Upper"] = fmt.Sprintf("0|%f", t.GenBearingToCoupling)
	chart["BD"].(map[string]interface{})["Clutch"] = fmt.Sprintf("%f|0", chartClutchBD)
	chart["BD"].(map[string]interface{})["Turbine"] = fmt.Sprintf("%f|%f", chartTurbineBD, t.CouplingToTurbine)

	chart["Upper"] = fmt.Sprintf("%f|%f", crockednessAC, crockednessBD)

	return &TurbineResponse{
		Id:        t.Id,
		TowerName: fmt.Sprintf("%v - %v", t.Tower.Name, t.Tower.UnitNumber),
		Shaft: TurbineShaft{
			GenBearingToKopling: t.GenBearingToCoupling,
			CouplingToTurbine:   t.CouplingToTurbine,
			Total:               total,
			Ratio:               ratio,
		},
		Chart:            chart,
		DetailData:       detailData,
		ACCrockedness:    crockednessAC,
		BDCrockedness:    crockednessBD,
		TotalCrockedness: math.Pow((crockednessAC + crockednessBD), 0.5),
		CreatedAt:        t.CreatedAt.Format(helpers.DefaultTimeFormat),
		CreatedBy:        t.User.Name,
	}
}

type TurbineResponseList struct {
	Id        string `json:"Id"`
	TowerName string `json:"TowerName"`
	CreatedAt string `json:"CreatedAt"`
}

func (t *Turbine) ToResponseList() *TurbineResponseList {
	return &TurbineResponseList{
		Id:        t.Id,
		TowerName: fmt.Sprintf("%v - %v", t.Tower.Name, t.Tower.UnitNumber),
		CreatedAt: t.CreatedAt.Format(helpers.DefaultTimeFormat),
	}
}
