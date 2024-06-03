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
	// var long, lat, coverageArea float64
	// var coverageType CoverageAreaType

	// if value, ok := c.Data["long"].(float64); !ok {
	// 	log.Error().Err(errors.New("CONFIG DATA long IS EMPTY OR IS NOT FLOAT64, CEK AGAIN ON CONFIG ROOT LOCATION DATABASE")).Msg("")
	// } else {
	// 	long = value
	// }

	// if value, ok := c.Data["lat"].(float64); !ok {
	// 	log.Error().Err(errors.New("CONFIG DATA lat IS EMPTY OR IS NOT FLOAT64, CEK AGAIN ON CONFIG ROOT LOCATION DATABASE")).Msg("")
	// } else {
	// 	lat = value
	// }

	// if value, ok := c.Data["coverage_area"].(float64); !ok {
	// 	log.Error().Err(errors.New("CONFIG DATA coverage_area IS EMPTY OR IS NOT FLOAT64, CEK AGAIN ON CONFIG ROOT LOCATION DATABASE")).Msg("")
	// } else {
	// 	coverageArea = value
	// }

	// if value, ok := c.Data["coverage_area_type"].(string); !ok {
	// 	log.Error().Err(errors.New("CONFIG DATA coverage_area_type IS EMPTY OR IS NOT FLOAT64, CEK AGAIN ON CONFIG ROOT LOCATION DATABASE")).Msg("")
	// } else {
	// 	coverageType = CoverageAreaType(value)
	// }

	// return &ConfigRootLocation{
	// 	Long:             long,
	// 	Lat:              lat,
	// 	CoverageArea:     coverageArea,
	// 	CoverageAreaType: coverageType,
	// }

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
	Long             float64          `gorm:"column:long" json:"Long"`
	Lat              float64          `gorm:"column:lat" json:"Lat"`
	CoverageArea     float64          `gorm:"column:coverage_area" json:"CoverageArea"`
	CoverageAreaType CoverageAreaType `gorm:"column:coverage_area_type" json:"CoverageAreaType"`
}

func (c *ConfigRootLocation) IsEmpty() bool {
	return c == nil
}
