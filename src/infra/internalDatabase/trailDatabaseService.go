package internalDbInfra

import (
	"errors"
	"fmt"
	"log/slog"

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

func (service *TrailDatabaseService) migrateOperatorAccountIdToSri() {
	if !service.Handler.Migrator().HasColumn(&tkInfraDbModel.ActivityRecord{}, "operator_account_id") {
		return
	}

	err := service.Handler.Exec(
		"UPDATE activity_records SET operator_sri = 'sri://accountId:account/' || operator_account_id WHERE operator_account_id IS NOT NULL AND (operator_sri IS NULL OR operator_sri = '')",
	).Error
	if err != nil {
		slog.Error("MigrateOperatorAccountIdToSriError", slog.String("err", err.Error()))
		return
	}

	err = service.Handler.Migrator().DropColumn(&tkInfraDbModel.ActivityRecord{}, "operator_account_id")
	if err != nil {
		slog.Error(
			"DropOperatorAccountIdColumnError",
			slog.String("err", fmt.Sprintf("%v", err)),
		)
	}
}

func (service *TrailDatabaseService) dbMigrate() error {
	err := service.Handler.AutoMigrate(
		&tkInfraDbModel.ActivityRecord{},
		&tkInfraDbModel.ActivityRecordAffectedResource{},
	)
	if err != nil {
		return errors.New("TrailDatabaseMigrationError: " + err.Error())
	}

	service.migrateOperatorAccountIdToSri()

	return nil
}
