package useCase

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
)

func GetVirtualHostsWithMappings(
	vhostQueryRepo repository.VirtualHostQueryRepo,
) ([]dto.VirtualHostWithMappings, error) {
	return vhostQueryRepo.GetWithMappings()
}
