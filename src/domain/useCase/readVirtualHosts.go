package useCase

import (
	"errors"
	"log/slog"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func ReadVirtualHosts(
	vhostQueryRepo repository.VirtualHostQueryRepo,
) ([]entity.VirtualHost, error) {
	vhosts, err := vhostQueryRepo.Read()
	if err != nil {
		slog.Error("ReadVirtualHostsError", slog.Any("err", err))
		return vhosts, errors.New("ReadVirtualHostsInfraError")
	}

	return vhosts, nil
}
