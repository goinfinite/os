package useCaseHelper

import (
	"log/slog"
	"slices"

	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func IsPhpVersionInstalled(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	phpVersion valueObject.PhpVersion,
) bool {
	phpVersions, err := runtimeQueryRepo.ReadPhpVersionsInstalled()
	if err != nil {
		slog.Error("ReadPhpVersionsInstalledError", slog.String("err", err.Error()))
		return false
	}

	return slices.Contains(phpVersions, phpVersion)
}
