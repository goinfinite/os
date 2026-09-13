package valueObject

import (
	"strings"
	"testing"
)

func TestNewPhpModuleName(t *testing.T) {
	t.Run("StringInput", func(t *testing.T) {
		maximumLengthName := strings.Repeat("a", 64)

		testCases := []struct {
			inputValue     any
			expectedOutput PhpModuleName
			expectError    bool
		}{
			{"ioncube", "ioncube", false},
			{"apcu", "apcu", false},
			{"pdo_sqlite", "pdo_sqlite", false},
			{"PHP_MODULE", "php_module", false},
			{maximumLengthName, PhpModuleName(maximumLengthName), false},
			{"", "", true},
			{"1module", "", true},
			{"ioncube_loader.so", "", true},
			{"posix/process", "", true},
			{"<script>alert('xss')</script>", "", true},
			{strings.Repeat("a", 65), "", true},
			{[]string{"pcntl"}, "", true},
		}

		for _, testCase := range testCases {
			actualOutput, conversionErr := NewPhpModuleName(testCase.inputValue)
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
			inputValue     PhpModuleName
			expectedOutput string
		}{
			{"ioncube", "ioncube"},
			{"pdo_sqlite", "pdo_sqlite"},
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
