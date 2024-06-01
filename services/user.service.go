package services

import (
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type userService struct {
	userRepo contract.IUserRepository
}

func NewUserService(userRepo contract.IUserRepository) contract.IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUserAdminByAdmin(c echo.Context, in *models.UserAdminCreateByAdminRequest) error {
	return nil
}

func (u *userService) UpdateByAdmin(c echo.Context, in *models.UserUpdateByAdminRequest) error {
	return nil
}

func (u *userService) Update(c echo.Context, in *models.UserUpdateRequest) error {
	return nil
}

func (u *userService) GetDetail(c echo.Context, id string) error {
	return nil
}

func (u *userService) DeleteByAdmin(c echo.Context, in *models.UserDeleteByAdminRequest) error {
	return nil
}

func (u *userService) GetListWithPaginateByAdmin(c echo.Context, cursor *helpers.Cursor) error {
	return nil
}

func (u *userService) ChangePassword(c echo.Context, in *models.UserChangePassword) error {
	return nil
}
