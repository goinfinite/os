package wsInfra

import (
	"errors"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/speedianet/os/src/domain/valueObject"
	infraEnvs "github.com/speedianet/os/src/infra/envs"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	o11yInfra "github.com/speedianet/os/src/infra/o11y"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

type WebServerSetup struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
}

func NewWebServerSetup(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *WebServerSetup {
	return &WebServerSetup{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
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
	_, err := infraHelper.RunCmd(
		"sed",
		"-i",
		"-e",
		"s/PHP_LSAPI_CHILDREN=[0-9]+/PHP_LSAPI_CHILDREN="+desiredChildProcessesStr+";/g",
		infraEnvs.PhpWebserverMainConfFilePath,
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

	aliases := []string{}
	err = infraHelper.CreateSelfSignedSsl(
		infraEnvs.PkiConfDir,
		primaryVhostStr,
		aliases,
	)
	if err != nil {
		log.Fatal("GenerateSelfSignedCertFailed")
	}

	log.Print("WebServerConfigured!")

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(ws.persistentDbSvc)
	serviceName, _ := valueObject.NewServiceName("nginx")
	err = servicesCmdRepo.Start(serviceName)
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
	cpuCoresStr := strconv.FormatInt(int64(math.Ceil(cpuCores)), 10)

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

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(ws.persistentDbSvc)
	serviceName, _ := valueObject.NewServiceName("nginx")
	err = servicesCmdRepo.Restart(serviceName)
	if err != nil {
		log.Fatalf("%sRestartNginxFailed", defaultLogPrefix)
	}

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(ws.persistentDbSvc)
	_, err = servicesQueryRepo.ReadByName("php-webserver")
	if err == nil {
		err = ws.updatePhpMaxChildProcesses(
			containerResources.HardwareSpecs.MemoryTotal,
		)
		if err != nil {
			log.Fatalf("%s%s", defaultLogPrefix, err.Error())
		}
	}
}
