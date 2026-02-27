package wsInfra

import (
	"errors"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkInfra "github.com/goinfinite/tk/src/infra"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
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

func (ws *WebServerSetup) updatePhpMaxChildProcesses(memoryTotal tkValueObject.Byte) error {
	log.Print("UpdatingMaxPhpChildProcesses...")

	maxChildProcesses := uint64(300)
	childProcessPerGb := uint64(5)

	memoryInGb := memoryTotal.ToGiB()
	desiredChildProcesses := memoryInGb * childProcessPerGb
	if desiredChildProcesses > maxChildProcesses {
		desiredChildProcesses = maxChildProcesses
	}

	desiredChildProcessesStr := strconv.FormatUint(desiredChildProcesses, 10)
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i", "-e",
			"s/PHP_LSAPI_CHILDREN=[0-9]+/PHP_LSAPI_CHILDREN=" + desiredChildProcessesStr + ";/g",
			infraEnvs.PhpWebserverMainConfFilePath,
		},
	}).Run()
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

	log.Print("FirstBootDetected! PleaseAwait...")

	primaryVirtualHostHostname, err := infraHelper.ReadPrimaryVirtualHostHostname()
	if err != nil {
		log.Fatal("PrimaryVirtualHostNotFound")
	}
	primaryHostnameStr := primaryVirtualHostHostname.String()

	log.Print("UpdatingPrimaryVirtualHost...")

	primaryConfFilePath := infraEnvs.VirtualHostsConfDir + "/primary.conf"
	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i",
			"s/" + infraEnvs.DefaultPrimaryVhost + "/" + primaryHostnameStr + "/g",
			primaryConfFilePath,
		},
	}).Run()
	if err != nil {
		log.Fatal("UpdatePrimaryVirtualHostFileFailed")
	}

	log.Print("GeneratingDhParams...")

	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "openssl",
		Args: []string{
			"dhparam", "-dsaparam", "-out", "/etc/nginx/dhparam.pem", "2048",
		},
	}).Run()
	if err != nil {
		log.Fatal("GenerateDhparamFailed")
	}

	log.Print("GeneratingSelfSignedCert...")

	pkiConfDir, err := valueObject.NewUnixFilePath(infraEnvs.PkiConfDir)
	if err != nil {
		log.Fatal("PkiConfDirNotFound")
	}

	aliasesHostnames := []tkValueObject.Fqdn{}
	err = infraHelper.CreateSelfSignedSsl(
		pkiConfDir, primaryVirtualHostHostname, aliasesHostnames,
	)
	if err != nil {
		log.Fatal("GenerateSelfSignedCertFailed: ", err.Error())
	}

	err = infraHelper.RestorePrimaryIndexFile()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Print("GenerateMappingSecurityRules...")

	mappingCmdRepo := vhostInfra.NewMappingCmdRepo(ws.persistentDbSvc)
	err = mappingCmdRepo.RecreateSecurityRuleFiles()
	if err != nil {
		log.Fatal("GenerateMappingSecurityRulesFailed: ", err.Error())
	}

	log.Print("ConfiguringWebServerAutoStart...")

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(ws.persistentDbSvc)
	nginxServiceName, _ := valueObject.NewServiceName("nginx")
	nginxAutoStart := true
	updateServiceDto := dto.NewUpdateService(
		nginxServiceName, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, &nginxAutoStart, nil, nil, nil, nil, nil, nil,
		tkValueObject.AccountIdSystem, tkValueObject.IpAddressLocal,
	)
	err = servicesCmdRepo.Update(updateServiceDto)
	if err != nil {
		if !strings.Contains(err.Error(), "Unauthorized") {
			log.Fatal("UpdateNginxAutoStartFailed: ", err.Error())
		}
	}

	log.Print("WebServerConfigured! RestartingServices...")

	// Do not write any code after this as supervisorctl reload will restart
	// the OS API and any remaining code will not be executed.
	_, _ = tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "supervisorctl",
		Args:    []string{"-p", "replacedOnFirstBoot", "reload"},
	}).Run()
}

func (ws *WebServerSetup) OnStartSetup() {
	defaultLogPrefix := "WsOnStartupSetup"

	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(ws.transientDbSvc)
	containerResources, err := o11yQueryRepo.ReadOverview(false)
	if err != nil {
		log.Fatalf("%sGetContainerResourcesFailed", defaultLogPrefix)
	}

	cpuCores := containerResources.HardwareSpecs.CpuCores
	cpuCoresStr := strconv.FormatInt(int64(math.Ceil(cpuCores)), 10)

	nginxConfFilePath := "/etc/nginx/nginx.conf"
	workerCount, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "awk",
		Args: []string{
			"/worker_processes/{gsub(/[^0-9]+/, \"\"); print}", nginxConfFilePath,
		},
	}).Run()
	if err != nil {
		log.Fatalf("%sGetNginxWorkersCountFailed", defaultLogPrefix)
	}

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(ws.persistentDbSvc)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(ws.persistentDbSvc)
	serviceName, _ := valueObject.NewServiceName("nginx")
	readFirstInstalledServiceRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &serviceName,
	}
	if workerCount == cpuCoresStr {
		nginxService, err := servicesQueryRepo.ReadFirstInstalledItem(
			readFirstInstalledServiceRequestDto,
		)
		if err != nil {
			log.Fatalf("ReadNginxServiceFailed: %s", err.Error())
		}

		if nginxService.Status.String() == "running" {
			return
		}

		err = servicesCmdRepo.Start(serviceName)
		if err != nil {
			log.Fatalf("StartNginxServiceFailed: %s", err.Error())
		}
		return
	}

	log.Print("UpdatingNginxWorkersCount...")

	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i", "-e",
			"s/^worker_processes.*/worker_processes " + cpuCoresStr + ";/g",
			nginxConfFilePath,
		},
	}).Run()
	if err != nil {
		log.Fatalf("%sUpdateNginxWorkersCountFailed", defaultLogPrefix)
	}

	err = servicesCmdRepo.Restart(serviceName)
	if err != nil {
		log.Fatalf("%sRestartNginxFailed", defaultLogPrefix)
	}

	phpWebServerSvcName, _ := valueObject.NewServiceName("php-webserver")
	readFirstInstalledServiceRequestDto.ServiceName = &phpWebServerSvcName
	_, err = servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledServiceRequestDto,
	)
	if err == nil {
		err = ws.updatePhpMaxChildProcesses(
			containerResources.HardwareSpecs.MemoryTotal,
		)
		if err != nil {
			log.Fatalf("%s%s", defaultLogPrefix, err.Error())
		}
	}
}
