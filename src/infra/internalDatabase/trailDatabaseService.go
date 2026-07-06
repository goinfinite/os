package internalDbInfra

import (
	"errors"
	"time"

	"github.com/glebarez/sqlite"
	tkInfraDbModel "github.com/goinfinite/tk/src/infra/db/model"

	infraEnvs "github.com/goinfinite/os/src/infra/envs"
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
		&gorm.Config{
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		},
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
		&tkInfraDbModel.ActivityRecord{},
		&tkInfraDbModel.ActivityRecordAffectedResource{},
	)
	if err != nil {
		return errors.New("TrailDatabaseMigrationError: " + err.Error())
	}

	return nil
}
