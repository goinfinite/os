package valueObject

import (
	"encoding/json"
	"strings"
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestInlineHtmlContent(t *testing.T) {
	t.Run("ValidInlineHtmlContent", func(t *testing.T) {
		validInlineHtmlContents := []string{
			"Some nice inline html content",
			"<h1>Nice title here</h1>",
			"<p>With some regular text here too...<h2>",
		}

		for _, ihc := range validInlineHtmlContents {
			_, err := NewInlineHtmlContent(ihc)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", ihc, err.Error())
			}
		}
	})

	t.Run("ValidInlineHtmlContent", func(t *testing.T) {
		invalidLength := 3600
		invalidInlineHtmlContents := []string{
			"",
			testHelpers.GenerateString(invalidLength),
		}

		for _, ihc := range invalidInlineHtmlContents {
			_, err := NewInlineHtmlContent(ihc)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", ihc)
			}
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest InlineHtmlContent
		}

		dataToTest := "Some nice inline html content"
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
			DataToTest InlineHtmlContent
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
