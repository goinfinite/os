package apiController

import (
	"errors"
	"net/http"
	"time"

	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	sslInfra "github.com/goinfinite/os/src/infra/ssl"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
	"github.com/labstack/echo/v4"
)

type SslController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	sslLiaison      *liaison.SslLiaison
}

func NewSslController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SslController {
	return &SslController{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		sslLiaison: liaison.NewSslLiaison(
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
// @Param        sslPairId query  string  false  "SslPairId"
// @Param        virtualHostHostname query  string  false  "VirtualHostHostname"
// @Param        altNames query  string  false  "AltNames"
// @Param        issuedBeforeAt query  string  false  "IssuedBeforeAt"
// @Param        issuedAfterAt query  string  false  "IssuedAfterAt"
// @Param        expiresBeforeAt query  string  false  "ExpiresBeforeAt"
// @Param        expiresAfterAt query  string  false  "ExpiresAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadSslPairsResponse
// @Router       /v1/ssl/ [get]
func (controller *SslController) Read(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if requestInputData["altNames"] != nil {
		requestInputData["altNames"] = sharedHelper.StringSliceValueObjectParser(
			requestInputData["altNames"], valueObject.NewSslHostname,
		)
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.sslLiaison.Read(requestInputData),
	)
}

func (controller *SslController) decodeContent(
	rawContent interface{},
) (decodedContent string, err error) {
	encodedContent, err := valueObject.NewEncodedContent(rawContent)
	if err != nil {
		return decodedContent, err
	}
	decodedContent, err = encodedContent.GetDecodedContent()
	if err != nil {
		return decodedContent, errors.New("CannotDecodeContent")
	}

	return decodedContent, nil
}

// CreateSslPair    	 godoc
// @Summary      CreateSslPair
// @Description  Create a new ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createSslPairDto 	  body    dto.CreateSslPair  true  "All props are required.<br />virtualHosts may be string or []string. Alias is not allowed.<br />certificate is a string field, i.e. ignore the structure shown.<br />certificate and key must be base64 encoded.<br />certificate should include the CA chain/bundle if not provided in the certificate field."
// @Success      201 {object} object{} "SslPairCreated"
// @Router       /v1/ssl/ [post]
func (controller *SslController) Create(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if requestInputData["virtualHostsHostnames"] == nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "VirtualHostHostnameIsRequired")
	}

	requestInputData["virtualHostsHostnames"] = sharedHelper.StringSliceValueObjectParser(
		requestInputData["virtualHostsHostnames"], valueObject.NewFqdn,
	)

	if requestInputData["encodedCertificate"] != nil {
		requestInputData["certificate"], err = controller.decodeContent(
			requestInputData["encodedCertificate"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "CannotDecodeSslCertificateContent")
		}
	}

	if requestInputData["encodedKey"] != nil {
		requestInputData["key"], err = controller.decodeContent(
			requestInputData["encodedKey"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "CannotDecodeSslKeyContent")
		}
	}

	if requestInputData["encodedChainCertificates"] != nil && requestInputData["encodedChainCertificates"] != "" {
		requestInputData["chainCertificates"], err = controller.decodeContent(
			requestInputData["encodedChainCertificates"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "CannotDecodeSslChainCertificatesContent")
		}
	}

	if requestInputData["chainCertificates"] != nil && requestInputData["chainCertificates"] == "" {
		requestInputData["chainCertificates"] = nil
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.sslLiaison.Create(requestInputData),
	)
}

// CreatePubliclyTrusted godoc
// @Summary      CreatePubliclyTrusted
// @Description  Create a new publicly trusted ssl pair.
// @Tags         ssl
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createPubliclyTrustedDto 	  body    dto.CreatePubliclyTrustedSslPair  true "All props are required."
// @Success      201 {object} object{} "PubliclyTrustedSslPairCreationScheduled"
// @Router       /v1/ssl/trusted/ [post]
func (controller *SslController) CreatePubliclyTrusted(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.sslLiaison.CreatePubliclyTrusted(requestInputData, true),
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

	return apiHelper.LiaisonResponseWrapper(
		c, controller.sslLiaison.Delete(requestInputData),
	)
}

func (controller *SslController) SslCertificateWatchdog() {
	validationIntervalMinutes := 60 / useCase.SslValidationsPerHour

	taskInterval := time.Duration(validationIntervalMinutes) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
	sslQueryRepo := sslInfra.NewSslQueryRepo()
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
