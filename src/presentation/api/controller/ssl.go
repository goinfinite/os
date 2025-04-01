package apiController

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	sslInfra "github.com/goinfinite/os/src/infra/ssl"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type SslController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	sslService      *service.SslService
}

func NewSslController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SslController {
	return &SslController{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		sslService: service.NewSslService(
			persistentDbSvc, transientDbSvc, trailDbSvc,
		),
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
) (rawVhostsStrSlice []string, err error) {
	var assertOk bool

	rawVhostsStrSlice, assertOk = rawVhostsInput.([]string)
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	rawVhosts, err := controller.parseRawVhosts(requestInputData["virtualHosts"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	requestInputData["virtualHosts"] = rawVhosts

	if _, exists := requestInputData["encodedCertificate"]; !exists {
		if _, exists = requestInputData["certificate"]; exists {
			requestInputData["encodedCertificate"] = requestInputData["certificate"]
		}
	}
	encodedCert, err := valueObject.NewEncodedContent(requestInputData["encodedCertificate"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	decodedCert, err := encodedCert.GetDecodedContent()
	if err != nil {
		return apiHelper.ResponseWrapper(
			c, http.StatusBadRequest, "CannotDecodeSslCertificateContent",
		)
	}
	requestInputData["certificate"] = decodedCert

	if _, exists := requestInputData["encodedKey"]; !exists {
		if _, exists = requestInputData["key"]; exists {
			requestInputData["encodedKey"] = requestInputData["key"]
		}
	}
	encodedKey, err := valueObject.NewEncodedContent(requestInputData["encodedKey"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}
	decodedKey, err := encodedKey.GetDecodedContent()
	if err != nil {
		return apiHelper.ResponseWrapper(
			c, http.StatusBadRequest, "CannotDecodeSslPrivateKeyContent",
		)
	}
	requestInputData["key"] = decodedKey

	return apiHelper.ServiceResponseWrapper(
		c, controller.sslService.Create(requestInputData),
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.sslService.Delete(requestInputData),
	)
}

func (controller *SslController) SslCertificateWatchdog() {
	validationIntervalMinutes := 60 / useCase.SslValidationsPerHour

	taskInterval := time.Duration(validationIntervalMinutes) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
	sslQueryRepo := sslInfra.SslQueryRepo{}
	sslCmdRepo := sslInfra.NewSslCmdRepo(
		controller.persistentDbSvc, controller.transientDbSvc,
	)
	sslWatchdogUseCase := useCase.NewSslCertificateWatchdog(
		vhostQueryRepo, sslQueryRepo, sslCmdRepo,
	)

	for range timer.C {
		sslWatchdogUseCase.Execute()
	}
}
