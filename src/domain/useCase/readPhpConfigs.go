package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func ReadPhpConfigs(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	hostname tkValueObject.Fqdn,
) (phpConfigs entity.PhpConfigs, err error) {
	phpConfigs, err = runtimeQueryRepo.ReadPhpConfigs(hostname)
	if errors.Is(err, repository.ErrPhpVirtualHostNotFound) {
		return phpConfigs, err
	}
	if err != nil {
		slog.Error("ReadPhpConfigsError", slog.String("err", err.Error()))
		return phpConfigs, errors.New("ReadPhpConfigsFailed")
	}

	return phpConfigs, nil
}
