package repository

import "github.com/speedianet/sam/src/domain/valueObject"

type WsQueryRepo interface {
	GetVirtualHosts() ([]valueObject.Fqdn, error)
}
