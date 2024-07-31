package valueObject

import "testing"

func TestNewMappingTargetValue(t *testing.T) {
	t.Run("ValidMappingTargetValueBasedOnType (Url)", func(t *testing.T) {
		urlTargetType, _ := NewMappingTargetType("url")
		validMappingTargetUrlValues := []interface{}{
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
		}

		for _, mtv := range validMappingTargetUrlValues {
			_, err := NewMappingTargetValue(mtv, urlTargetType)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mtv, err.Error())
			}
		}
	})

	t.Run("InvalidMappingTargetValueBasedOnType (Url)", func(t *testing.T) {
		urlTargetType, _ := NewMappingTargetType("url")
		invalidMappingTargetUrlValues := []interface{}{
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
		}

		for _, mtv := range invalidMappingTargetUrlValues {
			_, err := NewMappingTargetValue(mtv, urlTargetType)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mtv)
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

		for _, mtv := range validMappingTargetServiceNameValues {
			_, err := NewMappingTargetValue(mtv, svcNameTargetType)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mtv, err.Error())
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

		for _, mtv := range invalidMappingTargetServiceNameValues {
			_, err := NewMappingTargetValue(mtv, svcNameTargetType)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mtv)
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

		for _, mtv := range validMappingTargetResponseCodeValues {
			_, err := NewMappingTargetValue(mtv, responseCodeTargetType)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mtv, err.Error())
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

		for _, mtv := range invalidMappingTargetResponseCodeValues {
			_, err := NewMappingTargetValue(mtv, responseCodeTargetType)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mtv)
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

		for _, mtv := range validMappingTargetInlineHtmlContentValues {
			_, err := NewMappingTargetValue(mtv, inlineHtmlContentTargetType)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", mtv, err.Error())
			}
		}
	})

	t.Run("InvalidMappingTargetValueBasedOnType (Inline HTML Content)", func(t *testing.T) {
		inlineHtmlContentTargetType, _ := NewMappingTargetType("inline-html")
		invalidMappingTargetInlineHtmlContentValues := []interface{}{
			"",
		}

		for _, mtv := range invalidMappingTargetInlineHtmlContentValues {
			_, err := NewMappingTargetValue(mtv, inlineHtmlContentTargetType)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mtv)
			}
		}
	})
}
