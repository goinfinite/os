package internalDbInfra

import (
	"testing"
)

func TestTransientDbSvc(t *testing.T) {
	dbSvc, err := NewTransientDatabaseService()
	if err != nil {
		t.Fatalf("TransientDatabaseServiceInitFailed: %v", err)
	}

	testKey := "test_key"
	testValue := "test_value"
	nonExistentKey := "non_existent_key"
	updatedValue := "updated_value"

	t.Run("HasFunctionality", func(t *testing.T) {
		if dbSvc.Has(testKey) {
			t.Errorf("HasReturnedTrueForNonExistentKey")
		}

		err := dbSvc.Set(testKey, testValue)
		if err != nil {
			t.Errorf("SetFailed: %v", err)
		}

		if !dbSvc.Has(testKey) {
			t.Errorf("HasReturnedFalseForExistingKey")
		}
	})

	t.Run("ReadWriteFunctionality", func(t *testing.T) {
		value, err := dbSvc.Read(testKey)
		if err != nil {
			t.Errorf("ReadFailed: %v", err)
		}
		if value != testValue {
			t.Errorf("IncorrectValue: Got %v, Want %v", value, testValue)
		}

		_, err = dbSvc.Read(nonExistentKey)
		if err == nil {
			t.Errorf("ExpectedErrorForNonExistentKeyButGotNil")
		}

		err = dbSvc.Set(testKey, updatedValue)
		if err != nil {
			t.Errorf("UpdateFailed: %v", err)
		}

		value, err = dbSvc.Read(testKey)
		if err != nil {
			t.Errorf("ReadFailedAfterUpdate: %v", err)
		}
		if value != updatedValue {
			t.Errorf("IncorrectValueAfterUpdate: Got %v, Want %v", value, updatedValue)
		}
	})
}
