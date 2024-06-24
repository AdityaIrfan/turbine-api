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
	CreateUserAdminByAdmin(c echo.Context) error
	UpdateByAdmin(c echo.Context) error
	Update(c echo.Context) error
	GetDetailByAdmin(c echo.Context) error
	GetMyProfile(c echo.Context) error
	DeleteByAdmin(c echo.Context) error
	GetListWithPaginateByAdmin(c echo.Context) error
	ChangePassword(c echo.Context) error
	GeneratePasswordByAdmin(c echo.Context) error
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
}

type ITowerHandler interface {
	Create(c echo.Context) error
	Update(c echo.Context) error
	GetListMaster(c echo.Context) error
	Delete(c echo.Context) error
}
