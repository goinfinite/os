package internalSetupInfra

import (
	"errors"
	"log/slog"
	"os"
	"strconv"

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

func (ws *WebServerSetup) dhParamsGenerator() error {
	slog.Info("GeneratingDhParams")

	_, dhparamErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "openssl",
		Args: []string{
			"dhparam", "-dsaparam", "-out", "/etc/nginx/dhparam.pem", "2048",
		},
	}).Run()
	if dhparamErr != nil {
		return errors.New("DhParamsGeneratorError: " + dhparamErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) selfSignedCertGenerator() error {
	slog.Info("GeneratingSelfSignedCert")

	primaryVirtualHostHostname, readErr :=
		ws.vhostHelpers.ReadPrimaryVirtualHostHostname()
	if readErr != nil {
		return errors.New("ReadPrimaryVirtualHostHostnameError: " + readErr.Error())
	}

	pkiConfDir, pkiErr :=
		tkValueObject.NewUnixAbsoluteFilePath(infraEnvs.PkiConfDir, false)
	if pkiErr != nil {
		return errors.New("PkiConfDirError: " + pkiErr.Error())
	}

	aliasesHostnames := []tkValueObject.Fqdn{}
	certErr := infraHelper.CreateSelfSignedSsl(
		pkiConfDir, primaryVirtualHostHostname, aliasesHostnames,
	)
	if certErr != nil {
		return errors.New("SelfSignedCertGeneratorError: " + certErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) primaryIndexFileRestorer() error {
	restoreErr := infraHelper.RestorePrimaryIndexFile()
	if restoreErr != nil {
		return errors.New("PrimaryIndexFileRestorerError: " + restoreErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) mappingSecurityRulesGenerator() error {
	slog.Info("GeneratingMappingSecurityRules")

	mappingCmdRepo := vhostInfra.NewMappingCmdRepo(ws.persistentDbSvc)
	recreateErr := mappingCmdRepo.RecreateSecurityRuleFiles()
	if recreateErr != nil {
		return errors.New("MappingSecurityRulesGeneratorError: " + recreateErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) webServerAutoStartConfigurator() error {
	slog.Info("ConfiguringWebServerAutoStart")

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
		return errors.New("WebServerAutoStartConfiguratorError: " + updateErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) supervisorReloader() error {
	slog.Info("WebServerConfigured! RestartingServices")

	_, reloadErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "supervisorctl",
		Args:    []string{"-p", "replacedOnFirstBoot", "reload"},
	}).Run()
	if reloadErr != nil {
		return errors.New("SupervisorReloaderError: " + reloadErr.Error())
	}

	return nil
}

type firstSetupStep struct {
	name string
	fn   func() error
}

func (ws *WebServerSetup) firstSetupOrchestrator() {
	steps := []firstSetupStep{
		{name: "DhParamsGeneratorError", fn: ws.dhParamsGenerator},
		{name: "SelfSignedCertGeneratorError", fn: ws.selfSignedCertGenerator},
		{name: "PrimaryIndexFileRestorerError", fn: ws.primaryIndexFileRestorer},
		{name: "MappingSecurityRulesGeneratorError", fn: ws.mappingSecurityRulesGenerator},
		{name: "WebServerAutoStartConfiguratorError", fn: ws.webServerAutoStartConfigurator},
		{name: "SupervisorReloaderError", fn: ws.supervisorReloader},
	}
	for _, step := range steps {
		if err := step.fn(); err != nil {
			slog.Error(step.name, slog.String("err", err.Error()))
			os.Exit(1)
		}
	}
}

func (ws *WebServerSetup) phpMaxChildProcessesUpdater(
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
			infraEnvs.PhpWebServerMainConfFilePath,
		},
	}).Run()
	if sedErr != nil {
		return errors.New("PhpMaxChildProcessesUpdaterError: " + sedErr.Error())
	}

	slog.Debug("UpdateMaxPhpChildProcessesSuccess")
	return nil
}

func (ws *WebServerSetup) webServerRunningEnsurer(
	servicesQueryRepo *servicesInfra.ServicesQueryRepo,
	servicesCmdRepo *servicesInfra.ServicesCmdRepo,
) error {
	serviceName, _ := valueObject.NewServiceName("nginx")
	readRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &serviceName,
	}
	nginxService, readErr := servicesQueryRepo.ReadFirstInstalledItem(readRequestDto)
	if readErr != nil {
		return errors.New("WebServerRunningEnsurerError: " + readErr.Error())
	}

	if nginxService.Status.String() == "running" {
		return nil
	}

	startErr := servicesCmdRepo.Start(serviceName)
	if startErr != nil {
		return errors.New("WebServerRunningEnsurerError: " + startErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) phpChildProcessesConfigurator(
	servicesQueryRepo *servicesInfra.ServicesQueryRepo,
	memoryTotal tkValueObject.Byte,
) error {
	if !servicesQueryRepo.IsInstalled(valueObject.PhpWebServerServiceName) {
		slog.Debug(
			"SkippingPhpChildProcessesConfigurator",
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
			"SkippingPhpChildProcessesConfigurator",
			slog.String("reason", "EnvVarSet"),
		)
		return nil
	}

	err = ws.phpMaxChildProcessesUpdater(memoryTotal)
	if err != nil {
		return errors.New("PhpChildProcessesConfiguratorError: " + err.Error())
	}

	return nil
}

func (ws *WebServerSetup) onStartSetupOrchestrator() {
	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(ws.transientDbSvc)
	containerResources, resourcesErr := o11yQueryRepo.ReadOverview(false)
	if resourcesErr != nil {
		slog.Error(
			"ReadContainerResourcesFailed",
			slog.String("err", resourcesErr.Error()),
		)
		os.Exit(1)
	}

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(ws.persistentDbSvc)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(ws.persistentDbSvc)

	phpChildProcessesConfiguratorErr := ws.phpChildProcessesConfigurator(
		servicesQueryRepo, containerResources.HardwareSpecs.MemoryTotal,
	)
	if phpChildProcessesConfiguratorErr != nil {
		slog.Warn(
			"PhpChildProcessesConfiguratorError",
			slog.String("err", phpChildProcessesConfiguratorErr.Error()),
		)
	}

	webServerRunningEnsurerErr := ws.webServerRunningEnsurer(
		servicesQueryRepo, servicesCmdRepo,
	)
	if webServerRunningEnsurerErr != nil {
		slog.Error(
			"WebServerRunningEnsurerError",
			slog.String("err", webServerRunningEnsurerErr.Error()),
		)
		os.Exit(1)
	}
}
