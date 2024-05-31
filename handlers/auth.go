package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	helpers "turbine-api/helpers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
)

var RefreshTokenMap = map[string]interface{}{}

func InitRefreshToken() {
	RefreshTokenMap = make(map[string]interface{})

	fmt.Println("SUCCESS INIT REFRESH TOKEN")
}

type authHandler struct{}

func NewAuthHandler() *authHandler {
	return &authHandler{}
}

func (a *authHandler) Register(c echo.Context) error {
	if len(Users) == 3 {
		return Response(c, http.StatusBadRequest, "this is just mock api, do not try too much")
	}

	type registerPayload struct {
		Name       string `json:"Name" validate:"required"`
		DivisionId string `json:"DivisionId" validate:"required"`
		RoleId     string `json:"RoleId"`
		Username   string `json:"Username" validate:"required"`
		Password   string `json:"Password" validate:"required"`
	}

	var payload = new(registerPayload)

	if err := c.Bind(payload); err != nil {
		return Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return Response(c, http.StatusBadRequest, errMessage)
	}

	for _, u := range Users {
		if strings.EqualFold(u.Username, payload.Username) {
			return Response(c, http.StatusBadRequest, "username already in use")
		}
	}

	isDivisionExist := false
	var divisionName string
	for _, division := range Divisions {
		if payload.DivisionId == division.Id {
			isDivisionExist = true
			divisionName = division.Name
		}
	}

	if !isDivisionExist {
		return Response(c, http.StatusBadRequest, "division not found")
	}

	if payload.RoleId != "" {
	loopRoles1:
		for _, r := range Roles {
			if r.Id == payload.RoleId {
				payload.RoleId = r.Id
				break loopRoles1
			}
		}

		if payload.RoleId == "" {
			return Response(c, http.StatusBadRequest, "role not found")
		}
	} else {
	loopRoles2:
		for _, r := range Roles {
			if r.Name == RoleType_User {
				payload.RoleId = r.Id
				break loopRoles2
			}
		}
	}

	user := User{
		Id:         ulid.Make().String(),
		Name:       payload.Name,
		Username:   payload.Username,
		DivisionId: payload.DivisionId,
		Password:   payload.Password,
		RoleId:     payload.RoleId,
	}

	token, timeRefreshTokenActive, err := helpers.GenerateToken(user.Id)
	if err != nil {
		log.Println("FAILED TO GENERATE TOKEN : ", err.Error())
		return Response(c, http.StatusUnprocessableEntity, "register failed")
	}

	refreshToken, err := helpers.GenerateRefreshToken(user.Id)
	if err != nil {
		log.Println("FAILED TO GENERATE REFRESH TOKEN : ", err.Error())
		return Response(c, http.StatusUnprocessableEntity, "register failed")
	}

	RefreshTokenMap[user.Id] = map[string]interface{}{
		"RefreshToken": refreshToken,
		"StartActive":  timeRefreshTokenActive,
	}

	Users = append(Users, user)

	return Response(c, http.StatusOK, "success register", map[string]interface{}{
		"Name":         user.Name,
		"Division":     divisionName,
		"Token":        token,
		"RefreshToken": refreshToken,
	})
}

func (a *authHandler) Login(c echo.Context) error {
	type loginPayload struct {
		Username string `json:"Username" validate:"required"`
		Password string `json:"Password" validate:"required"`
	}

	var payload = new(loginPayload)

	if err := c.Bind(payload); err != nil {
		return Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return Response(c, http.StatusBadRequest, errMessage)
	}

	for _, u := range Users {
		if u.Username == payload.Username && u.Password == payload.Password {
			var division string

			for _, d := range Divisions {
				if d.Id == u.DivisionId {
					division = d.Name
				}
			}

			token, timeRefreshTokenActive, err := helpers.GenerateToken(u.Id)
			if err != nil {
				log.Println("FAILED TO GENERATE TOKEN : ", err.Error())
				return Response(c, http.StatusUnprocessableEntity, "login failed")
			}

			refreshToken, err := helpers.GenerateRefreshToken(u.Id)
			if err != nil {
				log.Println("FAILED TO GENERATE REFRESH TOKEN : ", err.Error())
				return Response(c, http.StatusUnprocessableEntity, "login failed")
			}

			RefreshTokenMap[u.Id] = map[string]interface{}{
				"RefreshToken": refreshToken,
				"StartActive":  timeRefreshTokenActive,
			}

			return Response(c, http.StatusOK, "success login", map[string]interface{}{
				"Name":         u.Name,
				"Division":     division,
				"Token":        token,
				"RefreshToken": refreshToken,
			})
		}
	}

	return Response(c, http.StatusBadRequest, "credential not found")
}

func (authHandler) RefreshToken(c echo.Context) error {
	type refreshTokenPayload struct {
		RefreshToken string `json:"RefreshToken" validate:"required"`
	}

	var payload = new(refreshTokenPayload)

	if err := c.Bind(payload); err != nil {
		return Response(c, http.StatusBadRequest, "error binding body request")
	}

	if err := c.Validate(payload); err != nil {
		errMessage := helpers.GenerateValidationErrorMessage(err)
		return Response(c, http.StatusBadRequest, errMessage)
	}

	token, err := helpers.VerifyRefreshToken(payload.RefreshToken)
	if err != nil {
		return Response(c, http.StatusBadRequest, "invalid refresh token")
	} else if !token.Valid {
		return Response(c, http.StatusBadRequest, "invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)
	userId, ok := claims["UserId"].(string)
	if !ok {
		log.Println("ERROR : THERE'S NO USER ID ON CLAIMS ")
		return Response(c, http.StatusBadRequest, "invalid refresh token")
	}

	var user User
	var division Division
loopUser:
	for _, u := range Users {
		if u.Id == userId {
		loopDivision:
			for _, d := range Divisions {
				if d.Id == u.DivisionId {
					division = d
					break loopDivision
				}
			}

			user = u
			break loopUser
		}
	}

	if helpers.IsTokenExpired(token) {
		delete(RefreshTokenMap, userId)

		return Response(c, http.StatusBadRequest, "refresh token expired")
	}

	if refresh, ok := RefreshTokenMap[userId]; ok {
		if time.Now().Before(time.Unix(refresh.(map[string]interface{})["StartActive"].(int64), 0)) {
			return Response(c, http.StatusBadRequest, "token still active")
		}
	} else {
		log.Println("ERROR : REFRESH TOKEN DATA IS NOT STORED ON RefreshTokenMap")
		return Response(c, http.StatusBadRequest, "invalid refresh token")
	}

	newToken, timeRefreshTokenActive, err := helpers.GenerateToken(userId)
	if err != nil {
		log.Println("FAILED TO GENERATE TOKEN : ", err.Error())
		return Response(c, http.StatusUnprocessableEntity, "login failed")
	}

	newRefreshToken, err := helpers.GenerateRefreshToken(userId)
	if err != nil {
		log.Println("FAILED TO GENERATE REFRESH TOKEN : ", err.Error())
		return Response(c, http.StatusUnprocessableEntity, "login failed")
	}

	RefreshTokenMap[userId] = map[string]interface{}{
		"RefreshToken": newRefreshToken,
		"StartActive":  timeRefreshTokenActive,
	}

	return Response(c, http.StatusOK, "success refresh token", map[string]interface{}{
		"Name":         user.Name,
		"Division":     division.Name,
		"Token":        newToken,
		"RefreshToken": newRefreshToken,
	})
}

func (authHandler) ResetAll(c echo.Context) error {
	InitDivisions()
	InitRole()
	InitRefreshToken()
	InitUsers()
	fmt.Println("===== RESET ALL ======")
	return Response(c, http.StatusOK, "success reset all data")
}
