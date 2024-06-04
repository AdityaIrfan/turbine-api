package services

import (
	"errors"
	"net/http"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/labstack/echo/v4"
	"github.com/phuslu/log"
)

type configService struct {
	configRepo      contract.IConfigRepository
	configRedisRepo contract.IConfigRedisRepository
	userRepo        contract.IUserRepository
}

func NewConfigService(
	configRepo contract.IConfigRepository,
	configRedisRepo contract.IConfigRedisRepository,
	userRepo contract.IUserRepository) contract.IConfigService {
	return &configService{
		configRepo:      configRepo,
		configRedisRepo: configRedisRepo,
		userRepo:        userRepo,
	}
}

func (cs *configService) SaveOrUpdate(c echo.Context, in *models.ConfigRootLocation) error {
	if err := cs.configRepo.SaveOrUpdateRootLocation(in); err != nil {
		return helpers.ResponseUnprocessableEntity(c)
	}

	go cs.configRedisRepo.SaveRootLocation(in)
	return helpers.Response(c, http.StatusOK, "success create or update config root location")
}

func (cs *configService) GetRootLocation(c echo.Context) error {
	rootLocation, err := cs.configRedisRepo.GetRootLocation()
	if err != nil || rootLocation.IsEmpty() {
		config, err := cs.configRepo.GetByType(models.ConfigType_RootLocation)
		if err != nil {
			return helpers.ResponseUnprocessableEntity(c)
		} else if config.IsEmpty() {
			log.Error().Err(errors.New("ERROR CONFIG ROOT LOCATION IS EMPTY ON DATABASE")).Msg("")
			return helpers.Response(c, http.StatusNotFound, "not found, call developer immediately")
		}

		rootLocation := config.ToConfigRootLocation()
		go cs.configRedisRepo.SaveRootLocation(rootLocation)

		return helpers.Response(c, http.StatusOK, "success get config location", rootLocation)
	}

	return helpers.Response(c, http.StatusOK, "success get config location", rootLocation)
}
