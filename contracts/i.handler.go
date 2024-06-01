package contract

import (
	"turbine-api/helpers"

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
	Delete(c echo.Context) error
}

type IUserHandler interface {
	CreateUserAdminByAdmin(c echo.Context) error
	UpdateByAdmin(c echo.Context) error
	Update(c echo.Context) error
	GetDetail(c echo.Context, id string) error
	DeleteByAdmin(c echo.Context) error
	GetListWithPaginateByAdmin(c echo.Context, cursor *helpers.Cursor) error
	ChangePassword(c echo.Context) error
}

type AuthHandler interface {
	Register(c echo.Context) error
	Login(c echo.Context) error
	RefreshToken(c echo.Context) error
}

type ITurbineInterface interface{}
