package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

type RoleType string

var RoleType_Admin = "admin"
var RoleType_User = "user"

var Roles []Role

type Role struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func InitRole() {
	Roles = []Role{
		{
			Id:   ulid.Make().String(),
			Name: RoleType_Admin,
		},
		{
			Id:   ulid.Make().String(),
			Name: RoleType_User,
		},
	}

	fmt.Println("SUCCESS INIT ROLES")
}

type roleHandler struct{}

func NewRoleHandler() *roleHandler {
	return &roleHandler{}
}

func (roleHandler) GetListMasterData(c echo.Context) error {
	if c.Request().Header.Get("Authorization") == "" {
		return Response(c, http.StatusBadRequest, "missing authorization header")
	}

	if len(strings.Split(c.Request().Header.Get("Authorization"), " ")) != 2 {
		return Response(c, http.StatusBadRequest, "invalid authorization")
	}

	return Response(c, http.StatusOK, "success get roles", Roles)
}
