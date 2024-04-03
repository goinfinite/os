package valueObject

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMappingTargetType(t *testing.T) {
	t.Run("ValidMappingTargetType", func(t *testing.T) {
		validMappingTargetTypes := []string{
			"url",
			"service",
			"response-code",
			"inline-html",
			"static-files",
		}

		for _, mtt := range validMappingTargetTypes {
			_, err := NewMappingTargetType(mtt)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mtt, err.Error())
			}
		}
	})

	t.Run("ValidMappingTargetType", func(t *testing.T) {
		invalidMappingTargetTypes := []string{
			"response-header",
			"reverse-proxy",
			"template",
		}

		for _, mtt := range invalidMappingTargetTypes {
			_, err := NewMappingTargetType(mtt)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mtt)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MappingTargetType
		}

		dataToTest := "service"
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
			DataToTest MappingTargetType
		}

		dataToTest := "reverse-proxy"
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
