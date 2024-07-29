package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	helpers "pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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

	if c.FormValue("Data") != "" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(c.FormValue("Data")), &data); err != nil {
			return helpers.Response(c, http.StatusBadRequest, "Data must be json stringify")
		}
		payload.Data = data
		c.Request().Form.Del("Data")
	}

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
	cursor, err := helpers.GenerateCursorPaginationByEcho(c, models.TurbineDefaultSortMap, models.TurbineDefaultSortFilter)
	if err != nil {
		if strings.Contains(err.Error(), "unavailable") {
			return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua data turbine", []models.Turbine{}, helpers.CursorPagination{})
		}
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return t.turbineService.GetListWithPaginate(c, cursor)
}

func (t *turbineHandler) GetLatest(c echo.Context) error {
	return t.turbineService.GetLatest(c)
}
