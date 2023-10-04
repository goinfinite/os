package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/infra"
	apiHelper "github.com/speedianet/sam/src/presentation/api/helper"
)

// GetSsls	 	 godoc
// @Summary      GetSsls
// @Description  List ssls.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.Ssl
// @Router       /ssl/ [get]
func GetSslsController(c echo.Context) error {
	sslsQueryRepo := infra.NewSslQueryRepo()
	sslsList, err := useCase.GetSsls(sslsQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, sslsList)
}
