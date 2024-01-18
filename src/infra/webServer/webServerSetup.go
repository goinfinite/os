package webServerInfra

import (
	"log"
	"os"
	"strconv"

	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func updatePhpMaxChildProcesses(
	memoryTotal valueObject.Byte,
) {
	log.Print("UpdatingMaxChildProcesses...")

	maxChildProcesses := int64(300)
	childProcessPerGb := int64(5)

	memoryInGb := memoryTotal.ToGiB()
	targetChildProcesses := memoryInGb * childProcessPerGb
	if targetChildProcesses > maxChildProcesses {
		targetChildProcesses = maxChildProcesses
	}
	targetChildProcessesStr := strconv.FormatInt(targetChildProcesses, 10)

	httpdConfFilePath := "/usr/local/lsws/conf/httpd_config.conf"
	_, err := infraHelper.RunCmd(
		"sed",
		"-i",
		"-e",
		"s/PHP_LSAPI_CHILDREN=[0-9]+/PHP_LSAPI_CHILDREN="+targetChildProcessesStr+";/g",
		httpdConfFilePath,
	)
	if err != nil {
		log.Fatal("WsOnStartupSetupUpdateMaxChildProcessesFailed")
	}
}

func WebServerFirstSetup() {
	_, err := os.Stat("/etc/nginx/dhparam.pem")
	if err == nil {
		return
	}

	log.Print("FirstBootDetected! Please await while the web server is configured...")

	vhost, err := valueObject.NewFqdn(os.Getenv("VIRTUAL_HOST"))
	if err != nil {
		log.Fatalf("VirtualHostEnvInvalidValue")
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
		log.Fatalf("UpdateVhostFailed")
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
		log.Fatalf("GenerateDhparamFailed")
	}

	log.Print("GeneratingSelfSignedCert...")

	_, err = infraHelper.RunCmd(
		"openssl",
		"req",
		"-x509",
		"-nodes",
		"-days",
		"365",
		"-newkey",
		"rsa:2048",
		"-keyout",
		"/app/conf/pki/"+vhostStr+".key",
		"-out",
		"/app/conf/pki/"+vhostStr+".crt",
		"-subj",
		"/C=US/ST=California/L=LosAngeles/O=Acme/CN="+vhostStr,
	)
	if err != nil {
		log.Fatalf("GenerateSelfSignedCertFailed")
	}

	log.Print("WebServerConfigured!")

	err = servicesInfra.SupervisordFacade{}.Start("nginx")
	if err != nil {
		log.Fatalf("StartNginxFailed")
	}
}

func WebServerOnStartSetup() {
	containerResources, err := infra.O11yQueryRepo{}.GetOverview()
	if err != nil {
		log.Fatalf("WsOnStartupSetupGetContainerResourcesFailed")
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
		log.Fatalf("WsOnStartupSetupGetNginxWorkersCountFailed")
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
		log.Fatalf("WsOnStartupSetupUpdateNginxWorkersCountFailed")
	}

	err = servicesInfra.SupervisordFacade{}.Restart("nginx")
	if err != nil {
		log.Fatalf("WsOnStartupSetupRestartNginxFailed")
	}

	_, err = servicesInfra.ServicesQueryRepo{}.GetByName("php")
	if err == nil {
		updatePhpMaxChildProcesses(
			containerResources.HardwareSpecs.MemoryTotal,
		)
	}
}
