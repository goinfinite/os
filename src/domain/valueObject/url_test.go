package valueObject

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNewUrl(t *testing.T) {
	t.Run("ValidUrl", func(t *testing.T) {
		validUrls := []string{
			// cSpell:disable
			"localhost",
			"localhost:8080",
			"speedia.net",
			"http://speedia.net/",
			"http://www.speedia.net",
			"https://speedia.net/",
			"https://www.speedia.net/",
			"http://localhost:8080/v1/ticket/253/attachment/b8680d5cc332672c649f4ff8d9e3b77f.svg",
			"https://www.cnn.com/politics/live-news/house-speaker-vote-10-20-23/index.html",
			"https://blog.goinfinite.net/2023/10/vulnerabilidades-top-10-da-owasp-parte-1/",
			"https://upload.wikimedia.org/wikipedia/commons/thumb/9/98/WordPress_blue_logo.svg/1200px-WordPress_blue_logo.svg.png",
			// cSpell:enable
		}

		for _, url := range validUrls {
			_, err := NewUrl(url)
			if err != nil {
				t.Errorf("Expected no error for '%s', got '%s'", url, err.Error())
			}
		}
	})

	t.Run("InvalidUrl", func(t *testing.T) {
		invalidUrls := []string{
			// cSpell:disable
			"",
			" ",
			"http://",
			"https://",
			"http://notãvalidurl.com/",
			"https://invalidmaçalink.com.br/",
			":8080:/",
			"www.GoOgle.com/",
			"/home/downloads/",
			"DROP TABLE users;",
			"SELECT * FROM users;",
			"<script>alert('XSS')</script>",
			"http://<script>alert('XSS')</script>",
			"https://<script>alert('XSS')</script>",
			"rm -rf /",
			"(){|:& };:",
			"INSERT INTO users (name, email) VALUES ('admin', 'admin@example.com');",
			"sudo rm -r /",
			// cSpell:enable
		}

		for _, url := range invalidUrls {
			_, err := NewUrl(url)
			if err == nil {
				t.Errorf("Expected error for '%s', got nil", url)
			}
		}
	})

	t.Run("GetPort", func(t *testing.T) {
		url, _ := NewUrl("localhost:8080")
		port, _ := url.GetPort()
		if port.Get() != 8080 {
			t.Errorf("Expected port '8080', got '%d'", port.Get())
		}
	})

	t.Run("ValidUnmarshalJSON", func(t *testing.T) {
		var testStruct struct {
			DataToTest Url
		}

		dataToTest := "https://speedia.net/"
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
			DataToTest Url
		}

		dataToTest := "https://invalidmaçalink.com.br/"
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
