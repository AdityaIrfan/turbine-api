package handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	helpers "pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
)

type turbineHandler struct {
	turbineService contract.ITurbineService
}

func NewTurbineHandler(turbineService contract.ITurbineService) contract.ITurbineHandler {
	return &turbineHandler{
		turbineService: turbineService,
	}
}

func (t *turbineHandler) Create(c echo.Context) error {
	payload := new(models.TurbineWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}
	if err := payload.ValidateData(); err != nil {
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	payload.CreatedBy = c.Get("claims").(jwt.MapClaims)["Id"].(string)
	return t.turbineService.Create(c, payload)
}

func (t *turbineHandler) GetDetail(c echo.Context) error {
	return t.turbineService.GetDetail(c, c.Param("id"))
}

func (t *turbineHandler) GetList(c echo.Context) error {
	cursor, err := helpers.GenerateCursorPaginationByEcho(c, models.TurbineDefaultMap)
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return t.turbineService.GetListWithPaginate(c, cursor)
}
