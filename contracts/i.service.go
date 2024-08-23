package contract

import (
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

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
	GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error
	Delete(c echo.Context, in *models.DivisionWriteRequest) error
}

type IUserService interface {
	// ADMIN BY SUPER ADMIN
	CreateUserBySuperAdmin(c echo.Context, in *models.UserCreateBySuperAdminRequest) error
	UpdateUserBySuperAdmin(c echo.Context, in *models.UserUpdateBySuperAdminRequest) error
	GetDetailUserBySuperAdmin(c echo.Context, in *models.UserGetDetailRequest) error
	DeleteUserBySuperAdmin(c echo.Context, in *models.UserDeleteBySuperAdminRequest) error
	GetListUserWithPaginateBySuperAdmin(c echo.Context, cursor *helpers.Cursor) error
	GenerateUserPasswordBySuperAdmin(c echo.Context, in *models.GeneratePasswordBySuperAdmin) error

	// USER BY ADMIN
	CreateUserByAdmin(c echo.Context, in *models.UserCreateByAdminRequest) error
	UpdateUserByAdmin(c echo.Context, in *models.UserUpdateByAdminRequest) error
	GetDetailUserByAdmin(c echo.Context, in *models.UserGetDetailRequest) error
	DeleteUserByAdmin(c echo.Context, in *models.UserDeleteByAdminRequest) error
	GetListUserWithPaginateByAdmin(c echo.Context, cursor *helpers.Cursor) error
	GenerateUserPasswordByAdmin(c echo.Context, in *models.GeneratePasswordByAdmin) error

	// USER ITSELF
	UpdateMyProfile(c echo.Context, in *models.UserUpdateRequest) error
	GetMyProfile(c echo.Context, id string) error
	ChangeMyPassword(c echo.Context, in *models.UserChangePasswordRequest) error
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
	GetLatest(c echo.Context) error
	Delete(c echo.Context, in *models.TurbineWriteRequest) error
	DownloadReport(c echo.Context, id string) error
}

type IPltaService interface {
	Create(c echo.Context, in *models.PltaCreateRequest) error
	Update(c echo.Context, in *models.PltaUpdateRequest) error
	Detail(c echo.Context, id string) error
	GetListMaster(c echo.Context, in *models.PltaGetListMasterRequest) error
	Delete(c echo.Context, in *models.PltaDeleteRequest) error
	GetListWithPaginate(c echo.Context, cursor *helpers.Cursor) error
}

type IPltaUnitService interface {
	CreateOrUpdate(c echo.Context, in *models.PltaUnitCreateOrUpdate) error
	Delete(c echo.Context, in *models.PltaUnitWriteRequest) error
	GetListMaster(c echo.Context, in *models.PltaGetListMasterRequest) error
}

type IDashboardService interface {
	GetDashboardData(c echo.Context) error
}
