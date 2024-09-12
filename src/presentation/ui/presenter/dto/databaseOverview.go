package presenterDto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type DatabaseOverview struct {
	Type        valueObject.DatabaseType
	IsInstalled bool
	Databases   []entity.Database
}

func NewDatabaseOverview(
	dbType valueObject.DatabaseType,
	isInstalled bool,
	databases []entity.Database,
) DatabaseOverview {
	return DatabaseOverview{
		Type:        dbType,
		IsInstalled: isInstalled,
		Databases:   databases,
	}
}
