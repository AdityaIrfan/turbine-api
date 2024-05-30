package helpers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func GenerateValidationErrorMessage(err error) string {
	castedObject, _ := err.(validator.ValidationErrors)
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
			messageSlice = append(messageSlice, fmt.Sprintf("%s is required", err.Field()))
		case "numeric":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be a number", err.Field()))
		case "lt":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be less than %s", err.Field(), err.Param()))
		case "gt":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be greater than %s", err.Field(), err.Param()))
		case "gte":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be greater than equal %s", err.Field(), err.Param()))
		case "max":
			messageSlice = append(messageSlice, fmt.Sprintf("%s length must be less than %s", err.Field(), err.Param()))
		case "min":
			messageSlice = append(messageSlice, fmt.Sprintf("%s length must greater than %s", err.Field(), err.Param()))
		case "eqfield":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be similar with %s", err.Field(), err.Param()))
		case "eq":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be '%v'", err.Field(), err.Param()))
		case "eq=active|eq=nonactive":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be 'active' or 'nonactive'", err.Field()))
		case "base64":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be base64 format", err.Field()))
		case "base64url":
			messageSlice = append(messageSlice, fmt.Sprintf("%s must be base64url string format", err.Field()))
		}
		mp[err.Field()] = true
	}

	errMessage := strings.Join(messageSlice, ", ")

	return errMessage
}
