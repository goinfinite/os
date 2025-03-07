package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
)

func ReadVirtualHosts(
	vhostQueryRepo repository.VirtualHostQueryRepo,
) ([]entity.VirtualHost, error) {
	vhosts, err := vhostQueryRepo.Read()
	if err != nil {
		slog.Error("ReadVirtualHostsError", slog.String("err", err.Error()))
		return vhosts, errors.New("ReadVirtualHostsInfraError")
	}

	return vhosts, nil
}
