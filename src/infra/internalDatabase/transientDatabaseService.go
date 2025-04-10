package internalDbInfra

import (
	"errors"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type TransientDatabaseService struct {
	Handler *gorm.DB
}

type KeyValueModel struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

func NewTransientDatabaseService() (*TransientDatabaseService, error) {
	ormSvc, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, errors.New("TransientDatabaseConnectionError")
	}

	err = ormSvc.AutoMigrate(&KeyValueModel{})
	if err != nil {
		return nil, errors.New("TransientDatabaseMigrationError: " + err.Error())
	}

	return &TransientDatabaseService{Handler: ormSvc}, nil
}

func (dbSvc *TransientDatabaseService) Has(key string) bool {
	var count int64
	result := dbSvc.Handler.Model(&KeyValueModel{}).
		Where("key = ?", key).Count(&count)
	if result.Error != nil {
		return false
	}

	return count > 0
}

func (dbSvc *TransientDatabaseService) Read(key string) (string, error) {
	var keyValueModel KeyValueModel
	result := dbSvc.Handler.Model(&KeyValueModel{}).
		Where("key = ?", key).Find(&keyValueModel)
	if result.Error != nil {
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		return "", errors.New("KeyNotFound")
	}

	return keyValueModel.Value, nil
}

func (dbSvc *TransientDatabaseService) Set(key string, value string) error {
	keyValueModel := KeyValueModel{Key: key, Value: value}

	result := dbSvc.Handler.Model(&KeyValueModel{}).
		Where("key = ?", key).FirstOrCreate(&keyValueModel)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		result = dbSvc.Handler.Model(&KeyValueModel{}).
			Where("key = ?", key).Update("value", value)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
