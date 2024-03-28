package wsInfra

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDatabaseInfra "github.com/speedianet/os/src/infra/internalDatabase"
	o11yInfra "github.com/speedianet/os/src/infra/o11y"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	envDataInfra "github.com/speedianet/os/src/infra/shared"
)

type WebServerSetup struct {
	transientDbSvc *internalDatabaseInfra.TransientDatabaseService
}

func NewWebServerSetup(
	transientDbSvc *internalDatabaseInfra.TransientDatabaseService,
) *WebServerSetup {
	return &WebServerSetup{
		transientDbSvc: transientDbSvc,
	}
}

func (ws *WebServerSetup) updatePhpMaxChildProcesses(memoryTotal valueObject.Byte) error {
	log.Print("UpdatingMaxPhpChildProcesses...")

	maxChildProcesses := int64(300)
	childProcessPerGb := int64(5)

	memoryInGb := memoryTotal.ToGiB()
	desiredChildProcesses := memoryInGb * childProcessPerGb
	if desiredChildProcesses > maxChildProcesses {
		desiredChildProcesses = maxChildProcesses
	}

	desiredChildProcessesStr := strconv.FormatInt(desiredChildProcesses, 10)
	httpdConfFilePath := "/usr/local/lsws/conf/httpd_config.conf"
	_, err := infraHelper.RunCmd(
		"sed",
		"-i",
		"-e",
		"s/PHP_LSAPI_CHILDREN=[0-9]+/PHP_LSAPI_CHILDREN="+desiredChildProcessesStr+";/g",
		httpdConfFilePath,
	)
	if err != nil {
		return errors.New("UpdateMaxChildProcessesFailed")
	}

	return nil
}

func (ws *WebServerSetup) FirstSetup() {
	_, err := os.Stat("/etc/nginx/dhparam.pem")
	if err == nil {
		return
	}

	log.Print("FirstBootDetected! Please await while the web server is configured...")

	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		log.Fatal("PrimaryVirtualHostNotFound")
	}

	primaryVhostStr := primaryVhost.String()

	log.Print("UpdatingVhost...")

	primaryConfFilePath := "/app/conf/nginx/primary.conf"
	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/speedia.net/"+primaryVhostStr+"/g",
		primaryConfFilePath,
	)
	if err != nil {
		log.Fatal("UpdateVhostFailed")
	}

	log.Print("GeneratingDhparam...")

	_, err = infraHelper.RunCmd(
		"openssl",
		"dhparam",
		"-dsaparam",
		"-out",
		"/etc/nginx/dhparam.pem",
		"2048",
	)
	if err != nil {
		log.Fatal("GenerateDhparamFailed")
	}

	log.Print("GeneratingSelfSignedCert...")

	err = infraHelper.CreateSelfSignedSsl(envDataInfra.PkiConfDir, primaryVhostStr)
	if err != nil {
		log.Fatal("GenerateSelfSignedCertFailed")
	}

	log.Print("WebServerConfigured!")

	err = servicesInfra.SupervisordFacade{}.Start("nginx")
	if err != nil {
		log.Fatal("StartNginxFailed")
	}
}

func (ws *WebServerSetup) OnStartSetup() {
	defaultLogPrefix := "WsOnStartupSetup"

	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(ws.transientDbSvc)
	containerResources, err := o11yQueryRepo.GetOverview()
	if err != nil {
		log.Fatalf("%sGetContainerResourcesFailed", defaultLogPrefix)
	}

	cpuCores := containerResources.HardwareSpecs.CpuCores
	cpuCoresStr := strconv.FormatUint(cpuCores, 10)

	nginxConfFilePath := "/etc/nginx/nginx.conf"
	workerCount, err := infraHelper.RunCmd(
		"awk",
		"/worker_processes/{gsub(/[^0-9]+/, \"\"); print}",
		nginxConfFilePath,
	)
	if err != nil {
		log.Fatalf("%sGetNginxWorkersCountFailed", defaultLogPrefix)
	}

	if workerCount == cpuCoresStr {
		return
	}

	log.Print("UpdatingNginxWorkersCount...")

	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"-e",
		"s/^worker_processes.*/worker_processes "+cpuCoresStr+";/g",
		nginxConfFilePath,
	)
	if err != nil {
		log.Fatalf("%sUpdateNginxWorkersCountFailed", defaultLogPrefix)
	}

	err = servicesInfra.SupervisordFacade{}.Restart("nginx")
	if err != nil {
		log.Fatalf("%sRestartNginxFailed", defaultLogPrefix)
	}

	_, err = servicesInfra.ServicesQueryRepo{}.GetByName("php")
	if err == nil {
		err = ws.updatePhpMaxChildProcesses(
			containerResources.HardwareSpecs.MemoryTotal,
		)
		if err != nil {
			log.Fatalf("%s%s", defaultLogPrefix, err.Error())
		}
	}
}
