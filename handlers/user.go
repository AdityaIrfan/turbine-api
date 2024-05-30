package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

var Users = []User{}

type User struct {
	Id         string `json:"id"`
	Name       string `json:"Name"`
	Username   string `json:"Username"`
	DivisionId string `json:"DivisionId,omitempty"`
	RoleId     string `json:"RoleId,omitempty"`
	Password   string `json:"Password,omitempty"`
}

func InitUsers() {
	var roleId, divisionId string

loopRoles:
	for _, r := range Roles {
		if r.Name == RoleType_Admin {
			roleId = r.Id
			break loopRoles
		}
	}

loopDivisions:
	for _, d := range Divisions {
		if d.Name == DivisionType_Engineer {
			divisionId = d.Id
			break loopDivisions
		}
	}

	Users = append(Users, User{
		Id:         ulid.Make().String(),
		Name:       "aditya",
		Username:   "aditya",
		DivisionId: divisionId,
		RoleId:     roleId,
		Password:   "password",
	})

	fmt.Println("SUCCESS INIT USERS")
}

type userHandler struct{}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (userHandler) GetList(c echo.Context) error {
	search := c.QueryParam("Search")
	type userResponse struct {
		Id       string `json:"id"`
		Name     string `json:"Name"`
		Division string `json:"Division"`
	}

	var users = []userResponse{}

	if search != "" {
		var tempDivision = make(map[string]string)

		for _, u := range Users {
			var division string

			if value, ok := tempDivision[u.DivisionId]; ok {
				division = value
			} else {
			loopThat:
				for _, d := range Divisions {
					if d.Id == u.DivisionId {
						division = d.Name
						tempDivision[d.Id] = d.Name
						break loopThat
					}
				}
			}

			if strings.Contains(strings.ToLower(u.Name), strings.ToLower(search)) ||
				strings.Contains(strings.ToLower(division), strings.ToLower(search)) {
				users = append(users, userResponse{
					Id:       u.Id,
					Name:     u.Name,
					Division: division,
				})
			}
		}
	} else {
		var tempDivision = make(map[string]string)

		for _, u := range Users {
			var division string

			if value, ok := tempDivision[u.DivisionId]; ok {
				division = value
			} else {
			loopThis:
				for _, d := range Divisions {
					if d.Id == u.DivisionId {
						division = d.Name
						tempDivision[d.Id] = d.Name
						break loopThis
					}
				}
			}

			users = append(users, userResponse{
				Id:       u.Id,
				Name:     u.Name,
				Division: division,
			})
		}
	}

	return Response(c, http.StatusOK, "success get user list", users)
}

func (userHandler) GetDetail(c echo.Context) error {
	type userResponse struct {
		Id       string `json:"id"`
		Name     string `json:"Name"`
		Division string `json:"Division"`
		Role     string `json:"Role"`
	}

	var user userResponse
	var isExist bool

loopUser:
	for _, u := range Users {
		if u.Id == c.Param("id") {
			user = userResponse{
				Id: u.Id,
			}

		loopDivision:
			for _, d := range Divisions {
				if d.Id == u.DivisionId {
					user.Division = d.Name
					break loopDivision
				}
			}

		loopRole:
			for _, r := range Roles {
				if u.RoleId == r.Id {
					user.Role = r.Name
					break loopRole
				}
			}

			isExist = true
			break loopUser
		}
	}

	if !isExist {
		return Response(c, http.StatusNotFound, "user not found", map[string]interface{}{})
	}

	return Response(c, http.StatusOK, "success get user detail", user)
}

func (userHandler) Delete(c echo.Context) error {
	for index, u := range Users {
		if u.Id == c.Param("id") {
			Users = append(Users[:index], Users[index+1:]...)
			return Response(c, http.StatusOK, "success delete user")
		}
	}

	return Response(c, http.StatusNotFound, "user not found")
}

func (userHandler) MyProfile(c echo.Context) error {
	type userResponse struct {
		Id       string `json:"id"`
		Name     string `json:"Name"`
		Division string `json:"Division"`
		Role     string `json:"Role"`
	}

	userId, ok := c.Get("claims").(jwt.MapClaims)["UserId"].(string)
	if !ok {
		log.Println("ERROR : Missing claims on context or user id on claims")
		return Response(c, http.StatusUnauthorized, "unauthorized")
	}

	for _, u := range Users {
		if u.Id == userId {
			user := userResponse{
				Id:   u.Id,
				Name: u.Name,
			}

		loopDivisions:
			for _, d := range Divisions {
				if d.Id == u.DivisionId {
					user.Division = d.Name
					break loopDivisions
				}
			}

		loopRole:
			for _, r := range Roles {
				if r.Id == u.RoleId {
					user.Role = r.Name
					break loopRole
				}
			}

			return Response(c, http.StatusOK, "success get my profile", user)
		}
	}

	log.Println("ERROR : User not found")
	return Response(c, http.StatusUnauthorized, "unauthorized")
}

func IsUserAdmin(userId string) bool {
	for _, u := range Users {
		if u.Id == userId {
			for _, r := range Roles {
				if r.Id == u.RoleId {
					if r.Name == RoleType_Admin {
						return true
					} else {
						return false
					}
				}
			}
		}
	}
	return false
}
