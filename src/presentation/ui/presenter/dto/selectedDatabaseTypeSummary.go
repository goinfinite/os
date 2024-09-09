package presenterDto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type SelectedDatabaseTypeSummary struct {
	Type        valueObject.DatabaseType
	IsInstalled bool
	Databases   []entity.Database
}

func NewSelectedDatabaseTypeSummary(
	dbType valueObject.DatabaseType,
	isInstalled bool,
	databases []entity.Database,
) SelectedDatabaseTypeSummary {
	return SelectedDatabaseTypeSummary{
		Type:        dbType,
		IsInstalled: isInstalled,
		Databases:   databases,
	}
}
