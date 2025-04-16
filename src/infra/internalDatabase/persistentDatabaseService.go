package internalDbInfra

import (
	"errors"
	"reflect"

	"github.com/glebarez/sqlite"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	"gorm.io/gorm"
)

type PersistentDatabaseService struct {
	Handler *gorm.DB
}

func NewPersistentDatabaseService() (*PersistentDatabaseService, error) {
	ormSvc, err := gorm.Open(
		sqlite.Open(infraEnvs.PersistentDatabaseFilePath),
		&gorm.Config{},
	)
	if err != nil {
		return nil, errors.New("DatabaseConnectionError")
	}

	dbSvc := &PersistentDatabaseService{Handler: ormSvc}
	err = dbSvc.dbMigrate()
	if err != nil {
		return nil, err
	}

	return dbSvc, nil
}

func (dbSvc *PersistentDatabaseService) isTableEmpty(model interface{}) (bool, error) {
	var count int64
	err := dbSvc.Handler.Model(&model).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (dbSvc *PersistentDatabaseService) seedDatabase(
	seedModels map[string]interface{},
) error {
	for modelName, seedModel := range seedModels {
		isTableEmpty, err := dbSvc.isTableEmpty(seedModel)
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
		initialEntriesMethodResults := seedModelInitialEntriesMethod.Call(
			[]reflect.Value{},
		)

		initialEntriesMethodErr := initialEntriesMethodResults[1]
		if !initialEntriesMethodErr.IsNil() {
			err = initialEntriesMethodErr.Interface().(error)
			if err != nil {
				return errors.New(
					"SeedModelInitialEntriesError (" + modelName + "): " + err.Error(),
				)
			}
		}

		initialEntries := initialEntriesMethodResults[0].Interface()
		for _, entry := range initialEntries.([]interface{}) {
			entryInnerStructure := reflect.ValueOf(entry)

			entryFormatHandlerWillAccept := reflect.New(seedModelType)
			entryFormatHandlerWillAccept.Elem().Set(entryInnerStructure)
			adjustedEntry := entryFormatHandlerWillAccept.Interface()

			err = dbSvc.Handler.Create(adjustedEntry).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (dbSvc *PersistentDatabaseService) dbMigrate() error {
	err := dbSvc.Handler.AutoMigrate(
		&dbModel.Account{},
		&dbModel.InstalledService{},
		&dbModel.Mapping{},
		&dbModel.MappingSecurityRule{},
		&dbModel.MarketplaceInstalledItem{},
		&dbModel.ScheduledTask{},
		&dbModel.ScheduledTaskTag{},
		&dbModel.SecureAccessPublicKey{},
		&dbModel.VirtualHost{},
	)
	if err != nil {
		return errors.New("PersistentDatabaseMigrationError: " + err.Error())
	}

	modelsWithInitialEntries := map[string]interface{}{
		"VirtualHost":         &dbModel.VirtualHost{},
		"InstalledService":    &dbModel.InstalledService{},
		"MappingSecurityRule": &dbModel.MappingSecurityRule{},
	}

	err = dbSvc.seedDatabase(modelsWithInitialEntries)
	if err != nil {
		return errors.New("CreateDefaultDatabaseEntriesError: " + err.Error())
	}

	return nil
}
