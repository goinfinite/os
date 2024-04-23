package dto

import "github.com/speedianet/os/src/domain/valueObject"

type DeleteSslPairVhosts struct {
	SslPairId    valueObject.SslId  `json:"sslPairId"`
	VirtualHosts []valueObject.Fqdn `json:"virtualHosts"`
}

func NewDeleteSslPairVhosts(
	sslPairId valueObject.SslId,
	virtualHosts []valueObject.Fqdn,
) DeleteSslPairVhosts {
	return DeleteSslPairVhosts{
		SslPairId:    sslPairId,
		VirtualHosts: virtualHosts,
	}
}
