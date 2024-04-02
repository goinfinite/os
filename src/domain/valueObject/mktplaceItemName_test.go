package valueObject

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMktplaceItemName(t *testing.T) {
	t.Run("ValidMktplaceItemName", func(t *testing.T) {
		validMktplaceItemNames := []string{
			"wordpress",
			"WordPress",
			"opencart",
			"OpenCart",
			"Magento",
			"magento",
			"Joomla",
			"joomla",
			"Drupal",
			"drupal",
			"Supabase",
			"supabase",
			"Laravel",
			"laravel",
			"rabbitmq",
			"RabbitMQ",
		}
		for _, min := range validMktplaceItemNames {
			_, err := NewMktplaceItemName(min)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", min, err.Error())
			}
		}
	})

	t.Run("ValidMktplaceItemName", func(t *testing.T) {
		invalidMktplaceItemNames := []string{
			"",
			".",
			"..",
			"/",
			"name with space",
			"A very long name without any reason just for the test",
			"ççççççç",
			"<root>",
		}
		for _, min := range invalidMktplaceItemNames {
			_, err := NewMktplaceItemName(min)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", min)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MktplaceItemName
		}

		dataToTest := "wordpress"
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
			DataToTest MktplaceItemName
		}

		dataToTest := "name with space"
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

	t.Run("ValidUnmarshalYAML", func(t *testing.T) {
		var testStruct struct {
			DataToTest MktplaceItemName `yaml:"dataToTest"`
		}

		dataToTest := "wordpress"
		mapToTest := map[string]string{
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
		if dataToTestFromStructStr != dataToTest {
			t.Errorf(
				"VO data '%s' after UnmarshalYAML is not the same as the original data '%s'",
				dataToTestFromStructStr,
				dataToTest,
			)
		}
	})

	t.Run("InvalidUnmarshalYAML", func(t *testing.T) {
		var testStruct struct {
			DataToTest MktplaceItemName `yaml:"dataToTest"`
		}

		dataToTest := "name with space"
		mapToTest := map[string]string{
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
