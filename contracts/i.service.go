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
	GetListWithPaginate(c echo.Context, cursor *helpers.Cursor, adminId string) error
	Delete(c echo.Context, in *models.DivisionWriteRequest) error
}

type IUserService interface {
	CreateUserAdminByAdmin(c echo.Context, in *models.UserAdminCreateByAdminRequest) error
	UpdateByAdmin(c echo.Context, in *models.UserUpdateByAdminRequest) error
	Update(c echo.Context, in *models.UserUpdateRequest) error
	GetDetailByAdmin(c echo.Context, in *models.UserGetDetailRequest) error
	GetMyProfile(c echo.Context, id string) error
	DeleteByAdmin(c echo.Context, in *models.UserDeleteByAdminRequest) error
	GetListWithPaginateByAdmin(c echo.Context, cursor *helpers.Cursor, adminId string) error
	ChangePassword(c echo.Context, in *models.UserChangePasswordRequest) error
	GeneratePasswordByAdmin(c echo.Context, in *models.GeneratePasswordByAdmin) error
}

type IAuthService interface {
	Register(c echo.Context, in *models.Register) error
	Login(c echo.Context, in *models.Login) error
	RefreshToken(c echo.Context, in *models.RefreshTokenRequest) error
}

type IConfigService interface {
	GetRootLocation(c echo.Context, adminId string) error
}

type ITurbineService interface{}
