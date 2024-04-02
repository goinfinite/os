package valueObject

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewMappingPath(t *testing.T) {
	t.Run("ValidMappingPath", func(t *testing.T) {
		validMappingPaths := []string{
			"/",
			"/img/",
			"/index.html",
			".(png|gif|ico|jpg|jpeg)",
			"/(media|images|cache|tmp|logs)/.*.(php|jsp|pl|py|asp|cgi|sh)$",
			"something",
			"@opencart",
			"/(uploads|files|wp-content|wp-includes|akismet)/.*.php",
			"\\.php(/|$)",
		}

		for _, path := range validMappingPaths {
			_, err := NewMappingPath(path)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", path, err.Error())
			}
		}
	})

	t.Run("InvalidMappingPath", func(t *testing.T) {
		invalidMappingPaths := []string{
			"",
			"UNION SELECT * FROM USERS",
			"/path\n/path",
			"?param=value",
			"https://www.google.com",
			"/path/'; DROP TABLE users; --",
		}

		for _, path := range invalidMappingPaths {
			_, err := NewMappingPath(path)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", path)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest MappingPath
		}

		dataToTest := "/index.html"
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
			DataToTest MappingPath
		}

		dataToTest := "?param=value"
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
