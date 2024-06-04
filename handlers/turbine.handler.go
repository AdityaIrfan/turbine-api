package handlers

import (
	"net/http"
	contract "turbine-api/contracts"
	helpers "turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type turbineHandler struct{}

func NewTurbineHandler() contract.ITurbineHandler {
	return &turbineHandler{}
}

func (t *turbineHandler) Create(c echo.Context) error {
	payload := new(models.TurbineWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}
	if err := payload.ValidateData(); err != nil {
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return helpers.Response(c, http.StatusOK, "good payload")
}
