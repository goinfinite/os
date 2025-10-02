package internalDbInfra

import (
	"errors"

	"github.com/glebarez/sqlite"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	"gorm.io/gorm"
)

type TrailDatabaseService struct {
	Handler *gorm.DB
}

func NewTrailDatabaseService() (*TrailDatabaseService, error) {
	ormSvc, err := gorm.Open(
		sqlite.Open("file:"+
			infraEnvs.TrailDatabaseFilePath+
			infraEnvs.PersistentDatabaseConnectionParams,
		),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	dbSvc := &TrailDatabaseService{Handler: ormSvc}
	err = dbSvc.dbMigrate()
	if err != nil {
		return nil, err
	}

	return dbSvc, nil
}

func (service *TrailDatabaseService) dbMigrate() error {
	err := service.Handler.AutoMigrate(
		&dbModel.ActivityRecord{},
		&dbModel.ActivityRecordAffectedResource{},
	)
	if err != nil {
		return errors.New("TrailDatabaseMigrationError: " + err.Error())
	}

	return nil
}
