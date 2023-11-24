package apiInit

import webServerInfra "github.com/speedianet/os/src/infra/webServer"

func WebServerSetup() {
	webServerInfra.WebServerFirstSetup()
	webServerInfra.WebServerOnStartSetup()
}
