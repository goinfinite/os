package uiHelper

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
)

func GetVhostHostnames(inputToGetVhostHostnames interface{}) []string {
	vhostHostnames := []string{}

	switch typedInput := inputToGetVhostHostnames.(type) {
	case []dto.VirtualHostWithMappings:
		for _, vhostWithMappings := range typedInput {
			vhostHostnames = append(vhostHostnames, vhostWithMappings.Hostname.String())
		}
	case []entity.SslPair:
		for _, sslPair := range typedInput {
			for _, vhostHostname := range sslPair.VirtualHostsHostnames {
				vhostHostnames = append(vhostHostnames, vhostHostname.String())
			}
		}
	}

	return vhostHostnames
}
