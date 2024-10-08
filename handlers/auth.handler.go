package handlers

import (
	"fmt"
	"net/http"
	"strings"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	helpers "pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/labstack/echo/v4"
)

var RefreshTokenMap = map[string]interface{}{}

func InitRefreshToken() {
	RefreshTokenMap = make(map[string]interface{})

	fmt.Println("SUCCESS INIT REFRESH TOKEN")
}

type authHandler struct {
	authService contract.IAuthService
}

func NewAuthHandler(authService contract.IAuthService) contract.IAuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (a *authHandler) Register(c echo.Context) error {
	payload := new(models.Register)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	if err := helpers.ValidatePhone(payload.Phone); err != nil {
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return a.authService.Register(c, payload)
}

func (a *authHandler) Login(c echo.Context) error {
	payload := new(models.Login)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	return a.authService.Login(c, payload)
}

func (a *authHandler) RefreshToken(c echo.Context) error {
	payload := new(models.RefreshTokenRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	return a.authService.RefreshToken(c, payload)
}

func (a *authHandler) Logout(c echo.Context) error {
	tokens := strings.Split(c.Request().Header.Get("Authorization"), " ")

	return a.authService.Logout(c, tokens[1])
}
