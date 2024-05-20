package presentation

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	wsInfra "github.com/speedianet/os/src/infra/webServer"
	"github.com/speedianet/os/src/presentation/api"
	"github.com/speedianet/os/src/presentation/ui"
)

type CustomLogger struct{}

func (*CustomLogger) Write(rawMessage []byte) (int, error) {
	messageStr := strings.TrimSpace(string(rawMessage))

	shouldLog := true
	if strings.HasSuffix(messageStr, "tls: unknown certificate") {
		shouldLog = false
	}

	messageLen := len(rawMessage)
	if !shouldLog {
		return messageLen, nil
	}

	return messageLen, log.Output(2, messageStr)
}

func NewCustomLogger() *log.Logger {
	return log.New(&CustomLogger{}, "", 0)
}

func webServerSetup(transientDbSvc *internalDbInfra.TransientDatabaseService) {
	ws := wsInfra.NewWebServerSetup(transientDbSvc)
	ws.FirstSetup()
	ws.OnStartSetup()
}

func HttpServerInit(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) {
	e := echo.New()

	api.ApiInit(e, persistentDbSvc, transientDbSvc)
	ui.UiInit(e)

	httpServer := http.Server{
		Addr:     ":1618",
		Handler:  e,
		ErrorLog: NewCustomLogger(),
	}

	webServerSetup(transientDbSvc)

	pkiDir := "/speedia/pki"
	certFile := pkiDir + "/os.crt"
	keyFile := pkiDir + "/os.key"
	if !infraHelper.FileExists(certFile) {
		err := infraHelper.MakeDir(pkiDir)
		if err != nil {
			log.Fatalf("MakePkiDirFailed: %v", err)
		}
		err = infraHelper.CreateSelfSignedSsl(pkiDir, "os")
		if err != nil {
			log.Fatalf("GenerateSelfSignedCertFailed: %v", err)
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
		log.Fatal(err)
	}
}
