package handlers

import (
	"net/http"

	contract "github.com/AdityaIrfan/turbine-api/contracts"
	helpers "github.com/AdityaIrfan/turbine-api/helpers"
	"github.com/AdityaIrfan/turbine-api/models"
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
	payload := new(models.UserAdminCreateByAdminRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	return u.userService.CreateUserAdminByAdmin(c, payload)
}

func (u *userHandler) UpdateByAdmin(c echo.Context) error {
	payload := new(models.UserUpdateByAdminRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	return u.userService.UpdateByAdmin(c, payload)
}

func (u *userHandler) Update(c echo.Context) error {
	payload := new(models.UserUpdateRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Get("claims").(jwt.MapClaims)["id"].(string)

	return u.userService.Update(c, payload)
}

func (u *userHandler) GetDetailByAdmin(c echo.Context) error {
	payload := &models.UserGetDetailRequest{
		Id: c.Param("id"),
	}
	return u.userService.GetDetailByAdmin(c, payload)
}

func (u *userHandler) GetMyProfile(c echo.Context) error {
	return u.userService.GetMyProfile(c, c.Get("claims").(jwt.MapClaims)["id"].(string))
}

func (u *userHandler) DeleteByAdmin(c echo.Context) error {
	payload := &models.UserDeleteByAdminRequest{
		Id:      c.Param("id"),
		AdminId: c.Get("claims").(jwt.MapClaims)["Id"].(string),
	}

	return u.userService.DeleteByAdmin(c, payload)
}

func (u *userHandler) GetListWithPaginateByAdmin(c echo.Context) error {
	cursor, err := helpers.GenerateCursorPaginationByEcho(c, models.UserDefaultSort)
	if err != nil {
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return u.userService.GetListWithPaginateByAdmin(c, cursor)
}

func (u *userHandler) ChangePassword(c echo.Context) error {
	payload := new(models.UserChangePasswordRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Get("claims").(jwt.MapClaims)["id"].(string)

	return u.userService.ChangePassword(c, payload)
}

func (u *userHandler) GeneratePasswordByAdmin(c echo.Context) error {
	payload := &models.GeneratePasswordByAdmin{
		Id:      c.Param("id"),
		AdminId: c.Get("claims").(jwt.MapClaims)["Id"].(string),
	}

	return u.userService.GeneratePasswordByAdmin(c, payload)
}
