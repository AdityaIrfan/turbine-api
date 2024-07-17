package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
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
	TowerId              string          `gorm:"column:tower_id"`
	GenBearingToCoupling float64         `gorm:"column:gen_bearing_to_coupling"`
	CouplingToTurbine    float64         `gorm:"column:coupling_to_turbine"`
	Data                 datatypes.JSON  `gorm:"column:data"`
	CreatedAt            *time.Time      `gorm:"column:created_at"`
	UpdatedAt            *time.Time      `gorm:"column:updated_at;<-:update"`
	DeletedAt            *gorm.DeletedAt `gorm:"column:deleted_at"`
	CreatedBy            string          `gorm:"column:created_by"`
	TotalBolts           uint32          `gorm:"column:total_bolts"`
	CurrentTorque        float64         `gorm:"column:current_torque"`
	MaxTorque            float64         `gorm:"column:max_torque"`

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
	TotalBolts           uint32                 `json:"TotalBolts" form:"TotalBolts" validate:"required,min=4"`
	CurrentTorque        float64                `json:"CurrentTorque" form:"CurrentTorque" validate:"required"`
	MaxTorque            float64                `json:"MaxTorque" form:"MaxTorque" validate:"required"`
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
		TotalBolts:           t.TotalBolts,
		CurrentTorque:        t.CurrentTorque,
		MaxTorque:            t.MaxTorque,
	}
}

type TurbineResponse struct {
	Id                string                 `json:"Id"`
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

	chart["Upper"] = fmt.Sprintf("%f|%f", resultanAC, resultanBD)
	totalCrockedness := math.Pow((crockednessAC + crockednessBD), 0.5)

	// TORQUE CALCULATION
	torqueGap := t.MaxTorque - t.CurrentTorque
	totalAngleInDegrees := uint32(360)
	degreeGap := float64(totalAngleInDegrees / t.TotalBolts)
	circleRadius := math.Sqrt(math.Pow(resultanAC, 2) + math.Pow(resultanBD, 2))
	bolt := 1
	TorqueCalculation := make(map[string]interface{})
	TorqueCalculationDetail := make(map[string]interface{})
	points := [][2]float64{}
	pointsTemp := make(map[string]uint32)
	fmt.Println("RADIUS : ", circleRadius)
	for i := float64(totalAngleInDegrees); i > 0; i -= degreeGap {
		angleInDegrees := float64(totalAngleInDegrees) - i
		angleInRadians := angleInDegrees * math.Pi / 180.0

		// Compute the cosine and sin of degrees
		cosValue := math.Cos(angleInRadians)
		sinValue := math.Sin(angleInRadians)

		fmt.Printf("cos(%f degrees) = %v\n", angleInDegrees, cosValue)
		fmt.Printf("sin(%f degrees) = %v\n", angleInDegrees, sinValue)
		fmt.Println()

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
			bolt = 8
		} else {
			bolt--
		}
	}
	TorqueCalculation["Details"] = TorqueCalculationDetail
	TorqueSuggestion := make(map[string]float64)

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

	TorqueCalculation["TorqueSuggestions"] = TorqueSuggestion

	return &TurbineResponse{
		Id:        t.Id,
		TowerName: fmt.Sprintf("%v - %v", t.Tower.Name, t.Tower.UnitNumber),
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
		CreatedBy:         t.User.Name,
		Status:            totalCrockedness <= 3,
		TotalBolts:        t.TotalBolts,
		CurrentTorque:     t.CurrentTorque,
		MaxTorque:         t.MaxTorque,
		TorqueGap:         torqueGap,
		TorqueCalculation: TorqueCalculation,
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
