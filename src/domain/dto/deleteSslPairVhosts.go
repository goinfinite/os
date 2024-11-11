package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteSslPairVhosts struct {
	SslPairId             valueObject.SslPairId `json:"sslPairId"`
	VirtualHostsHostnames []valueObject.Fqdn    `json:"virtualHostsHostnames"`
}

func NewDeleteSslPairVhosts(
	sslPairId valueObject.SslPairId,
	virtualHostsHostnames []valueObject.Fqdn,
) DeleteSslPairVhosts {
	return DeleteSslPairVhosts{
		SslPairId:             sslPairId,
		VirtualHostsHostnames: virtualHostsHostnames,
	}
}
