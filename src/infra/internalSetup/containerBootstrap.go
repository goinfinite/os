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
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	transientDbSvc   *internalDbInfra.TransientDatabaseService
	webServerSetup   *WebServerSetup
	fileClerk        tkInfra.FileClerk
	foundationalDirs []string
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
		foundationalDirs: []string{
			infraEnvs.PrimaryVirtualHostPublicDir,
			infraEnvs.CronLogDir,
			infraEnvs.WebServerLogDir,
			infraEnvs.TrashDir,
		},
	}
}

func (cb *ContainerBootstrap) isFirstBoot() bool {
	for _, dirPath := range cb.foundationalDirs {
		if cb.fileClerk.FileExists(dirPath) {
			return false
		}
	}
	return true
}

func (cb *ContainerBootstrap) foundationalDirsCreator() {
	for _, dirPath := range cb.foundationalDirs {
		if !cb.fileClerk.FileExists(dirPath) {
			createErr := cb.fileClerk.CreateDir(dirPath)
			if createErr != nil {
				slog.Error(
					"CreateFoundationalDirFailed",
					slog.String("path", dirPath),
					slog.String("err", createErr.Error()),
				)
				os.Exit(1)
			}
		}

		chownErr := infraHelper.UpdateOwnershipForWebServerUse(dirPath, false, false)
		if chownErr != nil {
			slog.Error(
				"ChownFoundationalDirFailed",
				slog.String("path", dirPath),
				slog.String("err", chownErr.Error()),
			)
			os.Exit(1)
		}
	}
}

func (cb *ContainerBootstrap) FirstBoot() {
	if !cb.isFirstBoot() {
		return
	}

	cb.foundationalDirsCreator()
	cb.webServerSetup.firstSetupOrchestrator()
}

func (cb *ContainerBootstrap) OnStart() {
	cb.webServerSetup.onStartSetupOrchestrator()
	NewPrimaryVirtualHostSynchronizer(cb.persistentDbSvc).Run()
}
