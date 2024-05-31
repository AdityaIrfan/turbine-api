package handlers

import (
	"fmt"
	"net/http"
	"strings"
	helpers "turbine-api/helpers"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

var Divisions []Division

type DivisionType string

const DivisionType_Engineer = "engineer"

type Division struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func InitDivisions() {
	Divisions = []Division{
		{
			Id:   ulid.Make().String(),
			Name: DivisionType_Engineer,
		},
		{
			Id:   ulid.Make().String(),
			Name: "Finance",
		},
		{
			Id:   ulid.Make().String(),
			Name: "Human Resource",
		},
		{
			Id:   ulid.Make().String(),
			Name: "Staff",
		},
	}

	fmt.Println("SUCCESS INIT DIVISIONS")
}

type divisionHandler struct{}

func NewDivisionHandler() *divisionHandler {
	return &divisionHandler{}
}

func (d *divisionHandler) GetListMasterData(c echo.Context) error {
	return Response(c, http.StatusOK, "success get list master data", Divisions)
}

func (d *divisionHandler) Add(c echo.Context) error {
	type addDivision struct {
		Name string `json:"Name" validate:"required"`
	}

	var payload = new(addDivision)

	if err := c.Bind(payload); err != nil {
		return Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return Response(c, http.StatusBadRequest, errMessage)
	}

	for _, d := range Divisions {
		if strings.EqualFold(d.Name, payload.Name) {
			return Response(c, http.StatusBadRequest, "division already exists")
		}
	}

	division := Division{
		Id:   ulid.Make().String(),
		Name: payload.Name,
	}

	Divisions = append(Divisions, division)

	return Response(c, http.StatusOK, "success add division", division)
}

func (d *divisionHandler) Update(c echo.Context) error {
	isDivisionExist := false
	var division Division
	var divisionIndex int
	for index, d := range Divisions {
		if d.Id == c.Param("id") {
			isDivisionExist = true
			division = d
			divisionIndex = index
		}
	}

	if !isDivisionExist {
		return Response(c, http.StatusBadRequest, "division not found")
	}

	type updateDivision struct {
		Name string `json:"Name" validate:"required"`
	}

	var payload = new(updateDivision)

	if err := c.Bind(payload); err != nil {
		return Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return Response(c, http.StatusBadRequest, errMessage)
	}

	for _, d := range Divisions {
		if strings.EqualFold(d.Name, payload.Name) && d.Id != division.Id {
			return Response(c, http.StatusBadRequest, "division already exists")
		}
	}

	division.Name = payload.Name
	Divisions[divisionIndex].Name = payload.Name

	return Response(c, http.StatusOK, "success update division", division)
}

func (d *divisionHandler) Delete(c echo.Context) error {
	isDivisionExist := false
	var divisionIndex int
	for index, d := range Divisions {
		if d.Id == c.Param("id") {
			isDivisionExist = true
			divisionIndex = index
		}
	}

	if !isDivisionExist {
		return Response(c, http.StatusBadRequest, "division not found")
	}

	Divisions = append(Divisions[:divisionIndex], Divisions[divisionIndex+1:]...)

	return Response(c, http.StatusOK, "success delete division")
}
