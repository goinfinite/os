package valueObject

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMktplaceItemType(t *testing.T) {
	t.Run("ValidMktplaceItemType", func(t *testing.T) {
		validMktplaceItemTypes := []string{
			"app",
			"framework",
			"stack",
		}
		for _, mit := range validMktplaceItemTypes {
			_, err := NewMktplaceItemType(mit)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mit, err.Error())
			}
		}
	})

	t.Run("InvalidMktplaceItemType", func(t *testing.T) {
		invalidMktplaceItemTypes := []string{
			"",
			"service",
			"mobile",
			"ml-model",
			"repository",
		}
		for _, mit := range invalidMktplaceItemTypes {
			_, err := NewMktplaceItemType(mit)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mit)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MktplaceItemType
		}

		dataToTest := "framework"
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
			DataToTest MktplaceItemType
		}

		dataToTest := "service"
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
