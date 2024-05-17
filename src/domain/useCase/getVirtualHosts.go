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
	vhosts, err := vhostQueryRepo.Read()
	if err != nil {
		log.Printf("ReadVirtualHostsError: %s", err.Error())
		return vhosts, errors.New("ReadVirtualHostsInfraError")
	}

	return vhosts, nil
}
