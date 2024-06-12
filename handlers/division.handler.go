package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	helpers "pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
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
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
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
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
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

func (d *divisionHandler) GetListWithPaginate(c echo.Context) error {
	cursor, err := helpers.GenerateCursorPaginationByEcho(c, models.DivisionDefaultSort)
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return d.divisionService.GetListWithPaginate(c, cursor)
}

func (d *divisionHandler) Delete(c echo.Context) error {
	payload := &models.DivisionWriteRequest{
		Id: c.Param("id"),
	}
	return d.divisionService.Delete(c, payload)
}
