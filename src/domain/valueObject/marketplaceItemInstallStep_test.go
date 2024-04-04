package valueObject

import (
	"encoding/json"
	"strings"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"gopkg.in/yaml.v3"
)

func TestMarketplaceItemInstallStep(t *testing.T) {
	t.Run("ValidMarketplaceItemInstallStep", func(t *testing.T) {
		validMarketplaceItemInstallSteps := []string{
			"ls -l",
			"cat file.txt | grep \"pattern\" | sort",
			"echo \"Today is $(date +%A)\"",
			"mkdir test_directory && cd test_directory && touch file1.txt file2.txt && ls",
			"certbot certonly --webroot --webroot-path /app/html --agree-tos --register-unsafely-without-email --cert-name speedia.net -d speedia.net",
			"wget https://github.com/speedianet/os -O $PATH",
		}

		for _, miis := range validMarketplaceItemInstallSteps {
			_, err := NewMarketplaceItemInstallStep(miis)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", miis, err.Error())
			}
		}
	})

	t.Run("ValidMarketplaceItemInstallStep", func(t *testing.T) {
		invalidLength := 700
		invalidMarketplaceItemInstallSteps := []string{
			"",
			testHelpers.GenerateString(invalidLength),
		}

		for _, miis := range invalidMarketplaceItemInstallSteps {
			_, err := NewMarketplaceItemInstallStep(miis)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", miis)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MarketplaceItemInstallStep
		}

		dataToTest := "cat file.txt | grep \"pattern\" | sort"
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
			DataToTest MarketplaceItemInstallStep
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
			DataToTest MarketplaceItemInstallStep `yaml:"dataToTest"`
		}

		dataToTest := "cat file.txt | grep \"pattern\" | sort"
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
			DataToTest MarketplaceItemInstallStep `yaml:"dataToTest"`
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
