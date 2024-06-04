package middleware

import (
	"log"
	"net/http"
	"strings"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func NewMiddleware(authRedisRepo contract.IAuthRedisRepository, userRepo contract.IUserRepository) *middleware {
	return &middleware{
		authRedisRepo: authRedisRepo,
		userRepo:      userRepo,
	}
}

type middleware struct {
	authRedisRepo contract.IAuthRedisRepository
	userRepo      contract.IUserRepository
}

func (middleware) ApplicationJson(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") == "" {
			return helpers.Response(c, http.StatusBadRequest, "missing content-type header application/json")
		}

		if c.Request().Header.Get("Content-Type") != "application/json" {
			return helpers.Response(c, http.StatusBadRequest, "content-type header not allowed")
		}

		return next(c)
	}
}

func (m *middleware) AuthSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := m.checkToken(c, models.UserRole_SuperAdmin)
		if err != nil {
			return err
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		return next(c)
	}
}

func (m *middleware) AuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := m.checkToken(c, models.UserRole_Admin)
		if err != nil {
			return err
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		return next(c)
	}
}

func (m *middleware) AuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := m.checkToken(c, models.UserRole_User)
		if err != nil {
			return err
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		return next(c)
	}
}

func (m *middleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := m.checkToken(c)
		if err != nil {
			return err
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		return next(c)
	}
}

func (m *middleware) checkToken(c echo.Context, roles ...models.UserRole) (*jwt.Token, error) {
	if c.Request().Header.Get("Authorization") == "" {
		return nil, helpers.Response(c, http.StatusBadRequest, "missing authorization header")
	}

	tokens := strings.Split(c.Request().Header.Get("Authorization"), " ")
	if len(tokens) != 2 {
		return nil, helpers.Response(c, http.StatusUnauthorized, "invalid token")
	} else if tokens[0] != "Bearer" {
		return nil, helpers.Response(c, http.StatusUnauthorized, "invalid token")
	}

	tokenString := tokens[1]

	token, err := helpers.VerifyToken(tokenString)
	if err != nil {
		log.Println("ERROR VERIFY TOKEN : " + err.Error())
		return nil, helpers.Response(c, http.StatusUnauthorized, "invalid token")
	}

	userId, ok := token.Claims.(jwt.MapClaims)["Id"].(string)
	if !ok || userId == "" {
		return nil, helpers.ResponseForbiddenAccess(c)
	}

	existingToken, err := m.authRedisRepo.GetToken(userId)
	if err != nil {
		return nil, helpers.ResponseUnprocessableEntity(c)
	} else if existingToken == "" || tokenString != existingToken {
		return nil, helpers.Response(c, http.StatusUnauthorized, "invalid token")
	}

	user, err := m.userRepo.GetByIdWithSelectedFields(userId, "id, status, role")
	if err != nil {
		return nil, helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() || !user.IsActive() {
		return nil, helpers.ResponseForbiddenAccess(c)
	}

	for _, role := range roles {
		if role == user.Role {
			return token, nil
		}
	}

	return nil, helpers.ResponseForbiddenAccess(c)
}

func (middleware) Signature(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
