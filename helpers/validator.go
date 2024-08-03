package helpers

import (
	"errors"
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/go-playground/validator/v10"
	geoIp "github.com/oschwald/geoip2-golang"
	"github.com/phuslu/log"
)

func GenerateValidationErrorMessage(err error) string {
	var castedObject validator.ValidationErrors
	errors.As(err, &castedObject)
	messageSlice := make([]string, 0)

	/*
		ensure that error message only shows once,
		if we do not create the map then the error message will be ex:
		"failed create role, type is required, type is required"
		it will shows up for x times the amount of same error encountered
	*/

	mp := make(map[string]bool)

	// add more validation rules below, ex: case "min": ...
	for _, err := range castedObject {
		if _, ok := mp[err.Field()]; ok {
			continue
		}
		switch err.Tag() {
		case "required":
			messageSlice = append(messageSlice, fmt.Sprintf("%s tidak boleh kosong", err.Field()))
		case "numeric":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus berupa angka", err.Field()))
		case "lt":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus kurang dari %s", err.Field(), err.Param()))
		case "gt":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus lebih dari %s", err.Field(), err.Param()))
		case "gte":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus lebih dari atau sama dengan %s", err.Field(), err.Param()))
		case "max":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus kurang dari %s", err.Field(), err.Param()))
		case "min":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus lebih dari %s", err.Field(), err.Param()))
		case "eqfield":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus sama dengan %s", err.Field(), err.Param()))
		case "eq":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus '%v'", err.Field(), err.Param()))
		case "eq=active|eq=nonactive":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus antara 'active' atau 'nonactive'", err.Field()))
		case "eq=meter|eq=kilometer":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus antara 'meter' atau 'kilometer'", err.Field()))
		case "base64":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus berformat base64", err.Field()))
		case "base64url":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus berupa format base64url string", err.Field()))
		case "longitude", "latitude":
			messageSlice = append(messageSlice, fmt.Sprintf("%s tidak valid", err.Tag()))
		}
		mp[err.Field()] = true
	}

	errMessage := strings.Join(messageSlice, ", ")

	return errMessage
}

func IsValidLatLong(lat, long float64) bool {
	if lat < -90 || lat > 90 {
		return false
	}
	if long < -180 || long > 180 {
		return false
	}
	return true
}

// haversine calculates the distance between two points on the Earth
// given their latitude and longitude in degrees.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // Earth radius in kilometers
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

// ip : from request
// lat, long and radius are set on PLTA data
// radius must be in kilometer
func IsIPWithinRadius(ip string, lat, long, radius float64) (bool, error) {
	// _, b, _, _ := runtime.Caller(0)
	// basepath := filepath.Dir(b)
	db, err := geoIp.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Error().Err(errors.New("ERROR GETTING GeoLite2-City.mmdb : " + err.Error())).Msg("")
		return false, err
	}
	defer db.Close()

	parsedIP := net.ParseIP(ip)
	record, err := db.City(parsedIP)
	if err != nil {
		log.Error().Err(errors.New("FAILED TO GET LOCATION DATA : " + err.Error())).Msg("")
		return false, err
	}

	distance := haversine(
		record.Location.Latitude,
		record.Location.Longitude,
		lat,
		long)

	fmt.Printf("LAT LONG FROM IP REQUEST : %f | %f\n", record.Location.Latitude, record.Location.Longitude)
	fmt.Printf("LAT LONG FROM ROOT : %f | %f\n", lat, long)
	fmt.Printf("MAX RADIUS : %f\n", radius)
	fmt.Printf("DISTANCE IN KM : %f\n", distance)

	if distance <= radius {
		return true, nil
	} else {
		return false, nil
	}
}
