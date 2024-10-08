package presentation

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
	wsInfra "github.com/goinfinite/os/src/infra/webServer"
	"github.com/goinfinite/os/src/presentation/api"
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

func HttpServerInit(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) {
	e := echo.New()

	api.ApiInit(e, persistentDbSvc, transientDbSvc, trailDbSvc)
	ui.UiInit(e, persistentDbSvc, transientDbSvc)

	httpServer := http.Server{Addr: ":1618", Handler: e}

	webServerSetup(persistentDbSvc, transientDbSvc)

	pkiDir := "/infinite/pki"
	certFile := pkiDir + "/os.crt"
	keyFile := pkiDir + "/os.key"
	if !infraHelper.FileExists(certFile) {
		err := infraHelper.MakeDir(pkiDir)
		if err != nil {
			slog.Error("MakePkiDirFailed", slog.Any("error", err))
			os.Exit(1)
		}

		aliases := []string{"localhost", "127.0.0.1"}
		err = infraHelper.CreateSelfSignedSsl(pkiDir, "os", aliases)
		if err != nil {
			slog.Error("GenerateSelfSignedCertFailed", slog.Any("error", err))
			os.Exit(1)
		}
	}

	osBanner := `Infinite OS server started on [::]:1618! 🎉`

	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(transientDbSvc)
	o11yOverview, err := o11yQueryRepo.ReadOverview()
	if err == nil {
		devModeStr := ""
		if isDevMode, _ := voHelper.InterfaceToBool(os.Getenv("DEV_MODE")); isDevMode {
			devModeStr = "(🚧 DevMode 🚧)"
		}

		osBanner = `
         INFINITE
    ▄▄█▀▀██▄  ▄█▀▀▀█▄█   |  🔒 HTTPS server started on [::]:1618! ` + devModeStr + `        
  ▄██▀    ▀██▄██    ▀█   |
  ██▀      ▀█████▄       |  🏠 Primary Hostname: ` + o11yOverview.Hostname.String() + `
  ██        ██ ▀█████▄   |  ⏰ Uptime: ` + o11yOverview.UptimeRelative.String() + `
  ▀██▄    ▄██▀█     ██   |  🌐 IP Address: ` + o11yOverview.PublicIpAddress.String() + `
    ▀▀████▀▀ █▀█████▀    |  ⚙️  ` + o11yOverview.HardwareSpecs.String() + `
`
	}

	fmt.Println(osBanner)

	err = httpServer.ListenAndServeTLS(certFile, keyFile)
	if err != http.ErrServerClosed {
		slog.Error("HttpServerError", slog.Any("error", err))
		os.Exit(1)
	}
}
