package repository

import "github.com/speedianet/os/src/domain/valueObject"

type WsQueryRepo interface {
	GetVirtualHosts() ([]valueObject.Fqdn, error)
}
