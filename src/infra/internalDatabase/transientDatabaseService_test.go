package internalDbInfra

import (
	"testing"
)

func TestTransientDbSvc(t *testing.T) {
	dbSvc, err := NewTransientDatabaseService()
	if err != nil {
		t.Fatalf("TransientDbSvcInitFailed: %v", err)
	}

	testKey := "test_key"
	testValue := "test_value"

	if dbSvc.Has(testKey) {
		t.Errorf("TransientDbSvcHasNonExistentKeyFailed: %v", err)
	}

	err = dbSvc.Set(testKey, testValue)
	if err != nil {
		t.Errorf("TransientDbSvcSetFailed: %v", err)
	}

	if !dbSvc.Has(testKey) {
		t.Errorf("TransientDbSvcHasExistingKeyFailed: %v", err)
	}
	value, err := dbSvc.Read(testKey)
	if err != nil {
		t.Errorf("TransientDbSvcReadFailed: %v", err)
	}
	if value != testValue {
		t.Errorf("TransientDbSvcReadIncorrectValue: got %v, want %v", value, testValue)
	}

	nonExistentKey := "non_existent_key"
	_, err = dbSvc.Read(nonExistentKey)
	if err == nil {
		t.Errorf("TransientDbSvcReadNonExistentKeyFailed: %v", err)
	}

	updatedValue := "updated_value"
	err = dbSvc.Set(testKey, updatedValue)
	if err != nil {
		t.Errorf("TransientDbSvcSetFailedWhenUpdating: %v", err)
	}

	value, err = dbSvc.Read(testKey)
	if err != nil {
		t.Errorf("TransientDbSvcReadFailedAfterUpdate: %v", err)
	}
	if value != updatedValue {
		t.Errorf("TransientDbSvcReadIncorrectValueAfterUpdate: got %v, want %v", value, updatedValue)
	}
}
