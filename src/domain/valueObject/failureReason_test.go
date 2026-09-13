package valueObject

import (
	"strings"
	"testing"
)

func TestNewFailureReason(t *testing.T) {
	truncatedExpected := strings.Repeat("x", 2048)
	truncatedInput := truncatedExpected + "overflow"

	testCases := []struct {
		name           string
		inputValue     any
		expectedOutput FailureReason
	}{
		{
			name:           "NonEmptyStringIsKept",
			inputValue:     "InvalidRecordId",
			expectedOutput: "InvalidRecordId",
		},
		{
			name:           "FreeFormSentenceIsKept",
			inputValue:     "This user should not be able to update API required policies",
			expectedOutput: "This user should not be able to update API required policies",
		},
		{
			name:           "EmptyStringBecomesMalformed",
			inputValue:     "",
			expectedOutput: "MalformedFailureReason",
		},
		{
			name:           "NonStringBecomesMalformed",
			inputValue:     []string{"InvalidRecordId"},
			expectedOutput: "MalformedFailureReason",
		},
		{
			name:           "OversizedStringIsTruncated",
			inputValue:     truncatedInput,
			expectedOutput: FailureReason(truncatedExpected),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualOutput := NewFailureReason(testCase.inputValue)
			if actualOutput != testCase.expectedOutput {
				t.Errorf(
					"UnexpectedOutputValue: '%v' vs '%v' [%v]",
					actualOutput, testCase.expectedOutput, testCase.inputValue,
				)
			}
		})
	}
}

func TestFailureReasonStringMethod(t *testing.T) {
	actualOutput := FailureReason("MalformedFailureReason").String()
	if actualOutput != "MalformedFailureReason" {
		t.Errorf(
			"UnexpectedOutputValue: '%v' vs 'MalformedFailureReason'",
			actualOutput,
		)
	}
}
