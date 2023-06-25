package apiHelper

import (
	"github.com/labstack/echo/v4"
)

type formattedResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

func ResponseWrapper(
	c echo.Context,
	responseStatus int,
	responseBody interface{},
) error {
	formattedResponse := formattedResponse{
		Status: responseStatus,
		Body:   responseBody,
	}
	return c.JSON(responseStatus, formattedResponse)
}
