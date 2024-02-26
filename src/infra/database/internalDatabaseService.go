package databaseInfra

import (
	"errors"
	"reflect"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

const DatabaseFilePath = "/speedia/sos.db"

type InternalDatabaseService struct {
	Handler *gorm.DB
}

func NewInternalDatabaseService() (*InternalDatabaseService, error) {
	ormSvc, err := gorm.Open(
		sqlite.Open(DatabaseFilePath),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	internalDbSvc := &InternalDatabaseService{Handler: ormSvc}
	err = internalDbSvc.dbMigrate()
	if err != nil {
		return nil, err
	}

	return internalDbSvc, nil
}

func (internalDbSvc InternalDatabaseService) isTableEmpty(model interface{}) (bool, error) {
	var count int64
	err := internalDbSvc.Handler.Model(&model).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (internalDbSvc InternalDatabaseService) seedDatabase(seedModels ...interface{}) error {
	for _, seedModel := range seedModels {
		isTableEmpty, err := internalDbSvc.isTableEmpty(seedModel)
		if err != nil {
			return err
		}

		if !isTableEmpty {
			continue
		}

		seedModelType := reflect.TypeOf(seedModel).Elem()
		seedModelFieldsAndMethods := reflect.ValueOf(seedModel)

		seedModelInitialEntriesMethod := seedModelFieldsAndMethods.MethodByName(
			"InitialEntries",
		)
		seedModelInitialEntriesMethodResults := seedModelInitialEntriesMethod.Call(
			[]reflect.Value{},
		)
		initialEntries := seedModelInitialEntriesMethodResults[0].Interface()

		for _, entry := range initialEntries.([]interface{}) {
			entryInnerStructure := reflect.ValueOf(entry)

			entryFormatHandlerWillAccept := reflect.New(seedModelType)
			entryFormatHandlerWillAccept.Elem().Set(entryInnerStructure)
			adjustedEntry := entryFormatHandlerWillAccept.Interface()

			err = internalDbSvc.Handler.Create(adjustedEntry).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (internalDbSvc InternalDatabaseService) dbMigrate() error {
	err := internalDbSvc.Handler.AutoMigrate()
	if err != nil {
		return errors.New("DatabaseMigrationError")
	}

	modelsWithInitialEntries := []interface{}{}

	err = internalDbSvc.seedDatabase(modelsWithInitialEntries...)
	if err != nil {
		return errors.New("AddDefaultDatabaseEntriesError")
	}

	return nil
}
