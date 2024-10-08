package valueObject

import "testing"

func TestNewMappingTargetValue(t *testing.T) {
	t.Run("ValidMappingTargetValueBasedOnType (Url)", func(t *testing.T) {
		urlTargetType, _ := NewMappingTargetType("url")
		validMappingTargetUrlValues := []interface{}{
			"localhost", "localhost:8080", "goinfinite.net", "http://goinfinite.net/",
			"http://www.goinfinite.net", "https://goinfinite.net/",
			"https://www.goinfinite.net/",
			"http://localhost:8080/v1/ticket/253/attachment/b8680d5cc332672c649f4ff8d9e3b77f.svg",
			"https://www.cnn.com/politics/live-news/house-speaker-vote-10-20-23/index.html",
			"https://blog.goinfinite.net/2023/10/vulnerabilidades-top-10-da-owasp-parte-1/",
			"https://upload.wikimedia.org/wikipedia/commons/thumb/9/98/WordPress_blue_logo.svg/1200px-WordPress_blue_logo.svg.png",
		}

		for _, urlValue := range validMappingTargetUrlValues {
			_, err := NewMappingTargetValue(urlValue, urlTargetType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", urlValue, err.Error())
			}
		}
	})

	t.Run("InvalidMappingTargetValueBasedOnType (Url)", func(t *testing.T) {
		urlTargetType, _ := NewMappingTargetType("url")
		invalidMappingTargetUrlValues := []interface{}{
			"", " ", "http://", "https://", "http://notãvalidurl.com/",
			"https://invalidmaçalink.com.br/", ":8080:/", "www.GoOgle.com/",
			"/home/downloads/", "DROP TABLE users;", "SELECT * FROM users;",
			"<script>alert('XSS')</script>", "http://<script>alert('XSS')</script>",
			"https://<script>alert('XSS')</script>", "rm -rf /", "(){|:& };:",
			"INSERT INTO users (name, email) VALUES ('admin', 'admin@example.com');",
		}

		for _, urlValue := range invalidMappingTargetUrlValues {
			_, err := NewMappingTargetValue(urlValue, urlTargetType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", urlValue)
			}
		}
	})

	t.Run("ValidMappingTargetValueBasedOnType (Service)", func(t *testing.T) {
		svcNameTargetType, _ := NewMappingTargetType("service")
		validMappingTargetServiceNameValues := []interface{}{
			"openlitespeed",
			"litespeed",
			"nginx",
			"node",
			"mysql",
			"nodejs",
			"redis-server",
		}

		for _, urlValue := range validMappingTargetServiceNameValues {
			_, err := NewMappingTargetValue(urlValue, svcNameTargetType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", urlValue, err.Error())
			}
		}
	})

	t.Run("InvalidMappingTargetValueBasedOnType (Service)", func(t *testing.T) {
		svcNameTargetType, _ := NewMappingTargetType("service")
		invalidMappingTargetServiceNameValues := []interface{}{
			"nginx@",
			"my<>sql",
			"php#fpm",
			"node(js)",
		}

		for _, urlValue := range invalidMappingTargetServiceNameValues {
			_, err := NewMappingTargetValue(urlValue, svcNameTargetType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", urlValue)
			}
		}
	})

	t.Run("ValidMappingTargetValueBasedOnType (HTTP Response Code)", func(t *testing.T) {
		responseCodeTargetType, _ := NewMappingTargetType("response-code")
		validMappingTargetResponseCodeValues := []interface{}{
			"100",
			"200",
			"300",
			"400",
			"500",
			100,
			200,
			300,
			400,
			500,
		}

		for _, urlValue := range validMappingTargetResponseCodeValues {
			_, err := NewMappingTargetValue(urlValue, responseCodeTargetType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", urlValue, err.Error())
			}
		}
	})

	t.Run("InvalidMappingTargetValueBasedOnType (HTTP Response Code)", func(t *testing.T) {
		responseCodeTargetType, _ := NewMappingTargetType("response-code")
		invalidMappingTargetResponseCodeValues := []interface{}{
			0,
			100000,
			"@blabla",
			"<script>alert('xss')</script>",
			"1000",
			"0",
			"-1",
			"UNION SELECT * FROM USERS",
		}

		for _, urlValue := range invalidMappingTargetResponseCodeValues {
			_, err := NewMappingTargetValue(urlValue, responseCodeTargetType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", urlValue)
			}
		}
	})

	t.Run("ValidMappingTargetValueBasedOnType (Inline HTML Content)", func(t *testing.T) {
		inlineHtmlContentTargetType, _ := NewMappingTargetType("inline-html")
		validMappingTargetInlineHtmlContentValues := []interface{}{
			"Some nice inline html content",
			"<h1>Nice title here</h1>",
			"<p>With some regular text here too...<h2>",
		}

		for _, urlValue := range validMappingTargetInlineHtmlContentValues {
			_, err := NewMappingTargetValue(urlValue, inlineHtmlContentTargetType)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", urlValue, err.Error())
			}
		}
	})

	t.Run("InvalidMappingTargetValueBasedOnType (Inline HTML Content)", func(t *testing.T) {
		inlineHtmlContentTargetType, _ := NewMappingTargetType("inline-html")
		invalidMappingTargetInlineHtmlContentValues := []interface{}{
			"",
		}

		for _, urlValue := range invalidMappingTargetInlineHtmlContentValues {
			_, err := NewMappingTargetValue(urlValue, inlineHtmlContentTargetType)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", urlValue)
			}
		}
	})
}
