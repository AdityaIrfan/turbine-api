package routes

import (
	"net/http"

	"pln/AdityaIrfan/turbine-api/handlers"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/middleware"
	"pln/AdityaIrfan/turbine-api/repositories"
	"pln/AdityaIrfan/turbine-api/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
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

func (api) Init(db *gorm.DB, client *redis.Client) *echo.Echo {
	route := echo.New()
	route.Validator = &CustomValidator{Validator: validator.New(validator.WithRequiredStructEnabled())}
	// route.Use(middleware.RateLimiterMiddleware)

	// Repositories
	roleRepository := repositories.NewRoleRepository(db)
	divisionRepository := repositories.NewDivisionRepository(db)
	userRepository := repositories.NewUserRepository(db)
	authRedisRepository := repositories.NewAuthRedisRepository(client)
	// configRepository := repositories.NewConfigRepository(db)
	// configRedisRepository := repositories.NewConfigRedisRepository(client)
	pltaRepository := repositories.NewPltaRepository(db)
	turbineRepository := repositories.NewTurbineRepository(db)
	pltaUnitRepository := repositories.NewPltaUnitRepo(db)

	// Services
	// roleService := services.NewRoleService(roleRepository)
	divisionService := services.NewDivisionService(divisionRepository, userRepository)
	userService := services.NewUserService(userRepository, divisionRepository, roleRepository)
	authService := services.NewAuthService(userRepository, authRedisRepository, divisionRepository)
	// configService := services.NewConfigService(configRepository, configRedisRepository, userRepository)
	pltaService := services.NewPltaService(pltaRepository, userRepository)
	turbineService := services.NewTurbineService(turbineRepository, pltaUnitRepository, userRepository)
	pltaUnitService := services.NewPltaUnitService(pltaUnitRepository, pltaRepository)

	// Handlers
	// roleHandler := handlers.NewRoleHandler(roleService)
	divisionHandler := handlers.NewDivisionHandler(divisionService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	// configHandler := handlers.NewConfigHandler(configService)
	turbineHandler := handlers.NewTurbineHandler(turbineService)
	pltaHandler := handlers.NewPltaHandler(pltaService)
	pltaUnitHandler := handlers.NewPltaUnitHandler(pltaUnitService)

	// Middleware
	middleware := middleware.NewMiddleware(authRedisRepository, userRepository)
	// applicationJson := middleware.ApplicationJson
	superAdmin := middleware.AuthSuperAdmin
	authAdmin := middleware.AuthAdmin
	authUser := middleware.AuthUser
	allAuth := middleware.Auth
	// signature := middleware.Signature

	route.GET("/", func(c echo.Context) error {
		return helpers.Response(c, http.StatusOK, "hello")
	})
	route.GET("/health", func(c echo.Context) error {
		return helpers.Response(c, http.StatusOK, "healthy")
	})

	v1 := route.Group("/v1")

	v1_DivisionRouting := v1.Group("/divisions")
	v1_DivisionRouting.GET("/master", divisionHandler.GetListMaster)
	v1_DivisionRouting.GET("", divisionHandler.GetListWithPaginate, authAdmin)
	v1_DivisionRouting.POST("", divisionHandler.Create /*, applicationJson*/, authAdmin)
	v1_DivisionRouting.PUT("/:id", divisionHandler.Update /*, applicationJson*/, authAdmin)
	v1_DivisionRouting.DELETE("/:id", divisionHandler.Delete, authAdmin)

	v1_AuthRouting := v1.Group("/auth")
	v1_AuthRouting.POST("/register", authHandler.Register /*, applicationJson*/)
	v1_AuthRouting.POST("/login", authHandler.Login /*, applicationJson*/)
	v1_AuthRouting.POST("/refresh-token", authHandler.RefreshToken /*, applicationJson*/)
	v1_AuthRouting.POST("/logout", authHandler.Logout, allAuth)

	// ADMIN BY SUPER ADMIN
	// v1_adminRoutingBySuperAdmin := v1.Group("/super/users")
	// v1_adminRoutingBySuperAdmin.PUT("/:id", userHandler.UpdateUserByAdmin, authAdmin /*, applicationJson*/)
	// v1_adminRoutingBySuperAdmin.DELETE("/:id", userHandler.DeleteUserByAdmin, authAdmin)
	// v1_adminRoutingBySuperAdmin.GET("", userHandler.GetListUserWithPaginateByAdmin, authAdmin)
	// v1_adminRoutingBySuperAdmin.GET("/:id", userHandler.GetDetailUserByAdmin, authAdmin)

	// USER BY ADMIN
	v1_UserRoutingByAdmin := v1.Group("/admin/users")
	v1_UserRoutingByAdmin.POST("", userHandler.CreateUserByAdmin, authAdmin /*, applicationJson*/)
	v1_UserRoutingByAdmin.PUT("/:id", userHandler.UpdateUserByAdmin, authAdmin /*, applicationJson*/)
	v1_UserRoutingByAdmin.GET("/:id", userHandler.GetDetailUserByAdmin, authAdmin)
	v1_UserRoutingByAdmin.DELETE("/:id", userHandler.DeleteUserByAdmin, authAdmin)
	v1_UserRoutingByAdmin.GET("", userHandler.GetListUserWithPaginateByAdmin, authAdmin)
	v1_UserRoutingByAdmin.POST("/generate-password/:id", userHandler.GenerateUserPasswordByAdmin, authAdmin /*, applicationJson*/)

	// USER ITSELF
	v1_UserRouting := v1.Group("/my")
	v1_UserRouting.PUT("", userHandler.UpdateMyProfile, authUser)
	v1_UserRouting.GET("", userHandler.GetMyProfile, allAuth)
	v1_UserRouting.POST("/change-password", userHandler.ChangeMyPassword, authUser /*, applicationJson*/)

	// roleRouting := v1.Group("/roles", authAdmin)
	// roleRouting.POST("/", roleHandler.Create/*, applicationJson*/)
	// roleRouting.PUT("/:id", roleHandler.Create/*, applicationJson*/)
	// roleRouting.GET("/", roleHandler.GetListMaster)
	// roleRouting.DELETE("/:id", roleHandler.Delete)

	// v1_ConfigRouting := v1.Group("/configs")
	// v1_ConfigRouting.POST("/root-location", configHandler.SaveOrUpdate, authAdmin /*, applicationJson*/)
	// v1_ConfigRouting.GET("/root-location", configHandler.GetRootLocation, authAdmin)

	v1_PltaRouting := v1.Group("/plta")
	v1_PltaRouting.POST("", pltaHandler.Create, superAdmin /*, applicationJson*/)
	v1_PltaRouting.PUT("/:id", pltaHandler.Update, superAdmin /*, applicationJson*/)
	v1_PltaRouting.GET("/:id", pltaHandler.Detail, authAdmin /*, applicationJson*/)
	v1_PltaRouting.GET("/master", pltaHandler.GetListMaster, allAuth)
	v1_PltaRouting.DELETE("/:id", pltaHandler.Delete, allAuth)
	v1_PltaRouting.GET("", pltaHandler.GetListWithPaginate, authAdmin)

	v1_TurbineRouting := v1.Group("/turbines")
	v1_TurbineRouting.POST("", turbineHandler.Create, allAuth /*, applicationJson*/)
	v1_TurbineRouting.GET("/:id", turbineHandler.GetDetail, allAuth)
	v1_TurbineRouting.GET("", turbineHandler.GetList, allAuth)
	v1_TurbineRouting.GET("/latest", turbineHandler.GetLatest, allAuth)
	v1_TurbineRouting.DELETE("/:id", turbineHandler.Delete, allAuth)

	v1_PltaUnitRouting := v1.Group("/plta-unit")
	v1_PltaUnitRouting.PUT("/:id", pltaUnitHandler.CreateOrUpdate, authAdmin)
	v1_PltaUnitRouting.DELETE("/:id", pltaUnitHandler.Delete, authAdmin)

	return route
}
