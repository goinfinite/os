package apiHelper

import (
	"net/http"

	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/labstack/echo/v4"
)

type newFormattedResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

func LiaisonResponseWrapper(
	c echo.Context,
	liaisonOutput liaison.LiaisonOutput,
) error {
	responseStatus := http.StatusOK
	switch liaisonOutput.Status {
	case liaison.Created:
		responseStatus = http.StatusCreated
	case liaison.MultiStatus:
		responseStatus = http.StatusMultiStatus
	case liaison.UserError:
		responseStatus = http.StatusBadRequest
	case liaison.InfraError:
		responseStatus = http.StatusInternalServerError
	}

	formattedResponse := newFormattedResponse{
		Status: responseStatus,
		Body:   liaisonOutput.Body,
	}
	return c.JSON(responseStatus, formattedResponse)
}
