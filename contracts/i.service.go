package contract

import (
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
)

type IRoleService interface {
	Create(c echo.Context, in *models.RoleWriteRequest) error
	Update(c echo.Context, in *models.RoleWriteRequest) error
	GetListMaster(c echo.Context, search string) error
	Delete(c echo.Context, id string) error
}

type IDivisionService interface {
	Create(c echo.Context, in *models.DivisionWriteRequest) error
	Update(c echo.Context, in *models.DivisionWriteRequest) error
	GetListMaster(c echo.Context, search string) error
	Delete(c echo.Context, id string) error
}

type IUserService interface {
	CreateUserAdminByAdmin(c echo.Context, in *models.UserAdminCreateByAdminRequest) error
	UpdateByAdmin(c echo.Context, in *models.UserUpdateByAdminRequest) error
	Update(c echo.Context, in *models.UserUpdateRequest) error
	GetDetail(c echo.Context, id string) error
	DeleteByAdmin(c echo.Context, in *models.UserDeleteByAdminRequest) error
	GetListWithPaginateByAdmin(c echo.Context, cursor *helpers.Cursor) error
	ChangePassword(c echo.Context, in *models.UserChangePassword) error
}

type IAuthService interface {
	Register(c echo.Context, in *models.Register) error
	Login(c echo.Context, in *models.Login) error
	RefreshToken(c echo.Context, in *models.RefreshToken) error
}

type ITurbineService interface{}
