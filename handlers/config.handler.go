package handlers

import (
	contract "turbine-api/contracts"

	"github.com/labstack/echo/v4"
)

type configHandler struct {
	configService contract.IConfigService
}

func NewConfigHandler(configService contract.IConfigService) contract.IConfigHandler {
	return &configHandler{
		configService: configService,
	}
}

func (ch *configHandler) GetRootLocation(c echo.Context) error {
	return ch.configService.GetRootLocation(c)
}
