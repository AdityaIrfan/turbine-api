package handlers

import (
	"net/http"

	contract "github.com/AdityaIrfan/turbine-api/contracts"
	helpers "github.com/AdityaIrfan/turbine-api/helpers"
	"github.com/AdityaIrfan/turbine-api/models"
	"github.com/labstack/echo/v4"
)

type towerHandler struct {
	towerService contract.ITowerService
}

func NewTowerHandler(towerService contract.ITowerService) contract.ITowerHandler {
	return &towerHandler{
		towerService: towerService,
	}
}

func (t *towerHandler) Create(c echo.Context) error {
	payload := new(models.TowerWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	return t.towerService.Create(c, payload)
}

func (t *towerHandler) Update(c echo.Context) error {
	payload := new(models.TowerWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Param("id")

	return t.towerService.Update(c, payload)
}

func (t *towerHandler) GetListMaster(c echo.Context) error {
	return t.towerService.GetListMaster(c, c.QueryParam("Search"))
}

func (t *towerHandler) Delete(c echo.Context) error {
	return t.towerService.Delete(c, c.Param("id"))
}
