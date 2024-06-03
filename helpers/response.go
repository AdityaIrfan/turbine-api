package helpers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type response struct {
	Success    bool        `json:"Success"`
	ErrorCode  int32       `json:"ErrorCode,omitempty"`
	StatusCode int32       `json:"StatusCode"`
	Message    string      `json:"Message"`
	Data       interface{} `json:"Data,omitempty"`
	Meta       interface{} `json:"Meta,omitempty"`
}

func Response(c echo.Context, statusCode int32, message string, ins ...interface{}) error {
	var data, meta interface{}
	var status bool

	for index, in := range ins {
		if index == 0 {
			data = in
		} else if index == 1 {
			meta = in
		} else {
			break
		}
	}

	if statusCode == http.StatusOK {
		status = true
	}

	res := &response{
		Success:    status,
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
		Meta:       meta,
	}

	return c.JSON(int(statusCode), res)
}

func ResponseUnprocessableEntity(c echo.Context) error {
	return Response(c, http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))
}

func ResponseNonAdminForbiddenAccess(c echo.Context) error {
	return Response(c, http.StatusUnauthorized, "who are you? you do not have permission for this access")
}
