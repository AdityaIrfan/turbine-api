package models

import (
	"encoding/json"
	"errors"

	"github.com/phuslu/log"
	"gorm.io/datatypes"
)

type ConfigType uint8

const (
	ConfigType_RootLocation = 1
)

type Config struct {
	Type ConfigType     `gorm:"column:type"`
	Data datatypes.JSON `gorm:"column:data"`
}

func (c *Config) IsEmpty() bool {
	return c == nil
}

func (c *Config) ToConfigRootLocation() *ConfigRootLocation {
	var rootLocation *ConfigRootLocation

	if err := json.Unmarshal(c.Data, &rootLocation); err != nil {
		log.Error().Err(errors.New("ERROR UNMARSHAL ROOT LOCATION FROM CONFIG VALUE : " + err.Error())).Msg("")
		return &ConfigRootLocation{}
	}

	return rootLocation
}

type CoverageAreaType string

const (
	CoverageAreaType_Kilometer = "kilometer"
)

type ConfigRootLocation struct {
	Long             float64          `gorm:"column:long" json:"Long" form:"Long" validate:"required"`
	Lat              float64          `gorm:"column:lat" json:"Lat" form:"Lat" validate:"required"`
	CoverageArea     float64          `gorm:"column:coverage_area" json:"CoverageArea" form:"CoverageArea" validate:"required"`
	CoverageAreaType CoverageAreaType `gorm:"column:coverage_area_type" json:"CoverageAreaType" form:"CoverageAreaType" validate:"required"`
}

func (c *ConfigRootLocation) IsEmpty() bool {
	return c == nil
}

func (c *ConfigRootLocation) ToConfigModel() *Config {
	resJson, err := json.Marshal(c)
	if err != nil {
		log.Error().Err(errors.New("ERROR MARSHAL CONFIG ROOT LOCATION : " + err.Error())).Msg("")
		return &Config{
			Type: ConfigType_RootLocation,
			Data: datatypes.JSON{},
		}
	}

	return &Config{
		Type: ConfigType_RootLocation,
		Data: resJson,
	}
}
