package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	contract "pln/AdityaIrfan/turbine-api/contracts"
	helpers "pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
)

type roleHandler struct {
	roleService contract.IRoleService
}

func NewRoleHandler(roleService contract.IRoleService) contract.IRoleHandler {
	return &roleHandler{
		roleService: roleService,
	}
}

func (r *roleHandler) Create(c echo.Context) error {
	payload := new(models.RoleWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	return r.roleService.Create(c, payload)
}
func (r *roleHandler) Update(c echo.Context) error {
	payload := new(models.RoleWriteRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Param("id")

	return r.roleService.Update(c, payload)
}

func (r *roleHandler) GetListMaster(c echo.Context) error {
	return r.roleService.GetListMaster(c, c.QueryParam("Search"))
}

func (r *roleHandler) Delete(c echo.Context) error {
	return r.roleService.Delete(c, c.Param("id"))
}
