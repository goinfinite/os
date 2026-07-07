package internalSetupInfra

import (
	"log/slog"
	"os"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

type ContainerBootstrap struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	webServerSetup  *WebServerSetup
	fileClerk       tkInfra.FileClerk
}

func NewContainerBootstrap(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *ContainerBootstrap {
	return &ContainerBootstrap{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		webServerSetup:  NewWebServerSetup(persistentDbSvc, transientDbSvc),
		fileClerk:       tkInfra.FileClerk{},
	}
}

func (cb *ContainerBootstrap) setupPrimaryPublicDir() {
	primaryPublicDir := infraEnvs.PrimaryVirtualHostPublicDir

	if !cb.fileClerk.FileExists(primaryPublicDir) {
		slog.Debug("CreatingPrimaryPublicDir")
		createDirErr := cb.fileClerk.CreateDir(primaryPublicDir)
		if createDirErr != nil {
			slog.Error(
				"CreatePrimaryPublicDirFailed",
				slog.String("err", createDirErr.Error()),
			)
			os.Exit(1)
		}

		permissions := 0755
		chmodErr := cb.fileClerk.UpdateFilePermissions(primaryPublicDir, &permissions)
		if chmodErr != nil {
			slog.Error(
				"UpdatePrimaryPublicDirPermissionsFailed",
				slog.String("err", chmodErr.Error()),
			)
			os.Exit(1)
		}
	}

	chownErr := infraHelper.UpdateOwnershipForWebServerUse(primaryPublicDir, false, false)
	if chownErr != nil {
		slog.Error(
			"ChownPrimaryPublicDirFailed",
			slog.String("err", chownErr.Error()),
		)
		os.Exit(1)
	}
}

func (cb *ContainerBootstrap) isFirstBoot() bool {
	return !cb.fileClerk.FileExists("/etc/nginx/dhparam.pem")
}

func (cb *ContainerBootstrap) FirstBoot() {
	if !cb.isFirstBoot() {
		return
	}

	cb.setupPrimaryPublicDir()
	cb.webServerSetup.FirstSetup()
}

func (cb *ContainerBootstrap) OnStart() {
	cb.webServerSetup.OnStartSetup()
}
