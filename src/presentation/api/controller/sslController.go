package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
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

// AddSsl    godoc
// @Summary      AddNewSsl
// @Description  Add a new ssl.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addSslDto 	  body    dto.AddSsl  true  "NewSsl"
// @Success      201 {object} object{} "SslCreated"
// @Router       /ssl/ [post]
func AddSslController(c echo.Context) error {
	requiredParams := []string{"hostname", "certificate", "key"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sslPair, err := entity.NewSslPair(requestBody["certificate"].(string))
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	sslPrivateKey, err := entity.NewSslPrivateKey(requestBody["key"].(string))
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	addCronDto := dto.NewAddSsl(
		valueObject.NewFqdnPanic(requestBody["hostname"].(string)),
		sslPair,
		sslPrivateKey,
	)

	sslCmdRepo := infra.SslCmdRepo{}

	err = useCase.AddSsl(
		sslCmdRepo,
		addCronDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "SslCreated")
}

// DeleteSsl	 godoc
// @Summary      DeleteSsl
// @Description  Delete a ssl.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        sslId 	  path   string  true  "SslId"
// @Success      200 {object} object{} "SslDeleted"
// @Router       /ssl/{sslId}/ [delete]
func DeleteSslController(c echo.Context) error {
	sslId := valueObject.NewSslIdPanic(c.Param("sslId"))

	sslQueryRepo := infra.NewSslQueryRepo()
	sslCmdRepo := infra.SslCmdRepo{}

	err := useCase.DeleteSsl(
		sslQueryRepo,
		sslCmdRepo,
		sslId,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "SslDeleted")
}
