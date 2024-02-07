package wsInfra

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	o11yInfra "github.com/speedianet/os/src/infra/o11y"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
)

type WebServerSetup struct{}

func (ws WebServerSetup) updatePhpMaxChildProcesses(memoryTotal valueObject.Byte) error {
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

func (ws WebServerSetup) FirstSetup() {
	_, err := os.Stat("/etc/nginx/dhparam.pem")
	if err == nil {
		return
	}

	log.Print("FirstBootDetected! Please await while the web server is configured...")

	vhost, err := valueObject.NewFqdn(os.Getenv("VIRTUAL_HOST"))
	if err != nil {
		log.Fatal("VirtualHostEnvInvalidValue")
	}
	vhostStr := vhost.String()

	log.Print("UpdatingVhost...")

	primaryConfFilePath := "/app/conf/nginx/primary.conf"
	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/speedia.net/"+vhostStr+"/g",
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

	sslCmdRepo := sslInfra.NewSslCmdRepo()
	err = sslCmdRepo.GenerateSelfSignedCert(vhost)
	if err != nil {
		log.Fatal("GenerateSelfSignedCertFailed")
	}

	log.Print("WebServerConfigured!")

	err = servicesInfra.SupervisordFacade{}.Start("nginx")
	if err != nil {
		log.Fatal("StartNginxFailed")
	}
}

func (ws WebServerSetup) OnStartSetup() {
	defaultLogPreffix := "WsOnStartupSetup"

	containerResources, err := o11yInfra.O11yQueryRepo{}.GetOverview()
	if err != nil {
		log.Fatalf("%sGetContainerResourcesFailed", defaultLogPreffix)
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
		log.Fatalf("%sGetNginxWorkersCountFailed", defaultLogPreffix)
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
		log.Fatalf("%sUpdateNginxWorkersCountFailed", defaultLogPreffix)
	}

	err = servicesInfra.SupervisordFacade{}.Restart("nginx")
	if err != nil {
		log.Fatalf("%sRestartNginxFailed", defaultLogPreffix)
	}

	_, err = servicesInfra.ServicesQueryRepo{}.GetByName("php")
	if err == nil {
		err = ws.updatePhpMaxChildProcesses(
			containerResources.HardwareSpecs.MemoryTotal,
		)
		if err != nil {
			log.Fatalf("%s%s", defaultLogPreffix, err.Error())
		}
	}
}
