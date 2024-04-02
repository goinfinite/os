package valueObject

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestServiceName(t *testing.T) {
	t.Run("ValidServiceNames", func(t *testing.T) {
		validNamesAndAliases := []string{
			"openlitespeed",
			"litespeed",
			"nginx",
			"node",
			"mysql",
			"nodejs",
			"redis-server",
		}
		for _, name := range validNamesAndAliases {
			_, err := NewServiceName(name)
			if err != nil {
				t.Errorf("Expected no error for %s, got %v", name, err)
			}
		}
	})

	t.Run("InvalidServiceNames", func(t *testing.T) {
		invalidNamesAndAliases := []string{
			"nginx@",
			"my<>sql",
			"php#fpm",
			"node(js)",
		}
		for _, name := range invalidNamesAndAliases {
			_, err := NewServiceName(name)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", name)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest ServiceName
		}

		dataToTest := "mariadb"
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
			DataToTest ServiceName
		}

		dataToTest := "nginx@"
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
