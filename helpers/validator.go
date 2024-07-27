package helpers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
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
		case "base64":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus berformat base64", err.Field()))
		case "base64url":
			messageSlice = append(messageSlice, fmt.Sprintf("%s harus berupa format base64url string", err.Field()))
		}
		mp[err.Field()] = true
	}

	errMessage := strings.Join(messageSlice, ", ")

	return errMessage
}
