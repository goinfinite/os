package uiPresenterHelper

import (
	"errors"
	tkPresentation "github.com/goinfinite/tk/src/presentation"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
)

func ReadVirtualHostHostnames(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) ([]string, error) {
	vhostHostnames := []string{}
	virtualHostLiaison := liaison.NewVirtualHostLiaison(persistentDbSvc, trailDbSvc)

	vhostResponseLiaisonOutput := virtualHostLiaison.Read(map[string]any{
		"itemsPerPage": 1000,
		"withMappings": false,
	})
	if vhostResponseLiaisonOutput.Status != tkPresentation.LiaisonResponseStatusSuccess {
		return vhostHostnames, errors.New("ReadVirtualHostLiaisonBadResponse")
	}

	vhostReadResponse, assertOk := vhostResponseLiaisonOutput.Body.(dto.ReadVirtualHostsResponse)
	if !assertOk {
		return vhostHostnames, errors.New("AssertReadVirtualHostsResponseFailed")
	}

	primaryHostname := ""
	for _, vhostEntity := range vhostReadResponse.VirtualHosts {
		if vhostEntity.Type == valueObject.VirtualHostTypeAlias {
			continue
		}
		if vhostEntity.IsPrimary {
			primaryHostname = vhostEntity.Hostname.String()
			continue
		}

		vhostHostnames = append(vhostHostnames, vhostEntity.Hostname.String())
	}

	if primaryHostname != "" {
		vhostHostnames = append([]string{primaryHostname}, vhostHostnames...)
	}

	return vhostHostnames, nil
}
