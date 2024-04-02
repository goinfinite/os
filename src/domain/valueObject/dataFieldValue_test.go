package valueObject

import (
	"encoding/json"
	"strings"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestDataFieldValue(t *testing.T) {
	t.Run("ValidDataFieldValue", func(t *testing.T) {
		validDataFieldValues := []string{
			"This is my username",
			"new_email@mail.net",
			"localhost:8000",
			"https://www.google.com/search",
		}

		for _, dfv := range validDataFieldValues {
			_, err := NewDataFieldValue(dfv)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", dfv, err.Error())
			}
		}
	})

	t.Run("InvalidDataFieldValue", func(t *testing.T) {
		invalidLength := 70
		invalidDataFieldValues := []string{
			"",
			"a",
			testHelpers.GenerateString(invalidLength),
		}

		for _, dfv := range invalidDataFieldValues {
			_, err := NewDataFieldValue(dfv)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", dfv)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest DataFieldValue
		}

		dataToTest := "SomeNiceDataFieldValue"
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
			DataToTest DataFieldValue
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
