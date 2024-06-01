package handlers

import (
	"net/http"
	contract "turbine-api/contracts"
	helpers "turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type divisionHandler struct {
	divisionService contract.IDivisionService
}

func NewDivisionHandler(divisionService contract.IDivisionService) contract.IDivisionHandler {
	return &divisionHandler{
		divisionService: divisionService,
	}
}

func (d *divisionHandler) Create(c echo.Context) error {
	payload := new(models.DivisionWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	return d.divisionService.Create(c, payload)
}
func (d *divisionHandler) Update(c echo.Context) error {
	payload := new(models.DivisionWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Param("id")

	return d.divisionService.Update(c, payload)
}

func (d *divisionHandler) GetListMaster(c echo.Context) error {
	return d.divisionService.GetListMaster(c, c.QueryParam("Search"))
}

func (d *divisionHandler) Delete(c echo.Context) error {
	return d.divisionService.Delete(c, c.Param("id"))
}
