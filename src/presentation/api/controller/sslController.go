package apiController

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
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

func parseVirtualHosts(vhostsBodyInput interface{}) []valueObject.Fqdn {
	var vhosts []valueObject.Fqdn

	switch rawVhosts := vhostsBodyInput.(type) {
	case string:
		vhosts = append(
			vhosts,
			valueObject.NewFqdnPanic(rawVhosts),
		)
	case []interface{}:
		for _, rawVhostInterface := range rawVhosts {
			vhostStr, assertOk := rawVhostInterface.(string)
			if !assertOk {
				continue
			}

			vhost, err := valueObject.NewFqdn(vhostStr)
			if err != nil {
				continue
			}

			vhosts = append(vhosts, vhost)
		}
	}

	return vhosts
}

// CreateSsl    	 godoc
// @Summary      CreateNewSslPair
// @Description  Create a new ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createSslPairDto 	  body    dto.CreateSslPair  true  "NewSslPair"
// @Success      201 {object} object{} "SslPairCreated"
// @Router       /ssl/ [post]
func CreateSslPairController(c echo.Context) error {
	requiredParams := []string{"virtualHosts", "certificate", "key"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sslCertificateEncoded := valueObject.NewEncodedContentPanic(
		requestBody["certificate"].(string),
	)
	sslCertificateContent := valueObject.NewSslCertificateContentFromEncodedContentPanic(
		sslCertificateEncoded,
	)
	sslCertificate := entity.NewSslCertificatePanic(sslCertificateContent)

	sslPrivateKeyEncoded := valueObject.NewEncodedContentPanic(requestBody["key"].(string))
	sslPrivateKey := valueObject.NewSslPrivateKeyFromEncodedContentPanic(sslPrivateKeyEncoded)

	virtualHosts := parseVirtualHosts(requestBody["virtualHosts"])

	createSslPairDto := dto.NewCreateSslPair(
		virtualHosts,
		sslCertificate,
		sslPrivateKey,
	)

	sslCmdRepo := sslInfra.NewSslCmdRepo()
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}

	err := useCase.CreateSslPair(
		sslCmdRepo,
		vhostQueryRepo,
		createSslPairDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "SslPairCreated")
}

// RemoveSslPairVhosts    	 godoc
// @Summary      RemoveSslPairVhosts
// @Description  Create a new ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        removeSslPairVhostsDto 	  body    dto.RemoveSslPairVhosts  true  "SslPairVhostsRemoved"
// @Success      200 {object} object{} "SslPairVhostsRemoved"
// @Router       /ssl/ [put]
func RemoveSslPairVhostsController(c echo.Context) error {
	requiredParams := []string{"sslPairId", "virtualHosts"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sslPairId := valueObject.NewSslIdPanic(requestBody["sslPairId"].(string))
	virtualHosts := parseVirtualHosts(requestBody["virtualHosts"])

	dto := dto.NewRemoveSslPairVhosts(sslPairId, virtualHosts)

	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslCmdRepo := sslInfra.NewSslCmdRepo()
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}

	err := useCase.RemoveSslPairVhosts(
		sslQueryRepo,
		sslCmdRepo,
		vhostQueryRepo,
		dto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "SslPairVhostsRemoved")
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
	sslCmdRepo := sslInfra.NewSslCmdRepo()

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

func SslCertificateWatchdogController() {
	validationIntervalMinutes := 60 / useCase.SslValidationsPerHour

	taskInterval := time.Duration(validationIntervalMinutes) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslCmdRepo := sslInfra.NewSslCmdRepo()
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	for range timer.C {
		sslCertificateWatchdog := useCase.NewSslCertificateWatchdog(
			sslQueryRepo,
			sslCmdRepo,
			vhostQueryRepo,
			vhostCmdRepo,
		)
		sslCertificateWatchdog.Execute()
	}
}
