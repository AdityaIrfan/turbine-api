package repositories

import (
	"errors"
	contract "turbine-api/contracts"
	"turbine-api/models"

	"github.com/phuslu/log"
	"gorm.io/gorm"
)

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) contract.IConfigRepository {
	return &configRepository{
		db: db,
	}
}

func (c *configRepository) SaveOrUpdateRootLocation(rootLocation *models.ConfigRootLocation) error {
	config, err := c.GetByType(models.ConfigType_RootLocation)
	if err != nil {
		return err
	} else if config.IsEmpty() {
		config = rootLocation.ToConfigModel()
		if err := c.db.Create(&config).Error; err != nil {
			log.Error().Err(errors.New("ERROR CONFIG ROOT LOCATION CREATING : " + err.Error())).Msg("")
			return err
		}
	}

	config = rootLocation.ToConfigModel()
	if err := c.db.Updates(&config).Error; err != nil {
		log.Error().Err(errors.New("ERROR CONFIG ROOT LOCATION UPDATING : " + err.Error())).Msg("")
		return err
	}

	return nil
}

func (c *configRepository) GetByType(configType models.ConfigType) (*models.Config, error) {
	var config *models.Config

	if err := c.db.Where("type = ?", configType).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY CONFIG BY TYPE : " + err.Error())).Msg("")
		return nil, err
	}

	return config, nil
}
