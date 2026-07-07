package cliMiddleware

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
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	newPrimaryHostname tkValueObject.Fqdn
	vhostHelpers     *vhostInfra.VirtualHostHelpers
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
	var primaryVhost dbModel.VirtualHost
	queryErr := sync.persistentDbSvc.Handler.Where(
		"is_primary = ?", true,
	).First(&primaryVhost).Error
	if queryErr != nil {
		return errors.New("DbReadFailed: " + queryErr.Error())
	}

	if primaryVhost.Hostname == sync.newPrimaryHostname.String() {
		return nil
	}

	slog.Info(
		"PrimaryVirtualHostSynchronizer",
		slog.String("phase", "db-update"),
		slog.String("oldHostname", primaryVhost.Hostname),
		slog.String("newHostname", sync.newPrimaryHostname.String()),
	)

	updateErr := sync.persistentDbSvc.Handler.Model(
		&primaryVhost,
	).Update("hostname", sync.newPrimaryHostname.String()).Error
	if updateErr != nil {
		return errors.New("DbWriteFailed: " + updateErr.Error())
	}

	return nil
}

func (sync *PrimaryVirtualHostSynchronizer) Run() {
	rawEnvHostnameStr := os.Getenv(infraEnvs.PrimaryVirtualHostEnvKey)
	if rawEnvHostnameStr == "" {
		slog.Info(
			"PrimaryVirtualHostSynchronizer",
			slog.String("phase", "skip"),
			slog.String("reason", "EmptyEnvVar"),
		)
		return
	}

	newPrimaryHostname, parseErr := tkValueObject.NewFqdn(rawEnvHostnameStr)
	if parseErr != nil {
		slog.Error(
			"PrimaryVirtualHostSynchronizer",
			slog.String("phase", "parse"),
			slog.String("error", parseErr.Error()),
		)
		os.Exit(1)
	}

	sync.newPrimaryHostname = newPrimaryHostname

	confErr := sync.confUpdater()
	if confErr != nil {
		slog.Error(
			"PrimaryVirtualHostSynchronizer",
			slog.String("phase", "conf"),
			slog.String("error", confErr.Error()),
		)
		os.Exit(1)
	}

	dbErr := sync.dbUpdater()
	if dbErr != nil {
		slog.Error(
			"PrimaryVirtualHostSynchronizer",
			slog.String("phase", "db"),
			slog.String("error", dbErr.Error()),
		)
		os.Exit(1)
	}
}
