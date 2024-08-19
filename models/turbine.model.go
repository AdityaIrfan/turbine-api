package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"pln/AdityaIrfan/turbine-api/helpers"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/oklog/ulid/v2"
	"github.com/phuslu/log"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var TurbineDefaultSortMap = map[string]string{
	"Title":     "title",
	"TowerName": "towername",
	"CreatedAt": "createdat",
}

var TurbineDefaultFilter = map[string]string{
	"Tower": "tower_id",
}

const TotalCrockednessToleration = 3

type Turbine struct {
	Id                   string          `gorm:"column:id"`
	Title                string          `gormc:"column:title"`
	PltaUnitId           string          `gorm:"column:plta_unit_id"`
	GenBearingToCoupling float64         `gorm:"column:gen_bearing_to_coupling"`
	CouplingToTurbine    float64         `gorm:"column:coupling_to_turbine"`
	Data                 datatypes.JSON  `gorm:"column:data"`
	CreatedAt            *time.Time      `gorm:"column:created_at"`
	CreatedBy            string          `gorm:"column:created_by"`
	DeletedAt            *gorm.DeletedAt `gorm:"column:deleted_at"`
	DeletedBy            string          `gorm:"deleted_by"`
	TotalBolts           uint32          `gorm:"column:total_bolts"`
	CurrentTorque        float64         `gorm:"column:current_torque"`
	MaxTorque            float64         `gorm:"column:max_torque"`

	PltaUnit *PltaUnit `gorm:"foreignKey:PltaUnitId;references:Id"`
}

func (t *Turbine) TableName() string {
	return "turbines"
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
	Title                string                 `json:"Title" form:"Title" validate:"required"`
	GenBearingToCoupling float64                `json:"GenBearingToCoupling" form:"GenBearingToCoupling" validate:"required"`
	CouplingToTurbine    float64                `json:"CouplingToTurbine" form:"CouplingToTurbine" validate:"required"`
	Data                 map[string]interface{} `json:"Data" form:"Data" validate:"required"`
	TotalBolts           uint32                 `json:"TotalBolts" form:"TotalBolts" validate:"required,min=4,max=24"`
	CurrentTorque        float64                `json:"CurrentTorque" form:"CurrentTorque" validate:"required"`
	MaxTorque            float64                `json:"MaxTorque" form:"MaxTorque" validate:"required"`
	WrittenBy            string
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
		Title:                t.Title,
		PltaUnitId:           t.TowerId,
		GenBearingToCoupling: t.GenBearingToCoupling,
		CouplingToTurbine:    t.CouplingToTurbine,
		Data:                 data,
		CreatedBy:            t.WrittenBy,
		TotalBolts:           t.TotalBolts,
		CurrentTorque:        t.CurrentTorque,
		MaxTorque:            t.MaxTorque,
	}
}

type TurbineResponse struct {
	Id                string                 `json:"Id"`
	Title             string                 `json:"Title"`
	TowerName         string                 `json:"TowerName"`
	Shaft             TurbineShaft           `json:"Shaft"`
	Chart             map[string]interface{} `json:"Chart"`
	DetailData        map[string]interface{} `json:"DetailData"`
	ACCrockedness     float64                `json:"ACCrockedness"`
	BDCrockedness     float64                `json:"BDCrockedness"`
	TotalCrockedness  float64                `json:"TotalCrockedness"`
	CreatedAt         string                 `json:"CreatedAt"`
	CreatedBy         string                 `json:"CreatedBy"`
	Status            bool                   `json:"Status"`
	TotalBolts        uint32                 `json:"TotalBolts"`
	CurrentTorque     float64                `json:"CurrentTorque"`
	MaxTorque         float64                `json:"MaxTorque"`
	TorqueGap         float64                `json:"TorqueGap"`
	TorqueCalculation map[string]interface{} `json:"TorqueCalculation"`
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

	resultanAC := float64(ratio*(averageTurbineAC-averageUpperAC)) - (averageClutchAC - averageUpperAC)
	resultanBD := float64(ratio*(averageTurbineBD-averageUpperBD)) - (averageClutchBD - averageUpperBD)
	crockednessAC := math.Pow(resultanAC, 2)
	crockednessBD := math.Pow(resultanBD, 2)

	defaultUpperAC := averageUpperAC
	chartClutchAC := averageClutchAC
	chartTurbineAC := averageTurbineAC
	if chartClutchAC > 0 {
		if defaultUpperAC > 0 {
			chartClutchAC -= defaultUpperAC
		} else {
			chartClutchAC += defaultUpperAC * -1
		}
	} else {
		chartClutchAC -= defaultUpperAC
	}
	if chartTurbineAC > 0 {
		if defaultUpperAC > 0 {
			chartTurbineAC -= defaultUpperAC
		} else {
			chartTurbineAC += defaultUpperAC * -1
		}
	} else {
		chartTurbineAC -= defaultUpperAC
	}
	chart["AC"].(map[string]interface{})["Upper"] = fmt.Sprintf("0|%f", t.GenBearingToCoupling)
	chart["AC"].(map[string]interface{})["Clutch"] = fmt.Sprintf("%f|0", chartClutchAC)
	chart["AC"].(map[string]interface{})["Turbine"] = fmt.Sprintf("%f|%f", chartTurbineAC, t.CouplingToTurbine)
	chart["AC"].(map[string]interface{})["Scale"] = GetScale(t.GenBearingToCoupling, t.CouplingToTurbine, chartClutchAC, chartTurbineAC)

	defaultUpperBD := averageUpperBD
	chartClutchBD := averageClutchBD
	chartTurbineBD := averageTurbineBD
	if chartClutchBD > 0 {
		if defaultUpperBD > 0 {
			chartClutchBD -= defaultUpperBD
		} else {
			chartClutchBD += defaultUpperBD
		}
	} else {
		chartClutchBD -= defaultUpperBD
	}
	if chartTurbineBD > 0 {
		if defaultUpperBD > 0 {
			chartTurbineBD -= defaultUpperBD
		} else {
			chartTurbineBD += defaultUpperBD
		}
	} else {
		chartTurbineBD -= defaultUpperBD
	}
	chart["BD"].(map[string]interface{})["Upper"] = fmt.Sprintf("0|%f", t.GenBearingToCoupling)
	chart["BD"].(map[string]interface{})["Clutch"] = fmt.Sprintf("%f|0", chartClutchBD)
	chart["BD"].(map[string]interface{})["Turbine"] = fmt.Sprintf("%f|%f", chartTurbineBD, t.CouplingToTurbine)
	chart["BD"].(map[string]interface{})["Scale"] = GetScale(t.GenBearingToCoupling, t.CouplingToTurbine, chartClutchBD, chartTurbineBD)

	chart["Upper"] = fmt.Sprintf("%f|%f", resultanAC, resultanBD)
	totalCrockedness := math.Pow((crockednessAC + crockednessBD), 0.5)
	chart["UpperScale"] = 7
	if totalCrockedness != 0 {
		chart["UpperScale"] = GetScale(resultanAC, resultanBD)
	}

	// TORQUE CALCULATION
	torqueGap := t.MaxTorque - t.CurrentTorque
	totalAngleInDegrees := uint32(360)
	degreeGap := float64(totalAngleInDegrees) / float64(t.TotalBolts)
	circleRadius := math.Sqrt(math.Pow(resultanAC, 2) + math.Pow(resultanBD, 2))
	noCrockedness := false
	if circleRadius == 0 {
		noCrockedness = true
		circleRadius += 7
	}
	bolt := 1
	TorqueCalculation := make(map[string]interface{})
	TorqueCalculationDetail := make(map[string]interface{})
	points := [][2]float64{}
	pointsTemp := make(map[string]uint32)
	// fmt.Println("RADIUS : ", circleRadius)
	for i := float64(totalAngleInDegrees); i > 0; i -= degreeGap {
		angleInDegrees := float64(totalAngleInDegrees) - i
		angleInRadians := angleInDegrees * math.Pi / 180.0

		// Compute the cosine and sin of degrees
		cosValue := math.Cos(angleInRadians)
		sinValue := math.Sin(angleInRadians)

		// fmt.Printf("cos(%f degrees) = %v\n", angleInDegrees, cosValue)
		// fmt.Printf("sin(%f degrees) = %v\n", angleInDegrees, sinValue)
		// fmt.Println()

		if math.Abs(cosValue) < 1e-9 { // Allowing a small tolerance
			cosValue = 0
		}
		if math.Abs(sinValue) < 1e-9 { // Allowing a small tolerance
			sinValue = 0
		}

		x := circleRadius * cosValue
		y := circleRadius * sinValue
		points = append(points, [2]float64{x, y})
		pointsTemp[fmt.Sprintf("%f|%f", x, y)] = uint32(bolt)
		TorqueCalculationDetail[fmt.Sprintf("%d", bolt)] = fmt.Sprintf("%f|%f", x, y)
		if bolt == 1 {
			bolt = int(t.TotalBolts)
		} else {
			bolt--
		}
	}

	TorqueCalculation["Details"] = TorqueCalculationDetail
	TorqueSuggestion := make(map[string]float64)

	// TORQUE SUGGESTION
	if !noCrockedness {
		// closestCoordinates := getClosestCoordinates(points, 0.0, -circleRadius)
		closestCoordinates := getClosestCoordinates(points, resultanAC, resultanBD)
		if len(closestCoordinates) != 0 {
			if len(closestCoordinates) == 1 {
				orderBolt, ok := pointsTemp[closestCoordinates[0]]
				if ok {
					TorqueSuggestion[fmt.Sprintf("%d", orderBolt)] = math.Round(t.CurrentTorque + (0.5 * torqueGap))
					prevOrderBolt := orderBolt - 1
					nextOrderBolt := orderBolt + 1

					if orderBolt == 1 {
						prevOrderBolt = t.TotalBolts
					}
					if orderBolt == t.TotalBolts {
						nextOrderBolt = 1
					}

					TorqueSuggestion[fmt.Sprintf("%d", prevOrderBolt)] = math.Round(t.CurrentTorque + (0.25 * torqueGap))
					TorqueSuggestion[fmt.Sprintf("%d", nextOrderBolt)] = math.Round(t.CurrentTorque + (0.25 * torqueGap))
				} else {
					log.Error().Err(fmt.Errorf("COORDINATE %s IS NOT EXIST", closestCoordinates[0]))
				}
			} else {
				var pressureOnTheClossestBolt float64
				closestBolt, ok := pointsTemp[closestCoordinates[0]]
				if ok {
					// The closes bolt have the main formula no matter what the distance
					// Closest Bolt Torque = current torque + (max torque - current torque) * 0.5
					pressureOnTheClossestBolt = t.CurrentTorque + (0.5 * torqueGap)
					TorqueSuggestion[fmt.Sprintf("%d", closestBolt)] = math.Round(pressureOnTheClossestBolt)
				} else {
					log.Error().Err(fmt.Errorf("COORDINATE %s IS NOT EXIST", closestCoordinates[0]))
				}

				// Find the angle in degrees between coordinate (resultanAC, resultanBD) and the coordinate of the closest bolt
				// Split closest coordinate from resultan from string type containing "|" to be x and y
				closestBoltPoint := strings.Split(closestCoordinates[0], "|")
				closesBoltPointX, _ := strconv.ParseFloat(closestBoltPoint[0], 64)
				closesBoltPointY, _ := strconv.ParseFloat(closestBoltPoint[1], 64)
				// set (0,0) as a center of the circle
				closestDegree := calculateAngleInDegrees(Point{0, 0}, Point{X: resultanAC, Y: resultanBD}, Point{X: closesBoltPointX, Y: closesBoltPointY})

				// Find the angle in degree every closest point from the second to fourth closes
				// Second Bolt
				SecondClosestBolt := strings.Split(closestCoordinates[1], "|")
				secondClosesBoltPointX, _ := strconv.ParseFloat(SecondClosestBolt[0], 64)
				secondClosesBoltPointY, _ := strconv.ParseFloat(SecondClosestBolt[1], 64)
				secondClosestDegree := calculateAngleInDegrees(Point{0, 0}, Point{X: resultanAC, Y: resultanBD}, Point{X: secondClosesBoltPointX, Y: secondClosesBoltPointY})
				// Third Bolt
				ThirdClosestBolt := strings.Split(closestCoordinates[2], "|")
				thirdClosesBoltPointX, _ := strconv.ParseFloat(ThirdClosestBolt[0], 64)
				thirdClosesBoltPointY, _ := strconv.ParseFloat(ThirdClosestBolt[1], 64)
				thirdClosestDegree := calculateAngleInDegrees(Point{0, 0}, Point{X: resultanAC, Y: resultanBD}, Point{X: thirdClosesBoltPointX, Y: thirdClosesBoltPointY})
				// Fourth Bolt
				FourthClosestBolt := strings.Split(closestCoordinates[3], "|")
				fourthClosesBoltPointX, _ := strconv.ParseFloat(FourthClosestBolt[0], 64)
				fourthClosesBoltPointY, _ := strconv.ParseFloat(FourthClosestBolt[1], 64)
				fourthClosestDegree := calculateAngleInDegrees(Point{0, 0}, Point{X: resultanAC, Y: resultanBD}, Point{X: fourthClosesBoltPointX, Y: fourthClosesBoltPointY})

				degreesPerBolt := totalAngleInDegrees / t.TotalBolts

				// Second Bolt
				secondtBolt, ok := pointsTemp[closestCoordinates[1]]
				if ok {
					// Except the closest bolt, there is fixed formula for others
					// Bolt Torque = current torque + (closest bolt degree / degrees per bolt) * ((pressure closest bolt - current torque) * degrees per bolt / degrees each bolt from resultan coordinate
					TorqueSuggestion[fmt.Sprintf("%d", secondtBolt)] = math.Round(t.CurrentTorque + ((closestDegree / float64(degreesPerBolt)) * ((pressureOnTheClossestBolt - t.CurrentTorque) * float64(degreesPerBolt) / secondClosestDegree)))
				} else {
					log.Error().Err(fmt.Errorf("COORDINATE %s IS NOT EXIST", closestCoordinates[1]))
				}

				// Second Bolt
				thirdBolt, ok := pointsTemp[closestCoordinates[2]]
				if ok {
					// Except the closest bolt, there is fixed formula for others
					// Bolt Torque = current torque + (closest bolt degree / degrees per bolt) * ((pressure closest bolt - current torque) * degrees per bolt / degrees each bolt from resultan coordinate
					TorqueSuggestion[fmt.Sprintf("%d", thirdBolt)] = math.Round(t.CurrentTorque + ((closestDegree / float64(degreesPerBolt)) * ((pressureOnTheClossestBolt - t.CurrentTorque) * float64(degreesPerBolt) / thirdClosestDegree)))
				} else {
					log.Error().Err(fmt.Errorf("COORDINATE %s IS NOT EXIST", closestCoordinates[2]))
				}

				// Second Bolt
				fourthtBolt, ok := pointsTemp[closestCoordinates[3]]
				if ok {
					// Except the closest bolt, there is fixed formula for others
					// Bolt Torque = current torque + (closest bolt degree / degrees per bolt) * ((pressure closest bolt - current torque) * degrees per bolt / degrees each bolt from resultan coordinate
					TorqueSuggestion[fmt.Sprintf("%d", fourthtBolt)] = math.Round(t.CurrentTorque + ((closestDegree / float64(degreesPerBolt)) * ((pressureOnTheClossestBolt - t.CurrentTorque) * float64(degreesPerBolt) / fourthClosestDegree)))
				} else {
					log.Error().Err(fmt.Errorf("COORDINATE %s IS NOT EXIST", closestCoordinates[3]))
				}
			}
		}
	}

	TorqueCalculation["TorqueSuggestions"] = TorqueSuggestion
	TorqueCalculation["Scale"] = circleRadius + 1

	return &TurbineResponse{
		Id:        t.Id,
		Title:     t.Title,
		TowerName: fmt.Sprintf("%v - Unit %v", t.PltaUnit.Plta.Name, t.PltaUnit.Name),
		Shaft: TurbineShaft{
			GenBearingToKopling: t.GenBearingToCoupling,
			CouplingToTurbine:   t.CouplingToTurbine,
			Total:               total,
			Ratio:               ratio,
		},
		Chart:             chart,
		DetailData:        detailData,
		ACCrockedness:     crockednessAC,
		BDCrockedness:     crockednessBD,
		TotalCrockedness:  totalCrockedness,
		CreatedAt:         t.CreatedAt.Format(helpers.DefaultTimeFormat),
		CreatedBy:         t.CreatedBy,
		Status:            totalCrockedness <= TotalCrockednessToleration,
		TotalBolts:        t.TotalBolts,
		CurrentTorque:     t.CurrentTorque,
		MaxTorque:         t.MaxTorque,
		TorqueGap:         torqueGap,
		TorqueCalculation: TorqueCalculation,
	}
}

type TurbineResponseList struct {
	Id        string `json:"Id"`
	Title     string `json:"Title"`
	TowerName string `json:"TowerName"`
	CreatedAt string `json:"CreatedAt"`
}

func (t *Turbine) ToResponseList() *TurbineResponseList {
	return &TurbineResponseList{
		Id:        t.Id,
		Title:     t.Title,
		TowerName: fmt.Sprintf("%v - Unit %v", t.PltaUnit.Plta.Name, t.PltaUnit.Name),
		CreatedAt: t.CreatedAt.Format(helpers.DefaultTimeFormat),
	}
}

func getClosestCoordinates(points [][2]float64, targetX, targetY float64) []string {
	// Slice to hold the points and their distances
	var pointDistances []PointDistance

	// Calculate the distance for each point
	for _, point := range points {
		if point[0] == targetX && point[1] == targetY {
			return []string{fmt.Sprintf("%f|%f", targetX, targetY)}
		}
		d := distance(targetX, targetY, point[0], point[1])
		pointDistances = append(pointDistances, PointDistance{Point: point, Distance: d})
	}

	// Sort the points by distance
	sort.Slice(pointDistances, func(i, j int) bool {
		return pointDistances[i].Distance < pointDistances[j].Distance
	})

	var coordinateString = []string{}
	for _, d := range pointDistances {
		coordinateString = append(coordinateString, fmt.Sprintf("%f|%f", d.Point[0], d.Point[1]))
		fmt.Printf("(%v, %v) with a distance of %v\n", d.Point[0], d.Point[1], d.Distance)
	}

	return coordinateString
}

// Function to calculate the distance between two points
func distance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}

// Struct to hold a point and its distance from (5,5)
type PointDistance struct {
	Point    [2]float64
	Distance float64
}

// Point represents a point in 2D space
type Point struct {
	X, Y float64
}

// calculateAngle calculates the angle between two points on a circle centered at the origin
func calculateAngleInDegrees(center, p1, p2 Point) float64 {
	// Calculate vectors A and B
	vectorA := Point{X: p1.X - center.X, Y: p1.Y - center.Y}
	vectorB := Point{X: p2.X - center.X, Y: p2.Y - center.Y}

	// Calculate the dot product of vectors A and B
	dotProduct := vectorA.X*vectorB.X + vectorA.Y*vectorB.Y

	// Calculate the magnitudes of vectors A and B
	magnitudeA := math.Sqrt(vectorA.X*vectorA.X + vectorA.Y*vectorA.Y)
	magnitudeB := math.Sqrt(vectorB.X*vectorB.X + vectorB.Y*vectorB.Y)

	// Calculate the cosine of the angle
	cosTheta := dotProduct / (magnitudeA * magnitudeB)

	// Calculate the angle in radians
	theta := math.Acos(cosTheta)

	// Convert the angle from radians to degrees
	thetaDegrees := theta * (180 / math.Pi)

	return thetaDegrees
}

func GetScale(data ...float64) float64 {
	switch len(data) {
	case 0:
		return 1
	case 1:
		return data[0] + 1
	}

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data)-i-1; j++ {
			if data[j] < 0 {
				data[j] *= -1
			}

			if data[j+1] < 0 {
				data[j+1] *= -1
			}

			if data[j] > data[j+1] {
				data[j], data[j+1] = data[j+1], data[j]
			}
		}
	}

	return data[len(data)-1] + 1
}

func (t *Turbine) GenerateReport() ([]byte, error) {
	turbineResponse := t.ToResponse()

	// upper
	upper := [][]string{}
	upperInit := false
	for _, u := range turbineResponse.DetailData["Upper"].(map[string]interface{}) {
		indexColumn := 0
		for _, data := range u.(map[string]interface{}) {
			if !upperInit {
				upper = append(upper, []string{})
			}

			upper[indexColumn] = append(upper[indexColumn], fmt.Sprintf("%v", data))
			indexColumn++
		}

		upperInit = true
	}

	// clutch
	clutch := [][]string{}
	clutchInit := false
	for _, u := range turbineResponse.DetailData["Clutch"].(map[string]interface{}) {
		indexColumn := 0
		for _, data := range u.(map[string]interface{}) {
			if !clutchInit {
				clutch = append(clutch, []string{})
			}

			clutch[indexColumn] = append(clutch[indexColumn], fmt.Sprintf("%v", data))
			indexColumn++
		}

		clutchInit = true
	}

	// turbine
	turbine := [][]string{}
	turbineInit := false
	for _, u := range turbineResponse.DetailData["Turbine"].(map[string]interface{}) {
		indexColumn := 0
		for _, data := range u.(map[string]interface{}) {
			if !turbineInit {
				turbine = append(turbine, []string{})
			}

			turbine[indexColumn] = append(turbine[indexColumn], fmt.Sprintf("%v", data))
			indexColumn++
		}

		turbineInit = true
	}

	dataTurbine := map[string][][]string{
		"Upper":   upper,
		"Clutch":  clutch,
		"Turbine": turbine,
	}

	return t.makePDF(
		dataTurbine,
		t.CreatedAt,
		turbineResponse.CreatedBy,
		turbineResponse.TotalBolts,
		turbineResponse.CurrentTorque,
		turbineResponse.MaxTorque,
		turbineResponse.Shaft.GenBearingToKopling,
		turbineResponse.Shaft.GenBearingToKopling,
		turbineResponse.Shaft.Ratio,
		turbineResponse.TotalCrockedness,
		turbineResponse.TotalCrockedness <= TotalCrockednessToleration,
	)
}

func (t *Turbine) makePDF(
	dataTurbine map[string][][]string,
	createdAt *time.Time,
	createdBy string,
	totalBolts uint32,
	currentTorque float64,
	maxTorque float64,
	genBearingToKoping float64,
	koplingToTurbine float64,
	ratio float64,
	totalRounOut float64,
	isTotalCrockednessSave bool) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "legal", "")

	// Set font for titles
	pdf.SetFont("Arial", "B", 12)

	// Add a page
	pdf.AddPage()

	// Set position for the first table (below the image)
	pdf.Ln(5) // Move 60mm down from the top

	// Add titletitle
	pdf.CellFormat(0, 7, "LAPORAN DATA TURBINE", "", 1, "C", false, 0, "")

	pdf.Ln(5) // Move 60mm down from the top

	details(pdf, createdBy, createdAt.Format("02/01/2006 - 15.04"), t.Title)

	// Create the first table with a title
	tableData(pdf, "Tabel Data Upper", dataTurbine["Upper"])

	// Create the second table with a title
	tableData(pdf, "Tabel Data Kopling", dataTurbine["Clutch"])

	// Create the third table with a title
	tableData(pdf, "Tabel Data Turbin", dataTurbine["Turbine"])

	// Conclussion
	conclussion(pdf,
		createdBy,
		totalBolts,
		currentTorque,
		maxTorque,
		genBearingToKoping,
		koplingToTurbine,
		ratio,
		totalRounOut,
		isTotalCrockednessSave)

	// Output the PDF to a buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func details(pdf *gofpdf.Fpdf, createdBy, createdAt, title string) {
	// set font to be bold
	pdf.SetFont("Arial", "B", 12)

	// Define column widths
	colWidths := []float64{35, 5, 132}
	tableWidth := 0.0
	for _, w := range colWidths {
		tableWidth += w
	}

	// Center the table by setting the X position
	pdf.SetX((210 - tableWidth) / 2) // 210 is the width of A4 paper in mm

	pdf.Ln(-1) // Move to the next linemain

	// Set font for the table content
	pdf.SetFont("Arial", "", 12)

	columns := [][]string{
		{"Penginput", ":", createdBy},
		{"Tanggal & Jam", ":", createdAt},
		{"Acara", ":", title},
	}
	// Add rows to the table with alternating row colors
	for _, column := range columns {
		// Center each row
		pdf.SetX((210 - tableWidth) / 2)
		for indexRow, w := range colWidths {
			// pdf.SetTextColor(255, 255, 255)
			pdf.CellFormat(w, 6.5, column[indexRow], "", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Add space after the table
	pdf.Ln(5)
}

func tableData(pdf *gofpdf.Fpdf, title string, data [][]string) {
	// // Calculate table height (assuming each row height is 10mm)
	// tableHeight := calculateTableHeight(0, float64(len(columns))*8)

	// // If the table height exceeds the remaining space on the page, add a new page
	// _, pageHeight := pdf.GetPageSize()
	// bottomMargin := 2.0 // Set your bottom margin here
	// if pdf.GetY()+tableHeight > pageHeight-bottomMargin {
	// 	pdf.AddPage()
	// }

	// set font to be bold
	pdf.SetFont("Arial", "B", 12)

	// Add title
	// pdf.CellFormat(0, 10, title, "", 1, "C", false, 0, "")
	pdf.SetX(210 / 11.5) // 210 is the width of A4 paper in mm
	pdf.CellFormat(0, 5, title, "", 1, "", false, 0, "")
	pdf.Ln(2) // Space after the title

	// Define column widths
	colWidths := []float64{10, 40, 40, 40, 40}
	tableWidth := 0.0
	for _, w := range colWidths {
		tableWidth += w
	}

	// Center the table by setting the X position
	pdf.SetX((210 - tableWidth) / 2) // 210 is the width of A4 paper in mm

	// Set fill color for the header (e.g., light gray)
	pdf.SetFillColor(255, 255, 255)

	// Set header for the table
	rows := []string{"No.", "A", "B", "C", "D"}
	pdf.SetFont("Arial", "B", 11)
	for index, w := range colWidths {
		pdf.CellFormat(w, 8, rows[index], "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1) // Move to the next line

	// Set font for the table content
	pdf.SetFont("Arial", "", 12)

	// Add rows to the table with alternating row colors
	for index, column := range data {
		if index%2 == 0 {
			pdf.SetFillColor(240, 240, 240) // Light gray for even rows
		} else {
			pdf.SetFillColor(255, 255, 255) // White for odd rows
		}

		// Center each row
		pdf.SetX((210 - tableWidth) / 2)
		for indexRow, w := range colWidths {
			if indexRow == 0 {
				// pdf.SetFillColor(255, 255, 255)
				pdf.SetFont("Arial", "B", 11)
				pdf.CellFormat(w, 8, strconv.Itoa(int(index+1)), "1", 0, "C", true, 0, "")
			} else {
				pdf.SetFont("Arial", "", 11)
				// if indexRow%2 == 0 {
				// 	pdf.SetFillColor(240, 128, 128)
				// } else {
				// 	pdf.SetFillColor(135, 206, 235)
				// }

				pdf.CellFormat(w, 8, column[indexRow-1], "1", 0, "C", true, 0, "")
			}
		}
		pdf.Ln(-1)
	}

	// Add space after the table
	pdf.Ln(5)
}

func conclussion(
	pdf *gofpdf.Fpdf,
	createdBy string,
	totalBolts uint32,
	currentTorque float64,
	maxTorque float64,
	genBearingToKoping float64,
	koplingToTurbine float64,
	ratio float64,
	totalRounOut float64,
	isTotalCrockednessSave bool) {
	details := [][]string{
		{"Total Baut", fmt.Sprintf("%v", totalBolts)},
		{"Torsi Terkini", fmt.Sprintf("%v", currentTorque)},
		{"Max Torsi", fmt.Sprintf("%v", maxTorque)},
		{"Selisih Torsi", fmt.Sprintf("%v", maxTorque-currentTorque)},
		{"Gen.Bearing - Kopling", fmt.Sprintf("%v", genBearingToKoping)},
		{"Kopling - Turbine", fmt.Sprintf("%v", koplingToTurbine)},
		{"Total", fmt.Sprintf("%v", genBearingToKoping+koplingToTurbine)},
		{"Rasio", fmt.Sprintf("%v", ratio)},
		{"Total Run Out", fmt.Sprintf("%v", totalRounOut)},
	}

	signature := [][]string{
		{"", ""},
		{"", ""},
		{"", ""},
		{"", ""},
		{"", ""},
		{"", fmt.Sprintf("( %s )", createdBy)},
	}

	details = append(details, signature...)

	// // Calculate table height (assuming each row height is 10mm)
	// tableHeight := calculateTableHeight(len(details), float64(len(details))*8)

	// // If the table height exceeds the remaining space on the page, add a new page
	// _, pageHeight := pdf.GetPageSize()
	// bottomMargin := 2.0 // Set your bottom margin here
	// if pdf.GetY()+tableHeight > pageHeight-bottomMargin {
	// 	pdf.AddPage()
	// }

	// set font to be bold
	pdf.SetFont("Arial", "B", 11)

	// Define column widths
	colWidths := []float64{85, 85}
	tableWidth := 0.0
	for _, w := range colWidths {
		tableWidth += w
	}

	// Center the table by setting the X position
	pdf.SetX((210 - tableWidth) / 2) // 210 is the width of A4 paper in mm

	// Set fill color for the header (e.g., light gray)
	pdf.SetFillColor(255, 255, 255)

	// Set font for the table content
	pdf.SetFont("Arial", "", 12)

	// Add rows to the table with alternating row colors
	for index, column := range details {
		if index <= len(details)-len(signature)-1 {
			if index%2 == 0 {
				pdf.SetFillColor(240, 240, 240) // Light gray for even rows
			} else {
				pdf.SetFillColor(255, 255, 255) // White for odd rows
			}

			// Center each row
			pdf.SetX((210 - tableWidth) / 2)
			for indexRow, w := range colWidths {
				if indexRow == 0 {
					if index == len(details)-len(signature)-1 {
						pdf.SetFont("Arial", "B", 13)
						pdf.SetTextColor(255, 255, 255)
						if isTotalCrockednessSave {
							pdf.SetFillColor(0, 128, 0) // GREEN
						} else {
							pdf.SetFillColor(255, 0, 0) // RED
						}
						pdf.CellFormat(w, 8, " "+column[indexRow], "1", 0, "L", true, 0, "")
						pdf.SetTextColor(0, 0, 0)
					} else {
						pdf.SetFont("Arial", "", 13)
						pdf.CellFormat(w, 8, " "+column[indexRow], "1", 0, "L", true, 0, "")
					}
				} else {
					if index == len(details)-len(signature)-1 {
						pdf.SetFont("Arial", "B", 13)
						pdf.SetTextColor(255, 255, 255)
						if isTotalCrockednessSave {
							pdf.SetFillColor(0, 128, 0) // GREEN
						} else {
							pdf.SetFillColor(255, 0, 0) // RED
						}
						pdf.CellFormat(w, 8, column[indexRow], "1", 0, "C", true, 0, "")
						pdf.SetTextColor(0, 0, 0)
					} else {
						pdf.SetFont("Arial", "B", 13)
						pdf.CellFormat(w, 8, column[indexRow], "1", 0, "C", true, 0, "")
					}
				}
			}
		} else {
			// Center each row
			pdf.SetX((210 - tableWidth) / 2)
			for indexRow, w := range colWidths {
				pdf.CellFormat(w, 8, column[indexRow], "", 0, "C", false, 0, "")
			}
		}

		pdf.Ln(-1)
	}
	pdf.Ln(10)
}

// Function to calculate the height of the table
func calculateTableHeight(rowCount int, rowHeight float64) float64 {
	// Header height + row heights
	return 10 + (float64(rowCount) * rowHeight)
}
