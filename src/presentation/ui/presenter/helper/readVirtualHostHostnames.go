package presenterHelper

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
)

func ReadVirtualHostHostnames(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) ([]string, error) {
	vhostHostnames := []string{}
	virtualHostService := service.NewVirtualHostService(persistentDbSvc, trailDbSvc)

	vhostResponseServiceOutput := virtualHostService.Read(map[string]interface{}{
		"itemsPerPage": 1000,
		"withMappings": false,
	})
	if vhostResponseServiceOutput.Status != service.Success {
		return vhostHostnames, errors.New("ReadVirtualHostServiceBadResponse")
	}

	vhostReadResponse, assertOk := vhostResponseServiceOutput.Body.(dto.ReadVirtualHostsResponse)
	if !assertOk {
		return vhostHostnames, errors.New("AssertReadVirtualHostsResponseFailed")
	}

	for _, vhost := range vhostReadResponse.VirtualHosts {
		vhostHostnames = append(vhostHostnames, vhost.Hostname.String())
	}

	return vhostHostnames, nil
}
