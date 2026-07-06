package internalSetupInfra

import (
	"errors"
	"log/slog"
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
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

type WebServerSetup struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	vhostHelpers    *vhostInfra.VirtualHostHelpers
}

func NewWebServerSetup(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *WebServerSetup {
	return &WebServerSetup{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		vhostHelpers:    vhostInfra.NewVirtualHostHelpers(),
	}
}

func (ws *WebServerSetup) generateDhParams() {
	slog.Info("GeneratingDhParams")

	_, dhparamErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "openssl",
		Args: []string{
			"dhparam", "-dsaparam", "-out", "/etc/nginx/dhparam.pem", "2048",
		},
	}).Run()
	if dhparamErr != nil {
		slog.Error("GenerateDhparamFailed", slog.String("err", dhparamErr.Error()))
		os.Exit(1)
	}
}

func (ws *WebServerSetup) generateSelfSignedCert() {
	slog.Info("GeneratingSelfSignedCert")

	primaryVirtualHostHostname, readErr := ws.vhostHelpers.ReadPrimaryVirtualHostHostname()
	if readErr != nil {
		slog.Error("PrimaryVirtualHostNotFound")
		os.Exit(1)
	}

	pkiConfDir, pkiErr := tkValueObject.NewUnixAbsoluteFilePath(infraEnvs.PkiConfDir, false)
	if pkiErr != nil {
		slog.Error("PkiConfDirNotFound")
		os.Exit(1)
	}

	aliasesHostnames := []tkValueObject.Fqdn{}
	certErr := infraHelper.CreateSelfSignedSsl(
		pkiConfDir, primaryVirtualHostHostname, aliasesHostnames,
	)
	if certErr != nil {
		slog.Error("GenerateSelfSignedCertFailed", slog.String("err", certErr.Error()))
		os.Exit(1)
	}
}

func (ws *WebServerSetup) restorePrimaryIndexFile() {
	restoreErr := infraHelper.RestorePrimaryIndexFile()
	if restoreErr != nil {
		slog.Error("RestorePrimaryIndexFileFailed", slog.String("err", restoreErr.Error()))
		os.Exit(1)
	}
}

func (ws *WebServerSetup) generateMappingSecurityRules() {
	slog.Info("GeneratingMappingSecurityRules")

	mappingCmdRepo := vhostInfra.NewMappingCmdRepo(ws.persistentDbSvc)
	recreateErr := mappingCmdRepo.RecreateSecurityRuleFiles()
	if recreateErr != nil {
		slog.Error("GenerateMappingSecurityRulesFailed", slog.String("err", recreateErr.Error()))
		os.Exit(1)
	}
}

func (ws *WebServerSetup) configureNginxAutoStart() {
	slog.Info("ConfiguringNginxAutoStart")

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(ws.persistentDbSvc)
	nginxServiceName, _ := valueObject.NewServiceName("nginx")
	nginxAutoStart := true
	updateServiceDto := dto.NewUpdateService(
		nginxServiceName, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, &nginxAutoStart, nil, nil, nil, nil, nil, nil,
		tkValueObject.AccountIdSystem, tkValueObject.IpAddressLocal,
	)
	updateErr := servicesCmdRepo.Update(updateServiceDto)
	if updateErr != nil {
		if !strings.Contains(updateErr.Error(), "Unauthorized") {
			slog.Error("UpdateNginxAutoStartFailed", slog.String("err", updateErr.Error()))
			os.Exit(1)
		}
	}
}

func (ws *WebServerSetup) reloadSupervisor() {
	slog.Info("WebServerConfigured! RestartingServices")

	_, reloadErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "supervisorctl",
		Args:    []string{"-p", "replacedOnFirstBoot", "reload"},
	}).Run()
	if reloadErr != nil {
		slog.Error("ReloadSupervisorFailed", slog.String("err", reloadErr.Error()))
		os.Exit(1)
	}
}

func (ws *WebServerSetup) FirstSetup() {
	err := ws.vhostHelpers.UpdatePrimaryVirtualHostPlaceholder()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ws.generateDhParams()
	ws.generateSelfSignedCert()
	ws.restorePrimaryIndexFile()
	ws.generateMappingSecurityRules()
	ws.configureNginxAutoStart()
	ws.reloadSupervisor()
}

func (ws *WebServerSetup) updatePhpMaxChildProcesses(memoryTotal tkValueObject.Byte) error {
	slog.Info("UpdatingMaxPhpChildProcessesInit")

	childProcsHardCap := uint64(300)
	childProcsPerGb := uint64(5)

	childProcsHealthyAmount := memoryTotal.ToGiB() * childProcsPerGb
	if childProcsHealthyAmount > childProcsHardCap {
		childProcsHealthyAmount = childProcsHardCap
	}

	childProcsHealthyAmountStr := strconv.FormatUint(childProcsHealthyAmount, 10)
	autoUpdateComment := "# AUTO CALCULATED. DO NOT EDIT. LAST EDIT: " +
		tkValueObject.NewUnixTimeNow().ReadRfcDate()
	childProcsNewValue := "PHP_LSAPI_CHILDREN=" + childProcsHealthyAmountStr + "; " +
		autoUpdateComment

	_, sedErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i", "-E",
			"s/PHP_LSAPI_CHILDREN=[0-9]+.*/" + childProcsNewValue + "/g",
			infraEnvs.PhpWebserverMainConfFilePath,
		},
	}).Run()
	if sedErr != nil {
		return errors.New("UpdatePhpMaxChildProcessesFailed: " + sedErr.Error())
	}

	slog.Info("UpdateMaxPhpChildProcessesSuccess")
	return nil
}

func (ws *WebServerSetup) startNginxIfNeeded(
	servicesQueryRepo *servicesInfra.ServicesQueryRepo,
	servicesCmdRepo *servicesInfra.ServicesCmdRepo,
) {
	serviceName, _ := valueObject.NewServiceName("nginx")
	readRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &serviceName,
	}
	nginxService, readErr := servicesQueryRepo.ReadFirstInstalledItem(readRequestDto)
	if readErr != nil {
		slog.Error("ReadNginxServiceFailed", slog.String("err", readErr.Error()))
		os.Exit(1)
	}

	if nginxService.Status.String() == "running" {
		return
	}

	startErr := servicesCmdRepo.Start(serviceName)
	if startErr != nil {
		slog.Error("StartNginxServiceFailed", slog.String("err", startErr.Error()))
		os.Exit(1)
	}
}

func (ws *WebServerSetup) configurePhpChildProcesses(
	servicesQueryRepo *servicesInfra.ServicesQueryRepo,
	memoryTotal tkValueObject.Byte,
) {
	phpWebServerSvcName, _ := valueObject.NewServiceName("php-webserver")
	_, readErr := servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{ServiceName: &phpWebServerSvcName},
	)
	if readErr != nil {
		slog.Debug("PhpWebServerNotInstalled. SkippingConfigurePhpChildProcesses")
		return
	}

	err := ws.updatePhpMaxChildProcesses(memoryTotal)
	if err != nil {
		slog.Error("ConfigurePhpChildProcessesFailed", slog.String("err", err.Error()))
		return
	}
}

func (ws *WebServerSetup) OnStartSetup() {
	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(ws.transientDbSvc)
	containerResources, resourcesErr := o11yQueryRepo.ReadOverview(false)
	if resourcesErr != nil {
		slog.Error("ReadContainerResourcesFailed", slog.String("err", resourcesErr.Error()))
		os.Exit(1)
	}

	cpuCores := containerResources.HardwareSpecs.CpuCores
	cpuCoresStr := strconv.FormatInt(int64(math.Ceil(cpuCores)), 10)

	nginxConfFilePath := "/etc/nginx/nginx.conf"
	workerCount, awkErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "awk",
		Args: []string{
			"/worker_processes/{gsub(/[^0-9]+/, \"\"); print}", nginxConfFilePath,
		},
	}).Run()
	if awkErr != nil {
		slog.Error("ReadNginxWorkersCountFailed", slog.String("err", awkErr.Error()))
		os.Exit(1)
	}

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(ws.persistentDbSvc)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(ws.persistentDbSvc)

	ws.configurePhpChildProcesses(
		servicesQueryRepo, containerResources.HardwareSpecs.MemoryTotal,
	)

	if workerCount == cpuCoresStr {
		ws.startNginxIfNeeded(servicesQueryRepo, servicesCmdRepo)
		return
	}
	err := ws.vhostHelpers.UpdateWebServerWorkerCount(cpuCoresStr, servicesCmdRepo)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
