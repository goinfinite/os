package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
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
	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslPairsList, err := useCase.GetSslPairs(sslQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, sslPairsList)
}

func parseVirtualHosts(vhosts []interface{}) []valueObject.Fqdn {
	var virtualHosts []valueObject.Fqdn

	for _, vhost := range vhosts {
		vhostStr, assertOk := vhost.(string)
		if !assertOk {
			panic("InvalidVirtualHosts")
		}

		virtualHosts = append(virtualHosts, valueObject.NewFqdnPanic(vhostStr))
	}

	return virtualHosts
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
	requiredParams := []string{"virtualHosts", "certificate", "key"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sslCertificateContent := valueObject.NewSslCertificateContentPanic(requestBody["certificate"].(string))
	sslCertificate := entity.NewSslCertificatePanic(sslCertificateContent)
	sslPrivateKey := valueObject.NewSslPrivateKeyPanic(requestBody["key"].(string))

	virtualHosts, assertOk := requestBody["virtualHosts"].([]interface{})
	if !assertOk {
		panic("InvalidVirtualHosts")
	}

	addSslPairDto := dto.NewAddSslPair(
		parseVirtualHosts(virtualHosts),
		sslCertificate,
		sslPrivateKey,
	)

	sslCmdRepo := sslInfra.SslCmdRepo{}

	err := useCase.AddSslPair(
		sslCmdRepo,
		addSslPairDto,
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

	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslCmdRepo := sslInfra.SslCmdRepo{}

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
