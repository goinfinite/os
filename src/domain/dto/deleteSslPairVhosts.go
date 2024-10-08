package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteSslPairVhosts struct {
	SslPairId             valueObject.SslId  `json:"sslPairId"`
	VirtualHostsHostnames []valueObject.Fqdn `json:"virtualHostsHostnames"`
}

func NewDeleteSslPairVhosts(
	sslPairId valueObject.SslId,
	virtualHostsHostnames []valueObject.Fqdn,
) DeleteSslPairVhosts {
	return DeleteSslPairVhosts{
		SslPairId:             sslPairId,
		VirtualHostsHostnames: virtualHostsHostnames,
	}
}
