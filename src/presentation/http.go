package presentation

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	wsInfra "github.com/speedianet/os/src/infra/webServer"
	"github.com/speedianet/os/src/presentation/api"
	"github.com/speedianet/os/src/presentation/ui"
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
) {
	e := echo.New()

	api.ApiInit(e, persistentDbSvc, transientDbSvc)
	ui.UiInit(e, persistentDbSvc)

	httpServer := http.Server{Addr: ":1618", Handler: e}

	webServerSetup(persistentDbSvc, transientDbSvc)

	pkiDir := "/speedia/pki"
	certFile := pkiDir + "/os.crt"
	keyFile := pkiDir + "/os.key"
	if !infraHelper.FileExists(certFile) {
		err := infraHelper.MakeDir(pkiDir)
		if err != nil {
			slog.Error("MakePkiDirFailed", slog.Any("error", err))
			os.Exit(1)
		}

		aliases := []string{}
		err = infraHelper.CreateSelfSignedSsl(pkiDir, "os", aliases)
		if err != nil {
			slog.Error("GenerateSelfSignedCertFailed", slog.Any("error", err))
			os.Exit(1)
		}
	}

	osBanner := `	
     ▒       ▒▓██████████████████████▒     ▓██████████████████████▓
   ▒█▓    ▒██████████      ▒██████████  ██████████▓             ▓██▒
  ▒█▓     ▓█████████▓      ██████████▓  ██████████▒
 ▓▓█▒▒   ▒██████████      ▓█████████▓    ▓▓███████████████████████
  ▒█▓    ▓█████████▓      ██████████▒   ▒▒             ▒██████████
   ▒    ▓██████████       ██████████▓  ████▓          ▒██████████
  ▒     ▒█████████████████████████▒   ██████████████████████████
_____________________________________________________________________

⇨ HTTPS server started on [::]:1618 and is ready to serve! 🎉
`

	fmt.Println(osBanner)

	err := httpServer.ListenAndServeTLS(certFile, keyFile)
	if err != http.ErrServerClosed {
		slog.Error("HttpServerError", slog.Any("error", err))
		os.Exit(1)
	}
}
