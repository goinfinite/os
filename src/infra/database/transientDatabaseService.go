package databaseInfra

import (
	"errors"

	"github.com/dgraph-io/badger/v4"
)

type TransientDatabaseService struct {
	Handler *badger.DB
}

func NewTransientDatabaseService() (*TransientDatabaseService, error) {
	dbSvcOptions := badger.DefaultOptions("").
		WithInMemory(true).
		WithNumVersionsToKeep(1).
		WithLoggingLevel(badger.ERROR)

	dbSvc, err := badger.Open(dbSvcOptions)
	if err != nil {
		return nil, errors.New("TransientDatabaseConnectionError")
	}

	return &TransientDatabaseService{
		Handler: dbSvc,
	}, nil
}

func (dbSvc *TransientDatabaseService) Has(key string) bool {
	hasKey := false
	err := dbSvc.Handler.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == nil {
			hasKey = true
		}

		return nil
	})
	if err != nil {
		return false
	}

	return hasKey
}

func (dbSvc *TransientDatabaseService) Get(key string) (string, error) {
	var value string
	err := dbSvc.Handler.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			value = string(val)
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return value, nil
}

func (dbSvc *TransientDatabaseService) Set(key string, value string) error {
	err := dbSvc.Handler.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(key), []byte(value))
		return err
	})
	return err
}
