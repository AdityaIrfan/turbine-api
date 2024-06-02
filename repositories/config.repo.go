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

func (c *configRepository) GetByType(configType models.ConfigType) (*models.Config, error) {
	var config *models.Config

	if err := c.db.Where("type = ?", configType).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error().Err(errors.New("ERROR QUERY CONFIG BY TYPE : " + err.Error()))
		return nil, err
	}

	return config, nil
}
