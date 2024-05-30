package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	helpers "turbine-backend-api-contract/halpers"
	"turbine-backend-api-contract/handlers"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func main() {
	port := 8081
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Panic Detected : " + fmt.Sprint(err))
		}
	}()
	route := echo.New()
	route.Validator = &CustomValidator{Validator: validator.New()}

	handlers.InitDivisions()
	handlers.InitRole()
	handlers.InitRefreshToken()
	handlers.InitUsers()

	// handle validator
	header := &header{}

	// handlers
	divisionHandler := handlers.NewDivisionHandler()
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()
	roleHandler := handlers.NewRoleHandler()

	route.GET("/", func(c echo.Context) error {
		return handlers.Response(c, http.StatusOK, "hello")
	})
	route.GET("/health", func(c echo.Context) error {
		return handlers.Response(c, http.StatusOK, "healthy")
	})

	applicationJson := header.ContentType.ApplicationJson
	authAdmin := header.Authorization.AuthAdmin
	allAuth := header.Authorization.AllAuth

	divisionRouting := route.Group("divisions")
	divisionRouting.GET("/master/list", divisionHandler.GetListMasterData)
	divisionRouting.POST("", divisionHandler.Add, applicationJson, authAdmin)
	divisionRouting.PUT("/:id", divisionHandler.Update, applicationJson, authAdmin)
	divisionRouting.DELETE("/:id", divisionHandler.Delete, authAdmin)

	authRouting := route.Group("auth")
	authRouting.POST("/register", authHandler.Register, applicationJson)
	authRouting.POST("/login", authHandler.Login, applicationJson)
	authRouting.POST("/refresh-token", authHandler.RefreshToken, applicationJson)

	userRouting := route.Group("/users")
	userRouting.GET("", userHandler.GetList, authAdmin)
	userRouting.GET("/:id", userHandler.GetDetail, authAdmin)
	userRouting.DELETE("/:id", userHandler.Delete, authAdmin)
	userRouting.GET("/my-profile", userHandler.MyProfile, allAuth)

	roleRouting := route.Group("/roles", authAdmin)
	roleRouting.GET("/master/list", roleHandler.GetListMasterData)

	route.GET("/reset-all", authHandler.ResetAll, authAdmin)

	fmt.Println("This server will run on localhost:", strconv.Itoa(port))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), route); err != nil {
		log.Println("FAILED TO RUN SERVER : " + err.Error())
	}
}

type header struct {
	ContentType   contentType
	Authorization authorization
}

type contentType struct{}

func (contentType) ApplicationJson(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Content-Type") == "" {
			return handlers.Response(c, http.StatusBadRequest, "missing content-type header")
		}

		if c.Request().Header.Get("Content-Type") != "application/json" {
			return handlers.Response(c, http.StatusBadRequest, "content-type header not allowed")
		}

		return next(c)
	}
}

type authorization struct{}

func (authorization) AuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" {
			return handlers.Response(c, http.StatusBadRequest, "missing authorization header")
		}

		tokens := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(tokens) != 2 {
			return handlers.Response(c, http.StatusUnauthorized, "invalid token")
		} else if tokens[0] != "Bearer" {
			return handlers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		token, err := helpers.VerifyToken(tokens[1])
		if err != nil {
			log.Println("ERROR VERIFY TOKEN : " + err.Error())
			return handlers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)

		userId, ok := claims["UserId"].(string)
		if !ok {
			log.Println("ERROR MISSING STRING USER ID")
			return handlers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		for _, u := range handlers.Users {
			if u.Id == userId {
				for _, r := range handlers.Roles {
					if r.Id == u.RoleId {
						c.Set("claims", claims)

						return next(c)
					} else {
						return handlers.Response(c, http.StatusUnauthorized, "unauthorized")
					}
				}
			}
		}

		return handlers.Response(c, http.StatusUnauthorized, "unauthorized")
	}
}

func (authorization) AllAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" {
			return handlers.Response(c, http.StatusBadRequest, "missing authorization header")
		}

		tokens := strings.Split(c.Request().Header.Get("Authorization"), " ")
		if len(tokens) != 2 {
			return handlers.Response(c, http.StatusUnauthorized, "invalid token")
		} else if tokens[0] != "Bearer" {
			return handlers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		token, err := helpers.VerifyToken(tokens[1])
		if err != nil {
			log.Println("ERROR VERIFY TOKEN : " + err.Error())
			return handlers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)

		userId, ok := claims["UserId"].(string)
		if !ok {
			log.Println("ERROR MISSING STRING USER ID")
			return handlers.Response(c, http.StatusUnauthorized, "invalid token")
		}

		for _, u := range handlers.Users {
			if u.Id == userId {
				c.Set("claims", claims)

				return next(c)
			} else {
				return handlers.Response(c, http.StatusUnauthorized, "unauthorized")
			}
		}

		return handlers.Response(c, http.StatusUnauthorized, "unauthorized")
	}
}
