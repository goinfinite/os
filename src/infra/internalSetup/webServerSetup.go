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
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
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

func (ws *WebServerSetup) generateDhParams() error {
	slog.Info("GeneratingDhParams")

	_, dhparamErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "openssl",
		Args: []string{
			"dhparam", "-dsaparam", "-out", "/etc/nginx/dhparam.pem", "2048",
		},
	}).Run()
	if dhparamErr != nil {
		return errors.New("GenerateDhParamsError: " + dhparamErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) generateSelfSignedCert() error {
	slog.Info("GeneratingSelfSignedCert")

	primaryVirtualHostHostname, readErr :=
		ws.vhostHelpers.ReadPrimaryVirtualHostHostname()
	if readErr != nil {
		return errors.New("ReadPrimaryVirtualHostHostnameError")
	}

	pkiConfDir, pkiErr :=
		tkValueObject.NewUnixAbsoluteFilePath(infraEnvs.PkiConfDir, false)
	if pkiErr != nil {
		return errors.New("PkiConfDirError")
	}

	aliasesHostnames := []tkValueObject.Fqdn{}
	certErr := infraHelper.CreateSelfSignedSsl(
		pkiConfDir, primaryVirtualHostHostname, aliasesHostnames,
	)
	if certErr != nil {
		return errors.New("GenerateSelfSignedCertError: " + certErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) restorePrimaryIndexFile() error {
	restoreErr := infraHelper.RestorePrimaryIndexFile()
	if restoreErr != nil {
		return errors.New("RestorePrimaryIndexFileError: " + restoreErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) generateMappingSecurityRules() error {
	slog.Info("GeneratingMappingSecurityRules")

	mappingCmdRepo := vhostInfra.NewMappingCmdRepo(ws.persistentDbSvc)
	recreateErr := mappingCmdRepo.RecreateSecurityRuleFiles()
	if recreateErr != nil {
		return errors.New("GenerateMappingSecurityRulesError: " + recreateErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) configureNginxAutoStart() error {
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
		if strings.Contains(updateErr.Error(), "Unauthorized") {
			return nil
		}

		return errors.New("ConfigureNginxAutoStartError: " + updateErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) reloadSupervisor() error {
	slog.Info("WebServerConfigured! RestartingServices")

	_, reloadErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "supervisorctl",
		Args:    []string{"-p", "replacedOnFirstBoot", "reload"},
	}).Run()
	if reloadErr != nil {
		return errors.New("ReloadSupervisorError: " + reloadErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) FirstSetup() {
	generateDhParamsErr := ws.generateDhParams()
	if generateDhParamsErr != nil {
		slog.Error(
			"GenerateDhParamsError",
			slog.String("err", generateDhParamsErr.Error()),
		)
		os.Exit(1)
	}

	generateSelfSignedCertErr := ws.generateSelfSignedCert()
	if generateSelfSignedCertErr != nil {
		slog.Error(
			"GenerateSelfSignedCertError",
			slog.String("err", generateSelfSignedCertErr.Error()),
		)
		os.Exit(1)
	}

	restorePrimaryIndexFileErr := ws.restorePrimaryIndexFile()
	if restorePrimaryIndexFileErr != nil {
		slog.Error(
			"RestorePrimaryIndexFileError",
			slog.String("err", restorePrimaryIndexFileErr.Error()),
		)
		os.Exit(1)
	}

	generateMappingSecurityRulesErr := ws.generateMappingSecurityRules()
	if generateMappingSecurityRulesErr != nil {
		slog.Error(
			"GenerateMappingSecurityRulesError",
			slog.String("err", generateMappingSecurityRulesErr.Error()),
		)
		os.Exit(1)
	}

	configureNginxAutoStartErr := ws.configureNginxAutoStart()
	if configureNginxAutoStartErr != nil {
		slog.Error(
			"ConfigureNginxAutoStartError",
			slog.String("err", configureNginxAutoStartErr.Error()),
		)
		os.Exit(1)
	}

	reloadSupervisorErr := ws.reloadSupervisor()
	if reloadSupervisorErr != nil {
		slog.Error(
			"ReloadSupervisorError",
			slog.String("err", reloadSupervisorErr.Error()),
		)
		os.Exit(1)
	}
}

func (ws *WebServerSetup) updatePhpMaxChildProcesses(
	memoryTotal tkValueObject.Byte,
) error {
	slog.Debug("UpdatingMaxPhpChildProcessesInit")

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
		return errors.New("UpdatePhpMaxChildProcessesError: " + sedErr.Error())
	}

	slog.Debug("UpdateMaxPhpChildProcessesSuccess")
	return nil
}

func (ws *WebServerSetup) startNginxIfNeeded(
	servicesQueryRepo *servicesInfra.ServicesQueryRepo,
	servicesCmdRepo *servicesInfra.ServicesCmdRepo,
) error {
	serviceName, _ := valueObject.NewServiceName("nginx")
	readRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &serviceName,
	}
	nginxService, readErr := servicesQueryRepo.ReadFirstInstalledItem(readRequestDto)
	if readErr != nil {
		return errors.New("StartNginxIfNeededError: " + readErr.Error())
	}

	if nginxService.Status.String() == "running" {
		return nil
	}

	startErr := servicesCmdRepo.Start(serviceName)
	if startErr != nil {
		return errors.New("StartNginxIfNeededError: " + startErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) configurePhpChildProcesses(
	servicesQueryRepo *servicesInfra.ServicesQueryRepo,
	memoryTotal tkValueObject.Byte,
) error {
	phpWebServerSvcName, _ := valueObject.NewServiceName("php-webserver")
	_, readErr := servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{ServiceName: &phpWebServerSvcName},
	)
	if readErr != nil {
		slog.Debug(
			"SkippingConfigurePhpChildProcesses",
			slog.String("reason", "PhpWebServerNotInstalled"),
		)
		return nil
	}

	skipProcUpdate := false
	envSkipProcUpdate, err := tkVoUtil.InterfaceToBool(
		os.Getenv(infraEnvs.PhpChildProcessesUpdateSkipEnvKey),
	)
	if err == nil && envSkipProcUpdate {
		skipProcUpdate = true
	}

	if skipProcUpdate {
		slog.Debug(
			"SkippingConfigurePhpChildProcesses",
			slog.String("reason", "EnvVarSet"),
		)
		return nil
	}

	err = ws.updatePhpMaxChildProcesses(memoryTotal)
	if err != nil {
		return errors.New("ConfigurePhpChildProcessesError: " + err.Error())
	}

	return nil
}

func (ws *WebServerSetup) OnStartSetup() {
	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(ws.transientDbSvc)
	containerResources, resourcesErr := o11yQueryRepo.ReadOverview(false)
	if resourcesErr != nil {
		slog.Error(
			"ReadContainerResourcesFailed",
			slog.String("err", resourcesErr.Error()),
		)
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
		slog.Error(
			"ReadNginxWorkersCountFailed", slog.String("err", awkErr.Error()),
		)
		os.Exit(1)
	}

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(ws.persistentDbSvc)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(ws.persistentDbSvc)

	configurePhpChildProcessesErr := ws.configurePhpChildProcesses(
		servicesQueryRepo, containerResources.HardwareSpecs.MemoryTotal,
	)
	if configurePhpChildProcessesErr != nil {
		slog.Warn(
			"ConfigurePhpChildProcessesError",
			slog.String("err", configurePhpChildProcessesErr.Error()),
		)
	}

	if workerCount == cpuCoresStr {
		startNginxIfNeededErr := ws.startNginxIfNeeded(servicesQueryRepo, servicesCmdRepo)
		if startNginxIfNeededErr != nil {
			slog.Error(
				"StartNginxIfNeededError",
				slog.String("err", startNginxIfNeededErr.Error()),
			)
			os.Exit(1)
		}

		return
	}

	slog.Debug("UpdatingWebServerWorkerCount")
	updateWorkerCountErr := ws.vhostHelpers.UpdateWebServerWorkerCount(
		cpuCoresStr, servicesCmdRepo,
	)
	if updateWorkerCountErr != nil {
		slog.Error(
			"UpdateWebServerWorkerCountError",
			slog.String("err", updateWorkerCountErr.Error()),
		)
		os.Exit(1)
	}
}
