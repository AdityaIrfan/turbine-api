package handlers

import (
	"encoding/json"
	"net/http"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	helpers "pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type pltaUnitHandler struct {
	pltaUnitService contract.IPltaUnitService
}

func NewPltaUnitHandler(pltaUnitService contract.IPltaUnitService) contract.IPltaUnitHandler {
	return &pltaUnitHandler{
		pltaUnitService: pltaUnitService,
	}
}

func (p *pltaUnitHandler) CreateOrUpdate(c echo.Context) error {
	payload := new(models.PltaUnitCreateOrUpdate)

	if c.FormValue("Units") != "" {
		var units []models.PltaUnitWriteRequest
		if err := json.Unmarshal([]byte(c.FormValue("Units")), &units); err != nil {
			return helpers.Response(c, http.StatusBadRequest, "Units must be json stringify")
		}
		payload.Units = units
		c.Request().Form.Del("Units")
	}

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.PltaId = c.Param("id")
	payload.WrittenBy = c.Get("claims").(jwt.MapClaims)["Id"].(string)
	return p.pltaUnitService.CreateOrUpdate(c, payload)
}

func (p *pltaUnitHandler) Delete(c echo.Context) error {
	payload := &models.PltaUnitWriteRequest{
		Id:        c.Param("id"),
		WrittenBy: c.Get("claims").(jwt.MapClaims)["Id"].(string),
	}

	return p.pltaUnitService.Delete(c, payload)
}
