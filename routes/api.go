package routes

import (
	"net/http"
	"turbine-api/handlers"
	"turbine-api/helpers"
	"turbine-api/middleware"
	"turbine-api/repositories"
	"turbine-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NewApi() *api {
	return &api{}
}

type api struct{}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func (api) Init(db *gorm.DB) *echo.Echo {
	route := echo.New()
	route.Validator = &CustomValidator{Validator: validator.New()}

	// // init middleware
	middleware := middleware.NewMiddleware()
	applicationJson := middleware.ApplicationJson
	authAdmin := middleware.AuthAdmin
	// allAuth := middleware.AllAuth

	// // handlers
	// divisionHandler := handlers.NewDivisionHandler()
	// authHandler := handlers.NewAuthHandler()
	// userHandler := handlers.NewUserHandler()
	// roleHandler := handlers.NewRoleHandler()

	// Repositories
	// roleRepository := repositories.NewRoleRepository(db)
	divisionRepository := repositories.NewDivisionRepository(db)

	// Services
	// roleService := services.NewRoleService(roleRepository)
	divisionService := services.NewDivisionService(divisionRepository)

	// Handlers
	// roleHandler := handlers.NewRoleHandler(roleService)
	divisionHandler := handlers.NewDivisionHandler(divisionService)

	route.GET("/", func(c echo.Context) error {
		return helpers.Response(c, http.StatusOK, "hello")
	})
	route.GET("/health", func(c echo.Context) error {
		return helpers.Response(c, http.StatusOK, "healthy")
	})

	divisionRouting := route.Group("divisions")
	divisionRouting.GET("/", divisionHandler.GetListMaster)
	divisionRouting.POST("", divisionHandler.Create, applicationJson, authAdmin)
	divisionRouting.PUT("/:id", divisionHandler.Update, applicationJson, authAdmin)
	divisionRouting.DELETE("/:id", divisionHandler.Delete, authAdmin)

	// authRouting := route.Group("auth")
	// authRouting.POST("/register", authHandler.Register, applicationJson)
	// authRouting.POST("/login", authHandler.Login, applicationJson)
	// authRouting.POST("/refresh-token", authHandler.RefreshToken, applicationJson)

	// userRouting := route.Group("/users")
	// userRouting.GET("", userHandler.GetList, authAdmin)
	// userRouting.GET("/:id", userHandler.GetDetail, authAdmin)
	// userRouting.DELETE("/:id", userHandler.Delete, authAdmin)
	// userRouting.GET("/my-profile", userHandler.MyProfile, allAuth)

	// roleRouting := route.Group("/roles", authAdmin)
	// roleRouting.POST("/", roleHandler.Create)
	// roleRouting.PUT("/:id", roleHandler.Create)
	// roleRouting.GET("/", roleHandler.GetListMaster)
	// roleRouting.DELETE("/:id", roleHandler.Delete)

	// route.GET("/reset-all", authHandler.ResetAll, authAdmin)

	return route
}
