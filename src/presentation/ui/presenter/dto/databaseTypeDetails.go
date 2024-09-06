package presenterDto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type DatabaseTypeDetails struct {
	Type        valueObject.DatabaseType
	IsInstalled bool
	Databases   []entity.Database
}

func NewDatabaseTypeDetails(
	dbType valueObject.DatabaseType,
	isInstalled bool,
	databases []entity.Database,
) DatabaseTypeDetails {
	return DatabaseTypeDetails{
		Type:        dbType,
		IsInstalled: isInstalled,
		Databases:   databases,
	}
}
