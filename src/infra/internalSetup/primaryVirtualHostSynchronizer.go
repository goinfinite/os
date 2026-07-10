package internalSetupInfra

import (
	"errors"
	"log/slog"
	"os"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	runtimeInfra "github.com/goinfinite/os/src/infra/runtime"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

var phpWebServerServiceName, phpWebServerServiceNameError = valueObject.NewServiceName(
	"php-webserver",
)

type PrimaryVirtualHostSynchronizer struct {
	persistentDbSvc         *internalDbInfra.PersistentDatabaseService
	previousPrimaryHostname tkValueObject.Fqdn
	newPrimaryHostname      tkValueObject.Fqdn
	servicesQueryRepo       *servicesInfra.ServicesQueryRepo
	runtimeCmdRepo          *runtimeInfra.RuntimeCmdRepo
	vhostCmdRepo            *vhostInfra.VirtualHostCmdRepo
	vhostHelpers            *vhostInfra.VirtualHostHelpers
}

func NewPrimaryVirtualHostSynchronizer(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *PrimaryVirtualHostSynchronizer {
	return &PrimaryVirtualHostSynchronizer{
		persistentDbSvc:   persistentDbSvc,
		servicesQueryRepo: servicesInfra.NewServicesQueryRepo(persistentDbSvc),
		runtimeCmdRepo:    runtimeInfra.NewRuntimeCmdRepo(persistentDbSvc),
		vhostCmdRepo:      vhostInfra.NewVirtualHostCmdRepo(persistentDbSvc),
		vhostHelpers:      vhostInfra.NewVirtualHostHelpers(),
	}
}

func (sync *PrimaryVirtualHostSynchronizer) phpConfUpdater() error {
	if !sync.servicesQueryRepo.IsInstalled(phpWebServerServiceName) {
		slog.Debug(
			"SkippingPrimaryVirtualHostPhpConfUpdater",
			slog.String("reason", "PhpWebServerNotInstalled"),
		)
		return nil
	}

	return sync.runtimeCmdRepo.UpdatePhpVirtualHostHostname(
		sync.previousPrimaryHostname, sync.newPrimaryHostname,
	)
}

func (sync *PrimaryVirtualHostSynchronizer) confUpdater() error {
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(sync.persistentDbSvc)
	vhostEntity, readErr := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &sync.previousPrimaryHostname,
	})
	if readErr != nil {
		return errors.New("ReadPreviousVirtualHostFailed: " + readErr.Error())
	}

	mappingsFilePath, readErr := vhostQueryRepo.ReadVirtualHostMappingsFilePath(
		sync.previousPrimaryHostname,
	)
	if readErr != nil {
		return errors.New("ReadMappingsFilePathFailed: " + readErr.Error())
	}

	pkiConfDir, parseErr := tkValueObject.NewUnixAbsoluteFilePath(
		infraEnvs.PkiConfDir, false,
	)
	if parseErr != nil {
		return errors.New("InvalidPkiConfDir: " + parseErr.Error())
	}

	createCertErr := infraHelper.CreateSelfSignedSsl(
		pkiConfDir, sync.newPrimaryHostname, vhostEntity.AliasesHostnames,
	)
	if createCertErr != nil {
		return errors.New("CreateSelfSignedSslFailed: " + createCertErr.Error())
	}

	modifiedEntity := vhostEntity
	modifiedEntity.Hostname = sync.newPrimaryHostname

	confContent, factoryErr := sync.vhostCmdRepo.WebServerUnitFileFactory(
		modifiedEntity, mappingsFilePath,
	)
	if factoryErr != nil {
		return errors.New("WebServerUnitFileFactoryFailed: " + factoryErr.Error())
	}

	fileClerk := tkInfra.FileClerk{}
	writeErr := fileClerk.UpdateFileContent(
		infraEnvs.PrimaryVirtualHostConfPath, confContent, true,
	)
	if writeErr != nil {
		return errors.New("WritePrimaryVirtualHostConfFailed: " + writeErr.Error())
	}

	return sync.vhostHelpers.ReloadWebServer()
}

func (sync *PrimaryVirtualHostSynchronizer) dbUpdater() error {
	var primaryVirtualHostModel dbModel.VirtualHost
	queryErr := sync.persistentDbSvc.Handler.Where(
		"is_primary = ?", true,
	).First(&primaryVirtualHostModel).Error
	if queryErr != nil {
		return errors.New("DbReadFailed: " + queryErr.Error())
	}

	if primaryVirtualHostModel.Hostname == sync.newPrimaryHostname.String() {
		return nil
	}

	updateErr := sync.persistentDbSvc.Handler.Model(
		&primaryVirtualHostModel,
	).Update("hostname", sync.newPrimaryHostname.String()).Error
	if updateErr != nil {
		return errors.New("DbWriteFailed: " + updateErr.Error())
	}

	return nil
}

func (sync *PrimaryVirtualHostSynchronizer) Run() {
	rawPrimaryVirtualHostEnvValue := os.Getenv(infraEnvs.PrimaryVirtualHostEnvKey)
	if rawPrimaryVirtualHostEnvValue == "" {
		slog.Debug(
			"SkippingPrimaryVirtualHostSynchronizer",
			slog.String("reason", "EmptyPrimaryVirtualHostEnvValue"),
		)
		return
	}

	primaryVirtualHostEnvValue, parseErr := tkValueObject.NewFqdn(rawPrimaryVirtualHostEnvValue)
	if parseErr != nil {
		slog.Error(
			"InvalidPrimaryVirtualHostEnvValue", slog.String("err", parseErr.Error()),
		)
		return
	}

	webServerPrimaryVirtualHost, err := sync.vhostHelpers.
		ReadPrimaryVirtualHostHostnameFromWebServerConf()
	if err != nil {
		slog.Error(
			"ReadPrimaryVirtualHostHostnameFailed", slog.String("err", err.Error()),
		)
		return
	}

	if primaryVirtualHostEnvValue == webServerPrimaryVirtualHost {
		slog.Debug(
			"SkippingPrimaryVirtualHostSynchronizer",
			slog.String("reason", "PrimaryVirtualHostEnvValueMatchesCurrent"),
		)
		return
	}

	sync.newPrimaryHostname = primaryVirtualHostEnvValue
	sync.previousPrimaryHostname = webServerPrimaryVirtualHost
	slog.Debug(
		"UpdatingPrimaryVirtualHost",
		slog.String("primaryVirtualHostEnvValue", primaryVirtualHostEnvValue.String()),
		slog.String("webServerPrimaryVirtualHost", webServerPrimaryVirtualHost.String()),
	)

	// phpConfUpdater must be first, as it relies on the previous primary hostname being set
	// on the web server config before other steps can be performed.
	phpConfErr := sync.phpConfUpdater()
	if phpConfErr != nil {
		slog.Error(
			"PrimaryVirtualHostPhpConfUpdaterFailed",
			slog.String("err", phpConfErr.Error()),
		)
		return
	}

	confErr := sync.confUpdater()
	if confErr != nil {
		slog.Error(
			"PrimaryVirtualHostConfUpdaterFailed", slog.String("err", confErr.Error()),
		)
		return
	}

	dbErr := sync.dbUpdater()
	if dbErr != nil {
		slog.Error(
			"PrimaryVirtualHostDbUpdaterFailed", slog.String("err", dbErr.Error()),
		)
		return
	}
}
