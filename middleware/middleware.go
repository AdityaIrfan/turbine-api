package middleware

import (
	"log"
	"net/http"
	"strings"
	"turbine-api/helpers"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func NewMiddleware() *middleware {
	return &middleware{}
}

type middleware struct{}

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

func (middleware) AuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" {
			return helpers.Response(c, http.StatusBadRequest, "missing authorization header")
		}

		tokens := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(tokens) != 2 {
			return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		} else if tokens[0] != "Bearer" {
			return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		token, err := helpers.VerifyToken(tokens[1])
		if err != nil {
			log.Println("ERROR VERIFY TOKEN : " + err.Error())
			return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("claims", claims)

		return next(c)
	}
}

func (middleware) AllAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" {
			return helpers.Response(c, http.StatusBadRequest, "missing authorization header")
		}

		tokens := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(tokens) != 2 {
			return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		} else if tokens[0] != "Bearer" {
			return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		// token, err := helpers.VerifyToken(tokens[1])
		// if err != nil {
		// 	log.Println("ERROR VERIFY TOKEN : " + err.Error())
		// 	return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		// }

		// claims := token.Claims.(jwt.MapClaims)

		// userId, ok := claims["UserId"].(string)
		// if !ok {
		// 	log.Println("ERROR MISSING STRING USER ID")
		// 	return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		// }

		// for _, u := range handlers.Users {
		// 	if u.Id == userId {
		// 		c.Set("claims", claims)

		// 		return next(c)
		// 	} else {
		// 		return helpers.Response(c, http.StatusUnauthorized, "unauthorized")
		// 	}
		// }

		return helpers.Response(c, http.StatusUnauthorized, "unauthorized")
	}
}

func (middleware) AuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" {
			return helpers.Response(c, http.StatusBadRequest, "missing authorization header")
		}

		tokens := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(tokens) != 2 {
			return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		} else if tokens[0] != "Bearer" {
			return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		// token, err := helpers.VerifyToken(tokens[1])
		// if err != nil {
		// 	log.Println("ERROR VERIFY TOKEN : " + err.Error())
		// 	return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		// }

		// claims := token.Claims.(jwt.MapClaims)

		// // userId, ok := claims["UserId"].(string)
		// // if !ok {
		// // 	log.Println("ERROR MISSING STRING USER ID")
		// // 	return helpers.Response(c, http.StatusUnauthorized, "invalid token")
		// // }

		// // for _, u := range handlers.Users {
		// // 	if u.Id == userId {
		// // 		for _, r := range handlers.Roles {
		// // 			if r.Id == u.RoleId {
		// // 				c.Set("claims", claims)

		// // 				return next(c)
		// // 			} else {
		// // 				return helpers.Response(c, http.StatusUnauthorized, "unauthorized")
		// // 			}
		// // 		}
		// // 	}
		// // }

		return helpers.Response(c, http.StatusUnauthorized, "unauthorized")
	}
}

func (middleware) Signature(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
