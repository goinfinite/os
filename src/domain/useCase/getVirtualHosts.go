package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/repository"
)

func GetVirtualHosts(
	vhostQueryRepo repository.VirtualHostQueryRepo,
) ([]entity.VirtualHost, error) {
	vhosts, err := vhostQueryRepo.Get()
	if err != nil {
		log.Printf("GetVirtualHostsError: %s", err.Error())
		return vhosts, errors.New("GetVirtualHostsInfraError")
	}

	return vhosts, nil
}
