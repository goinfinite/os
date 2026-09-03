package valueObject

import (
	"strings"
	"testing"
)

func TestNewScheduledTaskOutput(t *testing.T) {
	t.Run("StringInput", func(t *testing.T) {
		const expectedMaximumOutputBytes = 65536

		truncatedOutput := strings.Repeat("x", expectedMaximumOutputBytes)
		truncatingInput := truncatedOutput + "x"
		testCases := []struct {
			inputValue     any
			expectedOutput ScheduledTaskOutput
			expectError    bool
		}{
			{"validOutput", "validOutput", false},
			{"", "", false},
			{123, "123", false},
			{true, "true", false},
			{truncatingInput, ScheduledTaskOutput(truncatedOutput), false},
			{[]string{"output"}, "", true},
		}

		for _, testCase := range testCases {
			actualOutput, conversionErr := NewScheduledTaskOutput(
				testCase.inputValue,
			)
			if testCase.expectError && conversionErr == nil {
				t.Errorf("MissingExpectedError: [%v]", testCase.inputValue)
			}
			if !testCase.expectError && conversionErr != nil {
				t.Errorf(
					"UnexpectedError: '%s' [%v]",
					conversionErr.Error(),
					testCase.inputValue,
				)
			}
			if !testCase.expectError && actualOutput != testCase.expectedOutput {
				t.Errorf(
					"UnexpectedOutputValue: '%v' vs '%v' [%v]",
					actualOutput,
					testCase.expectedOutput,
					testCase.inputValue,
				)
			}
		}
	})

	t.Run("StringMethod", func(t *testing.T) {
		testCases := []struct {
			inputValue     ScheduledTaskOutput
			expectedOutput string
		}{
			{"validOutput", "validOutput"},
			{"everythingIsUpToDate", "everythingIsUpToDate"},
		}

		for _, testCase := range testCases {
			actualOutput := testCase.inputValue.String()
			if actualOutput != testCase.expectedOutput {
				t.Errorf(
					"UnexpectedOutputValue: '%v' vs '%v' [%v]",
					actualOutput,
					testCase.expectedOutput,
					testCase.inputValue,
				)
			}
		}
	})
}
