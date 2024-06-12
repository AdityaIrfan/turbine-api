package contract

import (
	"github.com/labstack/echo/v4"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"
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
	GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error
	Delete(c echo.Context, in *models.DivisionWriteRequest) error
}

type IUserService interface {
	CreateUserAdminByAdmin(c echo.Context, in *models.UserAdminCreateByAdminRequest) error
	UpdateByAdmin(c echo.Context, in *models.UserUpdateByAdminRequest) error
	Update(c echo.Context, in *models.UserUpdateRequest) error
	GetDetailByAdmin(c echo.Context, in *models.UserGetDetailRequest) error
	GetMyProfile(c echo.Context, id string) error
	DeleteByAdmin(c echo.Context, in *models.UserDeleteByAdminRequest) error
	GetListWithPaginateByAdmin(c echo.Context, cursor *helpers.Cursor) error
	ChangePassword(c echo.Context, in *models.UserChangePasswordRequest) error
	GeneratePasswordByAdmin(c echo.Context, in *models.GeneratePasswordByAdmin) error
}

type IAuthService interface {
	Register(c echo.Context, in *models.Register) error
	Login(c echo.Context, in *models.Login) error
	RefreshToken(c echo.Context, in *models.RefreshTokenRequest) error
	Logout(c echo.Context, token string) error
}

type IConfigService interface {
	SaveOrUpdate(c echo.Context, in *models.ConfigRootLocation) error
	GetRootLocation(c echo.Context) error
}

type ITurbineService interface {
	Create(c echo.Context, in *models.TurbineWriteRequest) error
	GetDetail(c echo.Context, id string) error
	GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error
}

type ITowerService interface {
	Create(c echo.Context, in *models.TowerWriteRequest) error
	Update(c echo.Context, in *models.TowerWriteRequest) error
	GetListMaster(c echo.Context, search string) error
	Delete(c echo.Context, id string) error
}
