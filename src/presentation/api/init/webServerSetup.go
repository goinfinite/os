package apiInit

import wsInfra "github.com/speedianet/os/src/infra/webServer"

func WebServerSetup() {
	wsInfra.WebServerFirstSetup()
	wsInfra.WebServerOnStartSetup()
}
