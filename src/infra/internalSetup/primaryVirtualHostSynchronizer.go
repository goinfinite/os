package internalSetupInfra

import (
	"errors"
	"log/slog"
	"os"

	"github.com/goinfinite/os/src/domain/dto"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

type PrimaryVirtualHostSynchronizer struct {
	persistentDbSvc         *internalDbInfra.PersistentDatabaseService
	previousPrimaryHostname tkValueObject.Fqdn
	newPrimaryHostname      tkValueObject.Fqdn
	vhostCmdRepo            *vhostInfra.VirtualHostCmdRepo
	vhostHelpers            *vhostInfra.VirtualHostHelpers
}

func NewPrimaryVirtualHostSynchronizer(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *PrimaryVirtualHostSynchronizer {
	return &PrimaryVirtualHostSynchronizer{
		persistentDbSvc: persistentDbSvc,
		vhostCmdRepo:    vhostInfra.NewVirtualHostCmdRepo(persistentDbSvc),
		vhostHelpers:    vhostInfra.NewVirtualHostHelpers(),
	}
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
		os.Exit(1)
	}

	webServerPrimaryVirtualHost, err := sync.vhostHelpers.
		ReadPrimaryVirtualHostHostnameFromWebServerConf()
	if err != nil {
		slog.Error(
			"ReadPrimaryVirtualHostHostnameFailed", slog.String("err", err.Error()),
		)
		os.Exit(1)
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

	confErr := sync.confUpdater()
	if confErr != nil {
		slog.Error(
			"ConfUpdaterFailed",
			slog.String("err", confErr.Error()),
		)
		os.Exit(1)
	}

	dbErr := sync.dbUpdater()
	if dbErr != nil {
		slog.Error(
			"DbUpdaterFailed",
			slog.String("err", dbErr.Error()),
		)
		os.Exit(1)
	}
}
