package dto

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestUpdatePhpConfigsResponseMarshalsTaskAndParsingFailures(
	test *testing.T,
) {
	taskId, err := valueObject.NewScheduledTaskId(42)
	if err != nil {
		test.Fatalf("TaskIdCreationFailed: %v", err)
	}
	reason := valueObject.NewFailureReason("InvalidPhpModuleName")

	response := NewUpdatePhpConfigsResponse(
		&taskId,
		[]PhpModuleParsingFailure{
			{
				Index: 1, Reason: reason,
			},
		},
	)
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		test.Fatalf("UpdatePhpConfigsResponseMarshalFailed: %v", err)
	}

	encodedResponseString := string(encodedResponse)
	if encodedResponseString == "" {
		test.Fatal("UpdatePhpConfigsResponseShouldNotBeEmpty")
	}
	if !strings.Contains(encodedResponseString, `"taskId":42`) {
		test.Error("MissingTaskIdField")
	}
	if !strings.Contains(
		encodedResponseString, `"failedModulesWithParsingErrors"`,
	) {
		test.Error("MissingParsingFailuresField")
	}
}
