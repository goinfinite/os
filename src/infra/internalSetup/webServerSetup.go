package internalSetupInfra

import (
	"errors"
	"log/slog"
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
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
)

type WebServerSetup struct {
	persistentDbSvc   *internalDbInfra.PersistentDatabaseService
	transientDbSvc    *tkInfraDb.TransientDatabaseService
	servicesQueryRepo *servicesInfra.ServicesQueryRepo
	servicesCmdRepo   *servicesInfra.ServicesCmdRepo
	vhostHelpers      *vhostInfra.VirtualHostHelpers
}

func NewWebServerSetup(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *tkInfraDb.TransientDatabaseService,
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

// The main web server is configured with auto-start disabled by default.
// This method enables auto-start in the database after mapping and
// configuration complete successfully.
func (ws *WebServerSetup) webServerAutoStartConfigurator() error {
	_, readErr := ws.servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{
			ServiceName: &valueObject.ServiceNameMainWebServer,
		},
	)
	if readErr != nil {
		return errors.New("MainWebServerServiceNotFound: " + readErr.Error())
	}

	serviceAutoStart := true
	updateServiceDto := dto.UpdateService{
		Name:              valueObject.ServiceNameMainWebServer,
		AutoStart:         &serviceAutoStart,
		OperatorAccountId: tkValueObject.AccountIdSystem,
		OperatorIpAddress: tkValueObject.IpAddressLocal,
	}
	err := ws.servicesCmdRepo.Update(updateServiceDto)
	if err != nil {
		// ProcessManagerAuthError is expected until processManagerStateGuarantor is
		// called. strings.Contains is used as the error message may be appended with
		// additional context.
		if strings.Contains(err.Error(), servicesInfra.ErrProcessManagerAuthError.Error()) {
			return nil
		}

		return errors.New("MainWebServerAutoStartUpdateError: " + err.Error())
	}

	return nil
}

// Important: Reloading the process manager (PID 1) restarts the container.
// Do not expect any code after this method to execute.
func (ws *WebServerSetup) processManagerStateGuarantor() error {
	fileClerk := tkInfra.FileClerk{}
	err := fileClerk.CreateSymlink(
		infraEnvs.ProcessManagerConfFilePath, "/etc/supervisord.conf", true,
	)
	if err != nil {
		return errors.New("ProcessManagerConfSymlinkError: " + err.Error())
	}

	err = ws.servicesCmdRepo.ProcessManagerConfRebuilder()
	if err != nil {
		if !errors.Is(err, servicesInfra.ErrProcessManagerAuthError) {
			return errors.New("ProcessManagerConfRebuilderError: " + err.Error())
		}

		_, reloadErr := tkInfra.NewShell(tkInfra.ShellSettings{
			Command:          infraEnvs.ProcessManagerBinaryPath,
			Args:             []string{"-p", "replacedOnFirstBoot", "reload"},
			WorkingDirectory: infraEnvs.InfiniteOsMainDir,
		}).Run()
		if reloadErr != nil {
			return errors.New("ProcessManagerReloadError: " + reloadErr.Error())
		}
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
			errorMessage: "ProcessManagerStateGuarantorError",
			executeFn:    ws.processManagerStateGuarantor,
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

func (ws *WebServerSetup) phpLsapiBlockUpdateExpressionFactory(
	substitutionExpression string,
) string {
	return "/^extprocessor lsphp[0-9]/,/^}/ " + substitutionExpression
}

func (ws *WebServerSetup) phpWebServerCapacityFileUpdater(
	confFilePath string,
	lsapiCapacity uint64,
	workersCount uint64,
) error {
	lsapiCapacityStr := strconv.FormatUint(lsapiCapacity, 10)
	workersCountStr := strconv.FormatUint(workersCount, 10)
	autoUpdateComment := "# AUTO CALCULATED. DO NOT EDIT. LAST EDIT: " +
		tkValueObject.NewUnixTimeNow().ReadRfcDate()

	childrenSetting := "PHP_LSAPI_CHILDREN=" + lsapiCapacityStr + "; " +
		autoUpdateComment
	connectionsSetting := "maxConns " + lsapiCapacityStr + "; " +
		autoUpdateComment
	workersSetting := "httpdworkers " + workersCountStr + "; " +
		autoUpdateComment

	childrenUpdateExpression := ws.phpLsapiBlockUpdateExpressionFactory(
		"s/PHP_LSAPI_CHILDREN=[0-9]+.*/" + childrenSetting + "/",
	)
	connectionsUpdateExpression := ws.phpLsapiBlockUpdateExpressionFactory(
		"s/^([[:space:]]*)maxConns[[:space:]]+[0-9]+.*/" +
			"\\1" + connectionsSetting + "/",
	)
	workersUpdateExpression := "s/^httpdworkers[[:space:]]+[0-9]+.*/" +
		workersSetting + "/"

	_, sedErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i", "-E",
			"-e", childrenUpdateExpression,
			"-e", connectionsUpdateExpression,
			"-e", workersUpdateExpression,
			confFilePath,
		},
	}).Run()
	if sedErr != nil {
		return errors.New("PhpWebServerCapacityFileUpdaterError: " + sedErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) phpWebServerCapacityCalculator(
	memoryTotal tkValueObject.Byte,
	cpuCores float64,
) (lsapiCapacity uint64, workersCount uint64) {
	openLiteSpeedCoresPerWorker := uint64(4)
	openLiteSpeedMaxWorkersCount := uint64(16)

	capacityHardCap := uint64(300)
	capacityPerGb := uint64(5)
	memoryGiB := max(memoryTotal.ToGiB(), 1)
	lsapiCapacity = min(memoryGiB*capacityPerGb, capacityHardCap)

	workersCount = min(
		max(uint64(cpuCores)/openLiteSpeedCoresPerWorker, 1),
		openLiteSpeedMaxWorkersCount,
	)

	return lsapiCapacity, workersCount
}

func (ws *WebServerSetup) webServerRunningEnsurer() error {
	mainWebServerService, readErr := ws.servicesQueryRepo.ReadFirstInstalledItem(
		dto.ReadFirstInstalledServiceItemsRequest{
			ServiceName: &valueObject.ServiceNameMainWebServer,
		},
	)
	if readErr != nil {
		return errors.New("MainWebServerServiceNotFound: " + readErr.Error())
	}

	if mainWebServerService.Status == valueObject.ServiceStatusRunning {
		return nil
	}

	startErr := ws.servicesCmdRepo.Start(valueObject.ServiceNameMainWebServer)
	if startErr != nil {
		return errors.New("WebServerStartError: " + startErr.Error())
	}

	return nil
}

func (ws *WebServerSetup) phpWebServerCapacityConfigurator(
	memoryTotal tkValueObject.Byte,
	cpuCores float64,
) error {
	phpWebServerIsInstalled, err := ws.servicesQueryRepo.IsInstalled(
		valueObject.ServiceNamePhpWebServer,
	)
	if err != nil {
		slog.Warn(
			"SkippingPhpWebServerCapacityConfigurator",
			slog.String("reason", "PhpWebServerInstallationCheckFailed"),
			slog.String("err", err.Error()),
		)
		return nil
	}
	if !phpWebServerIsInstalled {
		slog.Debug(
			"SkippingPhpWebServerCapacityConfigurator",
			slog.String("reason", "PhpWebServerNotInstalled"),
		)
		return nil
	}

	skipCapacityUpdate := false
	envSkipCapacityUpdate, err := tkVoUtil.InterfaceToBool(
		os.Getenv(infraEnvs.PhpWebServerCapacityUpdateSkipEnvKey),
	)
	if err == nil && envSkipCapacityUpdate {
		skipCapacityUpdate = true
	}

	if skipCapacityUpdate {
		slog.Debug(
			"SkippingPhpWebServerCapacityConfigurator",
			slog.String("reason", "EnvVarSet"),
		)
		return nil
	}

	lsapiCapacity, workersCount := ws.phpWebServerCapacityCalculator(
		memoryTotal, cpuCores,
	)

	err = ws.phpWebServerCapacityFileUpdater(
		infraEnvs.PhpWebServerMainConfFilePath, lsapiCapacity, workersCount,
	)
	if err != nil {
		return errors.New("PhpWebServerCapacityConfiguratorError: " + err.Error())
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

	phpWebServerCapacityConfiguratorErr := ws.phpWebServerCapacityConfigurator(
		containerResources.HardwareSpecs.MemoryTotal,
		containerResources.HardwareSpecs.CpuCores,
	)
	if phpWebServerCapacityConfiguratorErr != nil {
		slog.Warn(
			"PhpWebServerCapacityConfiguratorError",
			slog.String("err", phpWebServerCapacityConfiguratorErr.Error()),
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
