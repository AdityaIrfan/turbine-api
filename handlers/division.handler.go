package handlers

import (
	"net/http"
	contract "turbine-api/contracts"
	helpers "turbine-api/helpers"
	"turbine-api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type divisionHandler struct {
	divisionService contract.IDivisionService
}

func NewDivisionHandler(divisionService contract.IDivisionService) contract.IDivisionHandler {
	return &divisionHandler{
		divisionService: divisionService,
	}
}

func (d *divisionHandler) Create(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	payload := new(models.DivisionWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.AdminId = adminId

	return d.divisionService.Create(c, payload)
}
func (d *divisionHandler) Update(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	payload := new(models.DivisionWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Param("id")
	payload.AdminId = adminId

	return d.divisionService.Update(c, payload)
}

func (d *divisionHandler) GetListMaster(c echo.Context) error {
	return d.divisionService.GetListMaster(c, c.QueryParam("Search"))
}

func (d *divisionHandler) GetListWithPaginate(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	cursor, err := helpers.GenerateCursorPaginationByEcho(c)
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return d.divisionService.GetListWithPaginate(c, cursor, adminId)
}

func (d *divisionHandler) Delete(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	payload := &models.DivisionWriteRequest{
		Id:      c.Param("id"),
		AdminId: adminId,
	}
	return d.divisionService.Delete(c, payload)
}
