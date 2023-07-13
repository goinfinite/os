package repository

import (
	"github.com/speedianet/sam/src/domain/valueObject"
)

type ServicesCmdRepo interface {
	Start(name valueObject.ServiceName) error
	Stop(name valueObject.ServiceName) error
	Install(name valueObject.ServiceName, version *valueObject.ServiceVersion) error
	Uninstall(name valueObject.ServiceName) error
}
