package handlers

import (
	"net/http"
	"strings"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	helpers "pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type pltaHandler struct {
	pltaService contract.IPltaService
}

func NewPltaHandler(pltaService contract.IPltaService) contract.IPltaHandler {
	return &pltaHandler{
		pltaService: pltaService,
	}
}

func (t *pltaHandler) Create(c echo.Context) error {
	payload := new(models.PltaCreateRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	if !helpers.IsValidLatLong(payload.Lat, payload.Long) {
		return helpers.Response(c, http.StatusBadRequest, "latitude dan longitude tidak valid")
	}

	payload.WrittenBy = c.Get("claims").(jwt.MapClaims)["Id"].(string)
	return t.pltaService.Create(c, payload)
}

func (t *pltaHandler) Update(c echo.Context) error {
	payload := new(models.PltaUpdateRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	if !helpers.IsValidLatLong(payload.Lat, payload.Long) {
		return helpers.Response(c, http.StatusBadRequest, "latitude dan longitude tidak valid")
	}

	payload.Id = c.Param("id")
	payload.WrittenBy = c.Get("claims").(jwt.MapClaims)["Id"].(string)
	return t.pltaService.Update(c, payload)
}

func (t *pltaHandler) Detail(c echo.Context) error {
	return t.pltaService.Detail(c, c.Param("id"))
}

func (t *pltaHandler) GetListMaster(c echo.Context) error {
	return t.pltaService.GetListMaster(c, &models.PltaGetListMasterRequest{
		UserId: c.Get("claims").(jwt.MapClaims)["Id"].(string),
		Search: c.QueryParam("Search"),
	})
}

func (t *pltaHandler) Delete(c echo.Context) error {
	payload := &models.PltaDeleteRequest{
		Id:        c.Param("id"),
		DeletedBy: c.Get("claims").(jwt.MapClaims)["Id"].(string),
	}
	return t.pltaService.Delete(c, payload)
}

func (u *pltaHandler) GetListWithPaginate(c echo.Context) error {
	cursor, err := helpers.GenerateCursorPaginationByEcho(c, models.PltaDefaultSort, nil)
	if err != nil {
		if strings.Contains(err.Error(), "unavailable") {
			return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua user", []models.User{}, helpers.CursorPagination{})
		}
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return u.pltaService.GetListWithPaginate(c, cursor)
}
