package internalSetupInfra

import (
	"errors"
	"log/slog"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type PrimaryVirtualHostSynchronizer struct {
	persistentDbSvc    *internalDbInfra.PersistentDatabaseService
	newPrimaryHostname tkValueObject.Fqdn
	vhostHelpers       *vhostInfra.VirtualHostHelpers
}

func NewPrimaryVirtualHostSynchronizer(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *PrimaryVirtualHostSynchronizer {
	return &PrimaryVirtualHostSynchronizer{
		persistentDbSvc: persistentDbSvc,
		vhostHelpers:    vhostInfra.NewVirtualHostHelpers(),
	}
}

func (sync *PrimaryVirtualHostSynchronizer) confUpdater() error {
	updateErr := sync.vhostHelpers.UpdateWebServerPrimaryVirtualHost(
		sync.newPrimaryHostname,
	)
	if updateErr != nil {
		return errors.New("ConfUpdateFailed: " + updateErr.Error())
	}

	return nil
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
			"InvalidPrimaryVirtualHostEnvValue",
			slog.String("error", parseErr.Error()),
		)
		os.Exit(1)
	}

	webServerPrimaryVirtualHost, err := sync.vhostHelpers.
		ReadPrimaryVirtualHostHostnameFromWebServerConf()
	if err != nil {
		slog.Error(
			"ReadPrimaryVirtualHostHostnameFailed",
			slog.String("error", err.Error()),
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
	slog.Debug(
		"UpdatingPrimaryVirtualHost",
		slog.String("primaryVirtualHostEnvValue", primaryVirtualHostEnvValue.String()),
		slog.String("webServerPrimaryVirtualHost", webServerPrimaryVirtualHost.String()),
	)

	confErr := sync.confUpdater()
	if confErr != nil {
		slog.Error(
			"UpdatePrimaryVirtualHostConfFailed",
			slog.String("error", confErr.Error()),
		)
		os.Exit(1)
	}

	dbErr := sync.dbUpdater()
	if dbErr != nil {
		slog.Error(
			"UpdatePrimaryVirtualHostDbFailed",
			slog.String("error", dbErr.Error()),
		)
		os.Exit(1)
	}
}
