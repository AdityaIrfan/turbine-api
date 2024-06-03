package handlers

import (
	"net/http"
	contract "turbine-api/contracts"
	helpers "turbine-api/helpers"
	"turbine-api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userService contract.IUserService
}

func NewUserHandler(userService contract.IUserService) contract.IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (u *userHandler) CreateUserAdminByAdmin(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	payload := new(models.UserAdminCreateByAdminRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.AdminId = adminId

	return u.userService.CreateUserAdminByAdmin(c, payload)
}

func (u *userHandler) UpdateByAdmin(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	payload := new(models.UserUpdateByAdminRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.AdminId = adminId

	return u.userService.UpdateByAdmin(c, payload)
}

func (u *userHandler) Update(c echo.Context) error {
	payload := new(models.UserUpdateRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Get("claims").(jwt.MapClaims)["id"].(string)

	return u.userService.Update(c, payload)
}

func (u *userHandler) GetDetailByAdmin(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	payload := &models.UserGetDetailRequest{
		Id:      c.Param("id"),
		AdminId: adminId,
	}
	return u.userService.GetDetailByAdmin(c, payload)
}

func (u *userHandler) GetMyProfile(c echo.Context) error {
	return u.userService.GetMyProfile(c, c.Get("claims").(jwt.MapClaims)["id"].(string))
}

func (u *userHandler) DeleteByAdmin(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	payload := &models.UserDeleteByAdminRequest{
		Id:      c.Param("id"),
		AdminId: adminId,
	}

	return u.userService.DeleteByAdmin(c, payload)
}

func (u *userHandler) GetListWithPaginateByAdmin(c echo.Context) error {
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

	return u.userService.GetListWithPaginateByAdmin(c, cursor, adminId)
}

func (u *userHandler) ChangePassword(c echo.Context) error {
	payload := new(models.UserChangePasswordRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Get("claims").(jwt.MapClaims)["id"].(string)

	return u.userService.ChangePassword(c, payload)
}

func (u *userHandler) GeneratePasswordByAdmin(c echo.Context) error {
	var adminId string
	if value, ok := c.Get("claims").(jwt.MapClaims)["Id"].(string); !ok {
		return helpers.ResponseNonAdminForbiddenAccess(c)
	} else {
		adminId = value
	}

	payload := &models.GeneratePasswordByAdmin{
		Id:      c.Param("id"),
		AdminId: adminId,
	}

	return u.userService.GeneratePasswordByAdmin(c, payload)
}
