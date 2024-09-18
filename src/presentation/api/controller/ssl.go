package apiController

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	"github.com/speedianet/os/src/presentation/service"
)

type SslController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	sslService      *service.SslService
}

func NewSslController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *SslController {
	return &SslController{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		sslService:      service.NewSslService(persistentDbSvc, transientDbSvc),
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
	return apiHelper.ServiceResponseWrapper(c, controller.sslService.Read())
}

func (controller *SslController) parseRawVhosts(
	rawVhostsInput interface{},
) ([]string, error) {
	rawVhostsStrSlice, assertOk := rawVhostsInput.([]string)
	if assertOk {
		return rawVhostsStrSlice, nil
	}

	rawVhostsInterfaceSlice, assertOk := rawVhostsInput.([]interface{})
	if !assertOk {
		rawVhostUniqueStr, err := voHelper.InterfaceToString(rawVhostsInput)
		if err != nil {
			return rawVhostsStrSlice, errors.New("VirtualHostsMustBeStringOrStringSlice")
		}
		return append(rawVhostsStrSlice, rawVhostUniqueStr), err
	}

	for _, rawVhost := range rawVhostsInterfaceSlice {
		rawVhostStr, err := voHelper.InterfaceToString(rawVhost)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("vhost", rawVhost))
			continue
		}
		rawVhostsStrSlice = append(rawVhostsStrSlice, rawVhostStr)
	}

	return rawVhostsStrSlice, nil
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
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	rawVhosts, err := controller.parseRawVhosts(requestBody["virtualHosts"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	requestBody["virtualHosts"] = rawVhosts

	if _, exists := requestBody["encodedCertificate"]; !exists {
		if _, exists = requestBody["certificate"]; exists {
			requestBody["encodedCertificate"] = requestBody["certificate"]
		}
	}
	encodedCert, err := valueObject.NewEncodedContent(requestBody["encodedCertificate"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	decodedCert, err := encodedCert.GetDecodedContent()
	if err != nil {
		return apiHelper.ResponseWrapper(
			c, http.StatusBadRequest, "CannotDecodeSslCertificateContent",
		)
	}
	requestBody["certificate"] = decodedCert

	if _, exists := requestBody["encodedKey"]; !exists {
		if _, exists = requestBody["key"]; exists {
			requestBody["encodedKey"] = requestBody["key"]
		}
	}
	encodedKey, err := valueObject.NewEncodedContent(requestBody["encodedKey"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	decodedKey, err := encodedKey.GetDecodedContent()
	if err != nil {
		return apiHelper.ResponseWrapper(
			c, http.StatusBadRequest, "CannotDecodeSslPrivateKeyContent",
		)
	}
	requestBody["key"] = decodedKey

	return apiHelper.ServiceResponseWrapper(
		c, controller.sslService.Create(requestBody),
	)
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
	requestBody := map[string]interface{}{
		"id": c.Param("sslPairId"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.sslService.Delete(requestBody),
	)
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
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	rawVhosts, err := controller.parseRawVhosts(requestBody["virtualHosts"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	requestBody["virtualHosts"] = rawVhosts

	return apiHelper.ServiceResponseWrapper(
		c, controller.sslService.DeleteVhosts(requestBody),
	)
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

	for range timer.C {
		useCase.SslCertificateWatchdog(sslQueryRepo, sslCmdRepo)
	}
}
