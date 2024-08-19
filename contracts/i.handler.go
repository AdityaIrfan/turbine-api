package contract

import (
	"github.com/labstack/echo/v4"
)

type IRoleHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	GetListMaster(c echo.Context) error
	Delete(c echo.Context) error
}

type IDivisionHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	GetListMaster(c echo.Context) error
	GetListWithPaginate(c echo.Context) error
	Delete(c echo.Context) error
}

type IUserHandler interface {
	// ADMIN BY SUPER ADMIN
	CreateUserBySuperAdmin(c echo.Context) error
	UpdateUserBySuperAdmin(c echo.Context) error
	GetDetailUserBySuperAdmin(c echo.Context) error
	DeleteUserBySuperAdmin(c echo.Context) error
	GetListUserWithPaginateBySuperAdmin(c echo.Context) error
	GenerateUserPasswordBySuperAdmin(c echo.Context) error

	// USER BY ADMIN
	CreateUserByAdmin(c echo.Context) error
	UpdateUserByAdmin(c echo.Context) error
	GetDetailUserByAdmin(c echo.Context) error
	DeleteUserByAdmin(c echo.Context) error
	GetListUserWithPaginateByAdmin(c echo.Context) error
	GenerateUserPasswordByAdmin(c echo.Context) error

	// USER ITSELF
	UpdateMyProfile(c echo.Context) error
	GetMyProfile(c echo.Context) error
	ChangeMyPassword(c echo.Context) error
}

type IAuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	RefreshToken(c echo.Context) error
	Logout(c echo.Context) error
}

type IConfigHandler interface {
	SaveOrUpdate(c echo.Context) error
	GetRootLocation(c echo.Context) error
}

type ITurbineHandler interface {
	Create(c echo.Context) error
	GetDetail(c echo.Context) error
	GetList(c echo.Context) error
	GetLatest(c echo.Context) error
	Delete(c echo.Context) error
	DownloadReport(c echo.Context) error
}

type IPltaHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	Detail(c echo.Context) error
	GetListMaster(c echo.Context) error
	Delete(c echo.Context) error
	GetListWithPaginate(c echo.Context) error
}

type IPltaUnitHandler interface {
	CreateOrUpdate(c echo.Context) error
	Delete(c echo.Context) error
}

type IDashboardHandler interface {
	GetDashboardData(c echo.Context) error
}
