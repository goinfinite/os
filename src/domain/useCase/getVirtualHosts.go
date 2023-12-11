package useCase

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetVirtualHosts(
	vhostQueryRepo repository.VirtualHostQueryRepo,
) ([]entity.VirtualHost, error) {
	return vhostQueryRepo.Get()
}
