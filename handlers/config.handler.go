package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	helpers "pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
)

type configHandler struct {
	configService contract.IConfigService
}

func NewConfigHandler(configService contract.IConfigService) contract.IConfigHandler {
	return &configHandler{
		configService: configService,
	}
}

func (ch *configHandler) SaveOrUpdate(c echo.Context) error {
	payload := new(models.ConfigRootLocation)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	return ch.configService.SaveOrUpdate(c, payload)
}

func (ch *configHandler) GetRootLocation(c echo.Context) error {
	return ch.configService.GetRootLocation(c)
}
