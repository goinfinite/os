package presenterDto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type DatabaseTypeSummary struct {
	Type        valueObject.DatabaseType
	IsInstalled bool
	Databases   []entity.Database
}

func NewDatabaseTypeSummary(
	dbType valueObject.DatabaseType,
	isInstalled bool,
	databases []entity.Database,
) DatabaseTypeSummary {
	return DatabaseTypeSummary{
		Type:        dbType,
		IsInstalled: isInstalled,
		Databases:   databases,
	}
}
