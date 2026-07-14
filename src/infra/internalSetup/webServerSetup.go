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
	persistentDbSvc   *internalDbInfra.PersistentDatabaseService
	transientDbSvc    *internalDbInfra.TransientDatabaseService
	servicesQueryRepo *servicesInfra.ServicesQueryRepo
	servicesCmdRepo   *servicesInfra.ServicesCmdRepo
	vhostHelpers      *vhostInfra.VirtualHostHelpers
}

func NewWebServerSetup(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *WebServerSetup {
	return &WebServerSetup{
		persistentDbSvc:   persistentDbSvc,
		transientDbSvc:    transientDbSvc,
		servicesQueryRepo: servicesInfra.NewServicesQueryRepo(persistentDbSvc),
		servicesCmdRepo:   servicesInfra.NewServicesCmdRepo(persistentDbSvc),
		vhostHelpers:      vhostInfra.NewVirtualHostHelpers(),
	}
}

func (ws *WebServerSetup) dhParamsGenerator() error {
	_, dhparamErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "openssl",
		Args: []string{
			"dhparam", "-dsaparam", "-out", infraEnvs.WebServerDhParamFilePath, "2048",
		},
	}).Run()
	if dhparamErr != nil {
		return errors.New("DhParamsGeneratorError: " + dhparamErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) selfSignedCertGenerator() error {
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
	mappingCmdRepo := vhostInfra.NewMappingCmdRepo(ws.persistentDbSvc)
	recreateErr := mappingCmdRepo.RecreateSecurityRuleFiles()
	if recreateErr != nil {
		return errors.New("MappingSecurityRulesGeneratorError: " + recreateErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) webServerAutoStartConfigurator() error {
	_, readErr := ws.servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{
			ServiceName: &valueObject.MainWebServerServiceName,
		},
	)
	if readErr != nil {
		return errors.New("MainWebServerServiceNotFound: " + readErr.Error())
	}

	serviceAutoStart := true
	updateServiceDto := dto.NewUpdateService(
		valueObject.MainWebServerServiceName, nil, nil, nil, nil, nil, nil, nil,
		nil, nil, nil, nil, nil, nil, nil, &serviceAutoStart, nil, nil, nil, nil,
		nil, nil, tkValueObject.AccountIdSystem, tkValueObject.IpAddressLocal,
	)
	updateErr := ws.servicesCmdRepo.Update(updateServiceDto)
	if updateErr != nil {
		return errors.New("MainWebServerAutoStartUpdateError: " + updateErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) processManagerReloader() error {
	fileClerk := tkInfra.FileClerk{}
	err := fileClerk.CreateSymlink(
		infraEnvs.ProcessManagerConfFilePath, "/etc/supervisord.conf", true,
	)
	if err != nil {
		return errors.New("ProcessManagerConfSymlinkError: " + err.Error())
	}

	_, reloadErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:          infraEnvs.ProcessManagerBinaryPath,
		Args:             []string{"-p", "replacedOnFirstBoot", "reload"},
		WorkingDirectory: infraEnvs.InfiniteOsMainDir,
	}).Run()
	if reloadErr != nil {
		return errors.New("ProcessManagerReloaderError: " + reloadErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) firstSetupOrchestrator() {
	type firstSetupStep struct {
		errorMessage string
		executeFn    func() error
	}

	setupSteps := []firstSetupStep{
		{
			errorMessage: "DhParamsGeneratorError",
			executeFn:    ws.dhParamsGenerator,
		},
		{
			errorMessage: "SelfSignedCertGeneratorError",
			executeFn:    ws.selfSignedCertGenerator,
		},
		{
			errorMessage: "PrimaryIndexFileRestorerError",
			executeFn:    ws.primaryIndexFileRestorer,
		},
		{
			errorMessage: "MappingSecurityRulesGeneratorError",
			executeFn:    ws.mappingSecurityRulesGenerator,
		},
		{
			errorMessage: "WebServerAutoStartConfiguratorError",
			executeFn:    ws.webServerAutoStartConfigurator,
		},
		{
			errorMessage: "ProcessManagerReloaderError",
			executeFn:    ws.processManagerReloader,
		},
	}
	for _, setupStep := range setupSteps {
		executeErr := setupStep.executeFn()
		if executeErr != nil {
			slog.Error(
				setupStep.errorMessage,
				slog.String("err", executeErr.Error()),
			)
			os.Exit(1)
		}
	}
}

func (ws *WebServerSetup) phpMaxChildProcessesUpdater(
	memoryTotal tkValueObject.Byte,
) error {
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

	return nil
}

func (ws *WebServerSetup) webServerRunningEnsurer() error {
	mainWebServerService, readErr := ws.servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{
			ServiceName: &valueObject.MainWebServerServiceName,
		},
	)
	if readErr != nil {
		return errors.New("MainWebServerServiceNotFound: " + readErr.Error())
	}

	if mainWebServerService.Status == valueObject.ServiceStatusRunning {
		return nil
	}

	startErr := ws.servicesCmdRepo.Start(valueObject.MainWebServerServiceName)
	if startErr != nil {
		return errors.New("WebServerStartError: " + startErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) phpChildProcessesConfigurator(
	memoryTotal tkValueObject.Byte,
) error {
	if !ws.servicesQueryRepo.IsInstalled(valueObject.PhpWebServerServiceName) {
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

	phpChildProcessesConfiguratorErr := ws.phpChildProcessesConfigurator(
		containerResources.HardwareSpecs.MemoryTotal,
	)
	if phpChildProcessesConfiguratorErr != nil {
		slog.Warn(
			"PhpChildProcessesConfiguratorError",
			slog.String("err", phpChildProcessesConfiguratorErr.Error()),
		)
	}

	webServerRunningEnsurerErr := ws.webServerRunningEnsurer()
	if webServerRunningEnsurerErr != nil {
		slog.Error(
			"WebServerRunningEnsurerError",
			slog.String("err", webServerRunningEnsurerErr.Error()),
		)
		os.Exit(1)
	}
}
