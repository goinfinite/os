package dto

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestUpdatePhpModulesResponseMarshals(t *testing.T) {
	moduleName, err := valueObject.NewPhpModuleName("curl")
	if err != nil {
		t.Fatalf("PhpModuleNameCreationFailed: %v", err)
	}
	reason, err := valueObject.NewFailureReason("InstallPhpModulePackageFailed")
	if err != nil {
		t.Fatalf("FailureReasonCreationFailed: %v", err)
	}

	response := NewUpdatePhpModulesResponse(
		[]PhpModuleUpdate{NewPhpModuleUpdate(moduleName, true)},
		[]PhpModuleUpdateFailure{NewPhpModuleUpdateFailure(moduleName, false, reason)},
	)
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("UpdatePhpModulesResponseMarshalFailed: %v", err)
	}

	encodedResponseString := string(encodedResponse)
	if encodedResponseString == "" {
		t.Fatal("UpdatePhpModulesResponseShouldNotBeEmpty")
	}
	if !strings.Contains(encodedResponseString, `"modulesSuccessfullyUpdated"`) {
		t.Error("MissingSuccessfulModulesField")
	}
	if !strings.Contains(encodedResponseString, `"failedModulesWithReason"`) {
		t.Error("MissingFailedModulesField")
	}
}
