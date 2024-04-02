package valueObject

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMktplaceItemId(t *testing.T) {
	t.Run("ValidMktplaceItemId", func(t *testing.T) {
		validMktplaceItemIds := []interface{}{
			1,
			1000,
			65365,
			"12345",
		}

		for _, mii := range validMktplaceItemIds {
			_, err := NewMktplaceItemId(mii)
			if err != nil {
				t.Errorf("Expected no error for %v, got %s", mii, err.Error())
			}
		}
	})

	t.Run("InvalidMktplaceItemId", func(t *testing.T) {
		invalidMktplaceItemIds := []interface{}{
			-1,
			0,
			1000000000000000000,
			"-455",
		}

		for _, mii := range invalidMktplaceItemIds {
			_, err := NewMktplaceItemId(mii)
			if err == nil {
				t.Errorf("Expected error for %v, got nil", mii)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MktplaceItemId
		}

		dataToTest := 5
		mapToTest := map[string]int{
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
		dataToTestStr := strconv.Itoa(dataToTest)
		if dataToTestFromStructStr != dataToTestStr {
			t.Errorf(
				"VO data '%s' after UnmarshalJSON is not the same as the original data '%s'",
				dataToTestFromStructStr,
				dataToTestStr,
			)
		}
	})

	t.Run("InvalidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MktplaceItemId
		}

		dataToTest := 0
		mapToTest := map[string]int{
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

	t.Run("ValidUnmarshalYAML", func(t *testing.T) {
		var testStruct struct {
			DataToTest MktplaceItemId `yaml:"dataToTest"`
		}

		dataToTest := 5
		mapToTest := map[string]int{
			"dataToTest": dataToTest,
		}
		mapBytesToTest, _ := yaml.Marshal(mapToTest)

		reader := strings.NewReader(string(mapBytesToTest))
		yamlDecoder := yaml.NewDecoder(reader)
		err := yamlDecoder.Decode(&testStruct)
		if err != nil {
			t.Fatalf("Expected no error on UnmarshalYAML valid test, got %s", err.Error())
		}

		dataToTestFromStructStr := testStruct.DataToTest.String()
		dataToTestStr := strconv.Itoa(dataToTest)
		if dataToTestFromStructStr != dataToTestStr {
			t.Errorf(
				"VO data '%s' after UnmarshalYAML is not the same as the original data '%s'",
				dataToTestFromStructStr,
				dataToTestStr,
			)
		}
	})

	t.Run("InvalidUnmarshalYAML", func(t *testing.T) {
		var testStruct struct {
			DataToTest MktplaceItemId `yaml:"dataToTest"`
		}

		dataToTest := 0
		mapToTest := map[string]int{
			"dataToTest": dataToTest,
		}
		mapBytesToTest, _ := yaml.Marshal(mapToTest)

		reader := strings.NewReader(string(mapBytesToTest))
		yamlDecoder := yaml.NewDecoder(reader)
		err := yamlDecoder.Decode(&testStruct)
		if err == nil {
			t.Fatal("Expected error on UnmarshalYAML invalid test, got nil")
		}
	})
}
