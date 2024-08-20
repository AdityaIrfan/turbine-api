package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	contract "pln/AdityaIrfan/turbine-api/contracts"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/phuslu/log"
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
			return middlewareErrorResponse(c, err)
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		return next(c)
	}
}

func (m *middleware) AuthAdminAndSuperAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := m.checkToken(c, models.UserRole_SuperAdmin, models.UserRole_Admin)
		if err != nil {
			return middlewareErrorResponse(c, err)
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
			return middlewareErrorResponse(c, err)
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
			return middlewareErrorResponse(c, err)
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		return next(c)
	}
}

func (m *middleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := m.checkToken(c, models.UserRole_SuperAdmin, models.UserRole_Admin, models.UserRole_User)
		if err != nil {
			return middlewareErrorResponse(c, err)
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		return next(c)
	}
}

func middlewareErrorResponse(c echo.Context, err error) error {
	switch err.Error() {
	case "missing authorization token":
		return helpers.Response(c, http.StatusBadRequest, "missing authorization header")
	case "invalid token":
		return helpers.Response(c, http.StatusUnauthorized, "invalid token")
	case "forbidden access":
		return helpers.Response(c, http.StatusUnauthorized, "forbidden access")
	case "expired token":
		return helpers.Response(c, http.StatusUnauthorized, "expired token")
	default:
		return helpers.Response(c, http.StatusUnauthorized, "forbidden access")
	}
}

func (m *middleware) checkToken(c echo.Context, roles ...models.UserRole) (*jwt.Token, error) {
	if c.Request().Header.Get("Authorization") == "" {
		return nil, errors.New("missing authorization token")
	}

	tokens := strings.Split(c.Request().Header.Get("Authorization"), " ")
	if len(tokens) != 2 {
		return nil, errors.New("invalid token")
	} else if tokens[0] != "Bearer" {
		return nil, errors.New("invalid token")
	}

	tokenString := tokens[1]

	token, err := helpers.VerifyToken(tokenString)
	if err != nil {
		log.Error().Err(errors.New("ERROR VERIFY TOKEN : " + err.Error())).Msg("")
		return nil, errors.New("invalid token")
	}

	userId, ok := token.Claims.(jwt.MapClaims)["Id"].(string)
	if !ok || userId == "" {
		log.Error().Err(errors.New("THIS TOKEN DOES NOT HAVE [Id] IN CLAIMS : ")).Msg("")
		return nil, errors.New("forbidden access")
	}

	expUnix, ok := token.Claims.(jwt.MapClaims)["Exp"].(float64)
	if !ok || userId == "" {
		log.Error().Err(errors.New("THIS TOKEN DOES NOT HAVE [Exp] IN CLAIMS : ")).Msg("")
		return nil, errors.New("forbidden access")
	}
	expirationTime := time.Unix(int64(expUnix), 0)
	if time.Now().After(expirationTime) {
		return nil, errors.New("expired token")
	}

	// existingToken, err := m.authRedisRepo.GetToken(userId)
	// if err != nil {
	// 	return nil, helpers.ResponseUnprocessableEntity(c)
	// } else if existingToken == "" || tokenString != existingToken {
	// 	return nil, errors.New("invalid token")
	// }

	user, err := m.userRepo.GetByIdWithSelectedFields(userId, "id, status, role")
	if err != nil {
		return nil, helpers.ResponseUnprocessableEntity(c)
	} else if user.IsEmpty() || !user.IsActive() {
		return nil, errors.New("forbidden access")
	}

	for _, role := range roles {
		switch role {
		case models.UserRole_SuperAdmin:
			if user.Role == role {
				return token, nil
			}
		case models.UserRole_Admin:
			if user.Role == role {
				return token, nil
			}
		case models.UserRole_User:
			if user.Role == role {
				return token, nil
			}
		}
	}

	return nil, errors.New("forbidden access")
}

func (middleware) Signature(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
