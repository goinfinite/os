package presentation

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
	wsInfra "github.com/goinfinite/os/src/infra/webServer"
	"github.com/goinfinite/os/src/presentation/api"
	presentationMiddleware "github.com/goinfinite/os/src/presentation/middleware"
	"github.com/goinfinite/os/src/presentation/ui"
	"github.com/labstack/echo/v4"
)

func webServerSetup(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) {
	ws := wsInfra.NewWebServerSetup(persistentDbSvc, transientDbSvc)
	ws.FirstSetup()
	ws.OnStartSetup()
}

func initialSslSetup() (
	certFilePath tkValueObject.UnixAbsoluteFilePath,
	keyFilePath tkValueObject.UnixAbsoluteFilePath,
	err error,
) {
	rawOsPkiDir := "/infinite/pki"
	pkiDir, err := valueObject.NewUnixFilePath(rawOsPkiDir)
	if err != nil {
		return certFilePath, keyFilePath, errors.New("InvalidPkiDir")
	}
	osPkiDirStr := pkiDir.String()

	rawCertFilePath := osPkiDirStr + "/os.crt"
	certFilePath, err = valueObject.NewUnixFilePath(rawCertFilePath)
	if err != nil {
		return certFilePath, keyFilePath, errors.New("InvalidCertFilePath")
	}

	rawKeyFilePath := osPkiDirStr + "/os.key"
	keyFilePath, err = valueObject.NewUnixFilePath(rawKeyFilePath)
	if err != nil {
		return certFilePath, keyFilePath, errors.New("InvalidKeyFilePath")
	}

	if !infraHelper.FileExists(certFilePath.String()) {
		err := infraHelper.MakeDir(osPkiDirStr)
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
		if isDevMode, _ := tkVoUtil.InterfaceToBool(os.Getenv("DEV_MODE")); isDevMode {
			devModeStr = "(🚧 DevMode 🚧)"
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
	echoInstance.Use(presentationMiddleware.PanicHandler)

	rootBasePath := "/"
	apiBasePath := rootBasePath + "api"
	uiBasePath := rootBasePath + ""
	if uiBasePath == rootBasePath {
		uiBasePath = ""
	}

	echoInstance.Use(presentationMiddleware.BaseHref(rootBasePath, apiBasePath, uiBasePath))

	api.ApiInit(echoInstance, apiBasePath, persistentDbSvc, transientDbSvc, trailDbSvc)
	ui.UiInit(echoInstance, uiBasePath, persistentDbSvc, transientDbSvc, trailDbSvc)

	webServerSetup(persistentDbSvc, transientDbSvc)

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
