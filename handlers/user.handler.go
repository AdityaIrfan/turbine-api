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

type userHandler struct {
	userService contract.IUserService
}

func NewUserHandler(userService contract.IUserService) contract.IUserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (u *userHandler) CreateUserByAdmin(c echo.Context) error {
	payload := new(models.UserCreateByAdminRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.CreatedBy = c.Get("claims").(jwt.MapClaims)["Id"].(string)
	return u.userService.CreateUserByAdmin(c, payload)
}

func (u *userHandler) UpdateUserByAdmin(c echo.Context) error {
	payload := new(models.UserUpdateByAdminRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Param("id")
	payload.UpdatedBy = c.Get("claims").(jwt.MapClaims)["Id"].(string)
	return u.userService.UpdateUserByAdmin(c, payload)
}

func (u *userHandler) DeleteUserByAdmin(c echo.Context) error {
	payload := &models.UserDeleteByAdminRequest{
		Id:      c.Param("id"),
		AdminId: c.Get("claims").(jwt.MapClaims)["Id"].(string),
	}

	return u.userService.DeleteUserByAdmin(c, payload)
}

func (u *userHandler) GetDetailUserByAdmin(c echo.Context) error {
	payload := &models.UserGetDetailRequest{
		Id: c.Param("id"),
	}
	return u.userService.GetDetailUserByAdmin(c, payload)
}

func (u *userHandler) GetListUserWithPaginateByAdmin(c echo.Context) error {
	cursor, err := helpers.GenerateCursorPaginationByEcho(c, models.UserDefaultSort, models.UserDefaultFilter)
	if err != nil {
		if strings.Contains(err.Error(), "unavailable") {
			return helpers.Response(c, http.StatusOK, "berhasil mendapatkan semua user", []models.User{}, helpers.CursorPagination{})
		}
		return helpers.Response(c, http.StatusBadRequest, err.Error())
	}

	return u.userService.GetListUserWithPaginateByAdmin(c, cursor)
}

// func (u *userHandler) GenerateUserPasswordByAdmin(c echo.Context) error {
// 	payload := &models.GeneratePasswordByAdmin{
// 		Id:      c.Param("id"),
// 		AdminId: c.Get("claims").(jwt.MapClaims)["Id"].(string),
// 	}

// 	return u.userService.GenerateUserPasswordByAdmin(c, payload)
// }

func (u *userHandler) UpdateMyProfile(c echo.Context) error {
	payload := new(models.UserUpdateRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	payload.Id = c.Get("claims").(jwt.MapClaims)["Id"].(string)
	return u.userService.UpdateMyProfile(c, payload)
}

func (u *userHandler) GetMyProfile(c echo.Context) error {
	return u.userService.GetMyProfile(c, c.Get("claims").(jwt.MapClaims)["Id"].(string))
}

func (u *userHandler) ChangeMyPassword(c echo.Context) error {
	payload := new(models.UserChangePasswordRequest)

	if err := c.Bind(payload); err != nil {
		return helpers.Response(c, http.StatusBadRequest, "payload tidak valid")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return helpers.Response(c, http.StatusBadRequest, errMessage)
	}

	if payload.Password != payload.PasswordConfirmation {
		return helpers.Response(c, http.StatusBadRequest, "Password dan PasswordConfirmation tidak sama")
	}

	payload.Id = c.Get("claims").(jwt.MapClaims)["Id"].(string)

	return u.userService.ChangeMyPassword(c, payload)
}
