package valueObject

import (
	"encoding/json"
	"strings"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestDataFieldKey(t *testing.T) {
	t.Run("ValidDataFieldKey", func(t *testing.T) {
		validDataFieldKeys := []string{
			"username",
			"user-email",
			"Service-Name_With_Port80",
		}

		for _, dfk := range validDataFieldKeys {
			_, err := NewDataFieldKey(dfk)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dfk, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldKey", func(t *testing.T) {
		invalidLength := 40
		invalidDataFieldKeys := []string{
			"",
			"./test",
			testHelpers.GenerateString(invalidLength),
		}

		for _, dfk := range invalidDataFieldKeys {
			_, err := NewDataFieldKey(dfk)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dfk)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest DataFieldKey
		}

		dataToTest := "SomeNiceDataFieldKey"
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
			DataToTest DataFieldKey
		}

		dataToTest := ""
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
