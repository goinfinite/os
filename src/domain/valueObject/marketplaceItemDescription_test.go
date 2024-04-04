package valueObject

import (
	"encoding/json"
	"strings"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"gopkg.in/yaml.v3"
)

func TestMarketplaceItemDescription(t *testing.T) {
	t.Run("ValidMarketplaceItemDescription", func(t *testing.T) {
		validMarketplaceItemDescriptions := []string{
			"Build and grow your website with the best way to WordPress. Lightning-fast hosting, intuitive, flexible editing, and everything you need to grow your site and audience, baked right in.",
			"It's comprised of Elasticsearch, Kibana, Beats, and Logstash (also known as the ELK Stack) and more. Reliably and securely take data from any source, in any format, then search, analyze, and visualize.",
			"RabbitMQ is a reliable and mature messaging and streaming broker, which is easy to deploy on cloud environments, on-premises, and on your local machine. It is currently used by millions worldwide.",
		}
		for _, mid := range validMarketplaceItemDescriptions {
			_, err := NewMarketplaceItemDescription(mid)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mid, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemDescription", func(t *testing.T) {
		invalidLength := 600
		invalidMarketplaceItemDescriptions := []string{
			"",
			"a",
			testHelpers.GenerateString(invalidLength),
		}
		for _, mid := range invalidMarketplaceItemDescriptions {
			_, err := NewMarketplaceItemDescription(mid)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mid)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MarketplaceItemDescription
		}

		dataToTest := "Some nice marketplace item description."
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
			DataToTest MarketplaceItemDescription
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

	t.Run("ValidUnmarshalYAML", func(t *testing.T) {
		var testStruct struct {
			DataToTest MarketplaceItemDescription `yaml:"dataToTest"`
		}

		dataToTest := "Some nice marketplace item description."
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
			DataToTest MarketplaceItemDescription `yaml:"dataToTest"`
		}

		dataToTest := ""
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
