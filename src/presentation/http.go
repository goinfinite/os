package presentation

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
	"github.com/goinfinite/os/src/presentation/api"
	presentationMiddleware "github.com/goinfinite/os/src/presentation/middleware"
	"github.com/goinfinite/os/src/presentation/ui"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	tkInfra "github.com/goinfinite/tk/src/infra"
	tkPresentationMiddleware "github.com/goinfinite/tk/src/presentation/middleware"
	"github.com/labstack/echo/v4"
)

var fileClerk = tkInfra.FileClerk{}

func initialSslSetup() (
	certFilePath tkValueObject.UnixAbsoluteFilePath,
	keyFilePath tkValueObject.UnixAbsoluteFilePath,
	err error,
) {
	rawOsPkiDir := "/infinite/pki"
	pkiDir, err := tkValueObject.NewUnixAbsoluteFilePath(rawOsPkiDir, false)
	if err != nil {
		return certFilePath, keyFilePath, errors.New("InvalidPkiDir")
	}
	osPkiDirStr := pkiDir.String()

	rawCertFilePath := osPkiDirStr + "/os.crt"
	certFilePath, err = tkValueObject.NewUnixAbsoluteFilePath(rawCertFilePath, false)
	if err != nil {
		return certFilePath, keyFilePath, errors.New("InvalidCertFilePath")
	}

	rawKeyFilePath := osPkiDirStr + "/os.key"
	keyFilePath, err = tkValueObject.NewUnixAbsoluteFilePath(rawKeyFilePath, false)
	if err != nil {
		return certFilePath, keyFilePath, errors.New("InvalidKeyFilePath")
	}

	if !fileClerk.FileExists(certFilePath.String()) {
		err := fileClerk.CreateDir(osPkiDirStr)
		if err != nil {
			return certFilePath, keyFilePath, errors.New("CreatePkiDirFailed")
		}

		err = infraHelper.CreateSelfSignedSsl(pkiDir, "os", []tkValueObject.Fqdn{})
		if err != nil {
			return certFilePath, keyFilePath, errors.New("CreateSelfSignedSslFailed")
		}
	}

	return certFilePath, keyFilePath, nil
}

func initialBannerSetup(
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) {
	osBanner := `Infinite OS server started on [::]:` + infraEnvs.InfiniteOsApiHttpPublicPort + `! 🎉`

	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(transientDbSvc)
	o11yOverview, err := o11yQueryRepo.ReadOverview(false)
	if err == nil {
		devModeStr := ""
		isDevModeEnabled, err := tkVoUtil.InterfaceToBool(os.Getenv("DEV_MODE"))
		if err == nil && isDevModeEnabled {
			devModeStr = "(DevMode ON)"
		}

		osBanner = `
        INFINITE
    ▄▄█▀▀██▄  ▄█▀▀▀█▄█   |  🔒 HTTPS server started on [::]:` + infraEnvs.InfiniteOsApiHttpPublicPort + `! ` + devModeStr + `
  ▄██▀    ▀██▄██    ▀█   |
  ██▀      ▀█████▄       |  🏠 Primary Hostname: ` + o11yOverview.Hostname.String() + `
  ██        ██ ▀█████▄   |  ⏰ Uptime: ` + o11yOverview.UptimeRelative.String() + `
  ▀██▄    ▄██▀█     ██   |  🌐 IP Address: ` + o11yOverview.PublicIpAddress.String() + `
    ▀▀████▀▀ █▀█████▀    |  ⚙️  ` + o11yOverview.HardwareSpecs.String() + `
`
	}

	fmt.Println(osBanner)
}

func HttpServerInit(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) {
	echoInstance := echo.New()
	echoInstance.Use(tkPresentationMiddleware.ApiPanicHandler)

	rootBasePath := "/"
	apiBasePath := rootBasePath + "api"
	uiBasePath := rootBasePath + ""
	if uiBasePath == rootBasePath {
		uiBasePath = ""
	}

	echoInstance.Use(
		presentationMiddleware.BaseHref(rootBasePath, apiBasePath, uiBasePath),
	)

	api.ApiInit(echoInstance, apiBasePath, persistentDbSvc, transientDbSvc, trailDbSvc)
	ui.UiInit(echoInstance, uiBasePath, persistentDbSvc, transientDbSvc, trailDbSvc)

	certFilePath, keyFilePath, err := initialSslSetup()
	if err != nil {
		slog.Error("InitialSslSetupError", slog.String("err", err.Error()))
		os.Exit(1)
	}

	initialBannerSetup(transientDbSvc)

	httpServer := http.Server{
		Addr:    ":" + infraEnvs.InfiniteOsApiHttpPublicPort,
		Handler: echoInstance,
	}
	err = httpServer.ListenAndServeTLS(certFilePath.String(), keyFilePath.String())
	if err != http.ErrServerClosed {
		slog.Error("HttpServerError", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
