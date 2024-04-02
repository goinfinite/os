package valueObject

import (
	"encoding/json"
	"strings"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestMappingMatchPattern(t *testing.T) {
	t.Run("ValidMappingMatchPattern", func(t *testing.T) {
		validMappingMatchPatterns := []string{
			"begins-with",
			"contains",
			"equals",
			"ends-with",
		}

		for _, mmp := range validMappingMatchPatterns {
			_, err := NewMappingMatchPattern(mmp)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mmp, err.Error())
			}
		}
	})

	t.Run("InvalidMappingMatchPattern", func(t *testing.T) {
		invalidLength := 70
		invalidMappingMatchPatterns := []string{
			"",
			"bigger-then",
			"diff",
			"has-prefix",
			testHelpers.GenerateString(invalidLength),
		}

		for _, mmp := range invalidMappingMatchPatterns {
			_, err := NewMappingMatchPattern(mmp)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mmp)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MappingMatchPattern
		}

		dataToTest := "begins-with"
		mapToTest := map[string]string{
			"dataToTest": dataToTest,
		}
		mapBytesToTest, _ := json.Marshal(mapToTest)

		reader := strings.NewReader(string(mapBytesToTest))
		jsonDecoder := json.NewDecoder(reader)
		err := jsonDecoder.Decode(&testStruct)
		if err != nil {
			t.Fatalf("Expected no error on UnmarshalJSON valid test, got %s", err.Error())
		}

		dataToTestFromStructStr := testStruct.DataToTest.String()
		if dataToTestFromStructStr != dataToTest {
			t.Errorf(
				"VO data '%s' after UnmarshalJSON is not the same as the original data '%s'",
				dataToTestFromStructStr,
				dataToTest,
			)
		}
	})

	t.Run("InvalidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MappingMatchPattern
		}

		dataToTest := "has-prefix"
		mapToTest := map[string]string{
			"dataToTest": dataToTest,
		}
		mapBytesToTest, _ := json.Marshal(mapToTest)

		reader := strings.NewReader(string(mapBytesToTest))
		jsonDecoder := json.NewDecoder(reader)
		err := jsonDecoder.Decode(&testStruct)
		if err == nil {
			t.Fatal("Expected error on UnmarshalJSON invalid test, got nil")
		}
	})
}
