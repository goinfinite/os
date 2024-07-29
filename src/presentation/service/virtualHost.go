package service

import (
	"github.com/speedianet/os/src/domain/useCase"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

type VirtualHostService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewVirtualHostService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostService {
	return &VirtualHostService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *VirtualHostService) Read() ServiceOutput {
	vhostsQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)
	vhostsList, err := useCase.ReadVirtualHosts(vhostsQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, vhostsList)
}
