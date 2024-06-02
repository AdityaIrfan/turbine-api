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
}

func NewConfigService(configRepo contract.IConfigRepository, configRedisRepo contract.IConfigRedisRepository) contract.IConfigService {
	return &configService{
		configRepo:      configRepo,
		configRedisRepo: configRedisRepo,
	}
}

func (cs *configService) GetRootLocation(c echo.Context) error {
	rootLocation, err := cs.configRedisRepo.GetRootLocation()
	if err != nil || rootLocation.IsEmpty() {
		config, err := cs.configRepo.GetByType(models.ConfigType_RootLocation)
		if err != nil {
			return helpers.Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
		} else if config.IsEmpty() {
			log.Error().Err(errors.New("ERROR CONFIG ROOT LOCATION IS EMPTY ON DATABASE"))
			return helpers.Response(c, http.StatusNotFound, "not found, call developer immediately")
		}

		rootLocation := config.ToConfigRootLocation()
		go cs.configRedisRepo.SaveRootLocation(rootLocation)

		return helpers.Response(c, http.StatusOK, "success get config location", rootLocation)
	}

	return helpers.Response(c, http.StatusOK, "success get config location", rootLocation)
}
