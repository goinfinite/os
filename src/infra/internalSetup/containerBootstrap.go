package internalSetupInfra

import (
	"log/slog"
	"os"
	"slices"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkInfra "github.com/goinfinite/tk/src/infra"
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
)

type ContainerBootstrap struct {
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	transientDbSvc   *tkInfraDb.TransientDatabaseService
	webServerSetup   *WebServerSetup
	fileClerk        tkInfra.FileClerk
	foundationalDirs []string
}

func NewContainerBootstrap(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *tkInfraDb.TransientDatabaseService,
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
			infraEnvs.ToolchainDataDir,
		},
	}
}

func (cb *ContainerBootstrap) isFirstBoot() bool {
	return !slices.ContainsFunc(cb.foundationalDirs, cb.fileClerk.FileExists)
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
	}

	for _, webServerOwnedDir := range []string{
		infraEnvs.PrimaryVirtualHostPublicDir,
		infraEnvs.WebServerLogDir,
	} {
		chownErr := infraHelper.PathOwnership{}.UpdateForWebServerUse(
			webServerOwnedDir, true, false,
		)
		if chownErr != nil {
			slog.Error(
				"ChownWebServerOwnedDirFailed",
				slog.String("path", webServerOwnedDir),
				slog.String("err", chownErr.Error()),
			)
			os.Exit(1)
		}
	}
}

func (cb *ContainerBootstrap) hardenFileOwnership() {
	appRootDirPermissions := os.ModeSticky | os.FileMode(0o775)
	rootOwnedTargets := []struct {
		path        string
		groupName   string
		isRecursive bool
		permissions *os.FileMode
	}{
		{
			infraEnvs.ApplicationRootDir, infraEnvs.PhpWebServerGroupName,
			false, &appRootDirPermissions,
		},
		{infraEnvs.ApplicationConfDir, "root", false, nil},
		{infraEnvs.VirtualHostsConfDir, "root", true, nil},
		{infraEnvs.PkiConfDir, "root", true, nil},
	}

	for _, target := range rootOwnedTargets {
		if !cb.fileClerk.FileExists(target.path) {
			continue
		}

		chownErr := infraHelper.PathOwnership{}.UpdateForRootUse(
			target.path, target.groupName, target.isRecursive,
		)
		if chownErr != nil {
			slog.Error(
				"ChownRootOwnedPathFailed",
				slog.String("path", target.path),
				slog.String("err", chownErr.Error()),
			)
			os.Exit(1)
		}

		if target.permissions == nil {
			continue
		}

		chmodErr := os.Chmod(target.path, *target.permissions)
		if chmodErr != nil {
			slog.Error(
				"ChmodRootOwnedPathFailed",
				slog.String("path", target.path),
				slog.String("err", chmodErr.Error()),
			)
			os.Exit(1)
		}
	}

	webServerOwnedTargets := []string{
		infraEnvs.ToolchainDataDir,
	}
	for _, target := range webServerOwnedTargets {
		if !cb.fileClerk.FileExists(target) {
			continue
		}

		chownErr := infraHelper.PathOwnership{}.UpdateForWebServerUse(
			target, true, false,
		)
		if chownErr != nil {
			slog.Error(
				"ChownWebServerOwnedPathFailed",
				slog.String("path", target),
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
	cb.hardenFileOwnership()
	cb.webServerSetup.onStartSetupOrchestrator()
	NewPrimaryVirtualHostSynchronizer(cb.persistentDbSvc).Run()
}
