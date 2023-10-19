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

// GetSslPairs	 godoc
// @Summary      GetSslPair
// @Description  List ssl pairs.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.SslPair
// @Router       /ssl/ [get]
func GetSslPairsController(c echo.Context) error {
	sslQueryRepo := infra.SslQueryRepo{}
	sslPairsList, err := useCase.GetSslPairs(sslQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, sslPairsList)
}

// AddSsl    	 godoc
// @Summary      AddNewSslPair
// @Description  Add a new ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addSslPairDto 	  body    dto.AddSslPair  true  "NewSslPair"
// @Success      201 {object} object{} "SslPairCreated"
// @Router       /ssl/ [post]
func AddSslPairController(c echo.Context) error {
	requiredParams := []string{"hostname", "certificate", "key"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sslCertificate := entity.NewSslCertificatePanic(requestBody["certificate"].(string))
	sslPrivateKey := valueObject.NewSslPrivateKeyPanic(requestBody["key"].(string))

	addCronDto := dto.NewAddSslPair(
		valueObject.NewFqdnPanic(requestBody["hostname"].(string)),
		sslCertificate,
		sslPrivateKey,
	)

	sslCmdRepo := infra.SslCmdRepo{}

	err := useCase.AddSslPair(
		sslCmdRepo,
		addCronDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "SslPairCreated")
}

// DeleteSsl	 godoc
// @Summary      DeleteSslPair
// @Description  Delete a ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        sslPairId 	  path   string  true  "SslPairId"
// @Success      200 {object} object{} "SslPairDeleted"
// @Router       /ssl/{sslPairId}/ [delete]
func DeleteSslPairController(c echo.Context) error {
	sslSerialNumber := valueObject.NewSslIdPanic(c.Param("sslPairId"))

	sslQueryRepo := infra.SslQueryRepo{}
	sslCmdRepo := infra.SslCmdRepo{}

	err := useCase.DeleteSslPair(
		sslQueryRepo,
		sslCmdRepo,
		sslSerialNumber,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "SslPairDeleted")
}
