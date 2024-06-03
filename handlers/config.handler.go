package handlers

import (
	contract "turbine-api/contracts"
	helpers "turbine-api/helpers"

	"github.com/golang-jwt/jwt/v5"
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
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	return ch.configService.GetRootLocation(c, adminId)
}
