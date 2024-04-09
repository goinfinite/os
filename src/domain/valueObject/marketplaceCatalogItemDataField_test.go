package valueObject

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMarketplaceCatalogItemDataField(t *testing.T) {
	t.Run("ValidUnmarshalYAML", func(t *testing.T) {
		var testStruct struct {
			DataToTest MarketplaceCatalogItemDataField `yaml:"dataToTest"`
		}

		dataKeyToTest := "SomeNiceDataFieldKey"
		dataValueToTest := "SomeNiceDataFieldValue"
		mapToTest := map[string]interface{}{
			"dataToTest": map[string]string{
				"key":   dataKeyToTest,
				"value": dataValueToTest,
			},
		}
		mapBytesToTest, _ := yaml.Marshal(mapToTest)

		reader := strings.NewReader(string(mapBytesToTest))
		yamlDecoder := yaml.NewDecoder(reader)
		err := yamlDecoder.Decode(&testStruct)
		if err != nil {
			t.Fatalf("Expected no error on UnmarshalYAML valid test, got %s", err.Error())
		}

		dataKeyToTestFromStructStr := testStruct.DataToTest.Key.String()
		if dataKeyToTestFromStructStr != dataKeyToTest {
			t.Errorf(
				"VO data '%s' for 'Key' after UnmarshalYAML is not the same as the original data '%s'",
				dataKeyToTestFromStructStr,
				dataKeyToTest,
			)
		}

		dataValueToTestFromStructStr := testStruct.DataToTest.Value.String()
		if dataValueToTestFromStructStr != dataValueToTest {
			t.Errorf(
				"VO data '%s' for 'Value' after UnmarshalYAML is not the same as the original data '%s'",
				dataKeyToTestFromStructStr,
				dataKeyToTest,
			)
		}
	})

	t.Run("InvalidUnmarshalYAML", func(t *testing.T) {
		var testStruct struct {
			DataToTest MarketplaceCatalogItemDataField `yaml:"dataToTest"`
		}

		dataKeyToTest := ""
		dataValueToTest := ""
		mapToTest := map[string]interface{}{
			"dataToTest": map[string]string{
				"key":   dataKeyToTest,
				"value": dataValueToTest,
			},
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
