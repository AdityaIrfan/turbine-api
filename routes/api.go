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

	// Repositories
	roleRepository := repositories.NewRoleRepository(db)
	divisionRepository := repositories.NewDivisionRepository(db)
	userRepository := repositories.NewUserRepository(db)
	authRedisRepository := repositories.NewAuthRedisRepository(client)
	configRepository := repositories.NewConfigRepository(db)
	configRedisRepository := repositories.NewConfigRedisRepository(client)
	towerRepository := repositories.NewTowerRepository(db)
	turbineRepository := repositories.NewTurbineRepository(db)

	// Services
	// roleService := services.NewRoleService(roleRepository)
	divisionService := services.NewDivisionService(divisionRepository, userRepository)
	userService := services.NewUserService(userRepository, divisionRepository, roleRepository)
	authService := services.NewAuthService(userRepository, authRedisRepository, divisionRepository)
	configService := services.NewConfigService(configRepository, configRedisRepository, userRepository)
	towerService := services.NewTowerService(towerRepository)
	turbineService := services.NewTurbineService(turbineRepository, towerRepository)

	// Handlers
	// roleHandler := handlers.NewRoleHandler(roleService)
	divisionHandler := handlers.NewDivisionHandler(divisionService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	configHandler := handlers.NewConfigHandler(configService)
	turbineHandler := handlers.NewTurbineHandler(turbineService)
	towerHandler := handlers.NewTowerHandler(towerService)

	// Middleware
	middleware := middleware.NewMiddleware(authRedisRepository, userRepository)
	// applicationJson := middleware.ApplicationJson
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

	divisionRouting := route.Group("divisions")
	divisionRouting.GET("/master", divisionHandler.GetListMaster)
	divisionRouting.GET("", divisionHandler.GetListWithPaginate, authAdmin)
	divisionRouting.POST("", divisionHandler.Create /*, applicationJson*/, authAdmin)
	divisionRouting.PUT("/:id", divisionHandler.Update /*, applicationJson*/, authAdmin)
	divisionRouting.DELETE("/:id", divisionHandler.Delete, authAdmin)

	authRouting := route.Group("auth")
	authRouting.POST("/register", authHandler.Register /*, applicationJson*/)
	authRouting.POST("/login", authHandler.Login /*, applicationJson*/)
	authRouting.POST("/refresh-token", authHandler.RefreshToken /*, applicationJson*/)
	authRouting.POST("/logout", authHandler.Logout, allAuth)

	userRoutingByAdmin := route.Group("/admin/users")
	userRoutingByAdmin.POST("/", userHandler.CreateUserAdminByAdmin, authAdmin /*, applicationJson*/)
	userRoutingByAdmin.PUT("/:id", userHandler.UpdateByAdmin, authAdmin /*, applicationJson*/)
	userRoutingByAdmin.GET("/:id", userHandler.GetDetailByAdmin, authAdmin)
	userRoutingByAdmin.DELETE("/:id", userHandler.DeleteByAdmin, authAdmin)
	userRoutingByAdmin.GET("/", userHandler.GetListWithPaginateByAdmin, authAdmin)
	userRoutingByAdmin.POST("/generate-password", userHandler.GeneratePasswordByAdmin, authAdmin /*, applicationJson*/)

	userRouting := route.Group("/my")
	userRouting.PUT("/:id", userHandler.Update, authUser)
	userRouting.GET("", userHandler.GetMyProfile, authUser)
	userRouting.POST("/change-password", userHandler.ChangePassword, authUser /*, applicationJson*/)

	// roleRouting := route.Group("/roles", authAdmin)
	// roleRouting.POST("/", roleHandler.Create/*, applicationJson*/)
	// roleRouting.PUT("/:id", roleHandler.Create/*, applicationJson*/)
	// roleRouting.GET("/", roleHandler.GetListMaster)
	// roleRouting.DELETE("/:id", roleHandler.Delete)

	configRouting := route.Group("/configs")
	configRouting.POST("/root-location", configHandler.SaveOrUpdate, authAdmin /*, applicationJson*/)
	configRouting.GET("/root-location", configHandler.GetRootLocation, authAdmin)

	towerRouting := route.Group("/towers")
	towerRouting.POST("", towerHandler.Create, authAdmin /*, applicationJson*/)
	towerRouting.PUT("/:id", towerHandler.Update, authAdmin /*, applicationJson*/)
	towerRouting.GET("/master", towerHandler.GetListMaster, allAuth)
	towerRouting.DELETE("/:id", towerHandler.Delete, authAdmin)

	turbineRouting := route.Group("/turbines")
	turbineRouting.POST("", turbineHandler.Create, allAuth /*, applicationJson*/)
	turbineRouting.GET("/:id", turbineHandler.GetDetail, allAuth)
	turbineRouting.GET("", turbineHandler.GetList, allAuth)
	turbineRouting.GET("/latest", turbineHandler.GetLatest, allAuth)

	return route
}
