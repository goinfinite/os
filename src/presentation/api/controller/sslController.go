package apiController

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

type SslController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
}

func NewSslController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *SslController {
	return &SslController{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
	}
}

// ReadSslPairs	 godoc
// @Summary      ReadSslPairs
// @Description  List ssl pairs.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.SslPair
// @Router       /v1/ssl/ [get]
func (controller *SslController) Read(c echo.Context) error {
	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslPairsList, err := useCase.ReadSslPairs(sslQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, sslPairsList)
}

func parseVirtualHosts(vhostsBodyInput interface{}) []valueObject.Fqdn {
	_, isStringType := vhostsBodyInput.(string)
	if isStringType {
		vhostsBodyInput = []interface{}{vhostsBodyInput}
	}

	rawVhosts, isInterfaceSliceType := vhostsBodyInput.([]interface{})
	if !isInterfaceSliceType {
		panic("InvalidVirtualHosts")
	}

	vhosts := []valueObject.Fqdn{}
	for _, rawVhost := range rawVhosts {
		rawVhostStr, assertOk := rawVhost.(string)
		if !assertOk {
			continue
		}

		vhosts = append(vhosts, valueObject.NewFqdnPanic(rawVhostStr))
	}

	return vhosts
}

// CreateSslPair    	 godoc
// @Summary      CreateSslPair
// @Description  Create a new ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createSslPairDto 	  body    dto.CreateSslPair  true  "All props are required.<br />virtualHosts may be string or []string. Alias is not allowed.<br />certificate is a string field, i.e. ignore the structure shown.<br />certificate and key must be base64 encoded."
// @Success      201 {object} object{} "SslPairCreated"
// @Router       /v1/ssl/ [post]
func (controller *SslController) Create(c echo.Context) error {
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

	sslCmdRepo := sslInfra.NewSslCmdRepo(
		controller.persistentDbSvc, controller.transientDbSvc,
	)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

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

// DeleteSslPair	 godoc
// @Summary      DeleteSslPair
// @Description  Delete a ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        sslPairId 	  path   string  true  "SslPairId to delete."
// @Success      200 {object} object{} "SslPairDeleted"
// @Router       /v1/ssl/{sslPairId}/ [delete]
func (controller *SslController) Delete(c echo.Context) error {
	sslSerialNumber := valueObject.NewSslIdPanic(c.Param("sslPairId"))

	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslCmdRepo := sslInfra.NewSslCmdRepo(
		controller.persistentDbSvc, controller.transientDbSvc,
	)

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

func (controller *SslController) SslCertificateWatchdog() {
	validationIntervalMinutes := 60 / useCase.SslValidationsPerHour

	taskInterval := time.Duration(validationIntervalMinutes) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslCmdRepo := sslInfra.NewSslCmdRepo(
		controller.persistentDbSvc, controller.transientDbSvc,
	)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
	vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(controller.persistentDbSvc)

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

// DeleteSslPairVhosts    	 godoc
// @Summary      DeleteSslPairVhosts
// @Description  Delete vhosts from a ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        deleteSslPairVhostsDto 	  body    dto.DeleteSslPairVhosts  true  "All props are required."
// @Success      200 {object} object{} "SslPairVhostsRemoved"
// @Router       /v1/ssl/vhost/ [put]
func (controller *SslController) DeleteVhosts(c echo.Context) error {
	requiredParams := []string{"sslPairId", "virtualHosts"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sslPairId := valueObject.NewSslIdPanic(requestBody["sslPairId"].(string))
	virtualHosts := parseVirtualHosts(requestBody["virtualHosts"])

	dto := dto.NewDeleteSslPairVhosts(sslPairId, virtualHosts)

	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslCmdRepo := sslInfra.NewSslCmdRepo(
		controller.persistentDbSvc, controller.transientDbSvc,
	)

	err := useCase.DeleteSslPairVhosts(
		sslQueryRepo,
		sslCmdRepo,
		dto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "SslPairVhostsDeleted")
}
