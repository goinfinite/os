package valueObject

import "testing"

func TestNewUrl(t *testing.T) {
	t.Run("ValidUrl", func(t *testing.T) {
		validUrls := []interface{}{
			// cSpell:disable
			"localhost", "localhost:8080", "speedia.net", "http://speedia.net/",
			"http://www.speedia.net", "https://speedia.net/", "https://www.speedia.net/",
			"http://localhost:8080/v1/ticket/253/attachment/b8680d5cc332672c649f4ff8d9e3b77f.svg",
			"https://www.cnn.com/politics/live-news/house-speaker-vote-10-20-23/index.html",
			"https://blog.goinfinite.net/2023/10/vulnerabilidades-top-10-da-owasp-parte-1/",
			"https://upload.wikimedia.org/wikipedia/commons/thumb/9/98/WordPress_blue_logo.svg/1200px-WordPress_blue_logo.svg.png",
			// cSpell:enable
		}

		for _, url := range validUrls {
			_, err := NewUrl(url)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", url, err.Error())
			}
		}
	})

	t.Run("InvalidUrl", func(t *testing.T) {
		invalidUrls := []interface{}{
			// cSpell:disable
			"", " ", "http://", "https://", "http://notãvalidurl.com/",
			"https://invalidmaçalink.com.br/", ":8080:/", "www.GoOgle.com/",
			"/home/downloads/", "DROP TABLE users;", "SELECT * FROM users;",
			"<script>alert('XSS')</script>", "http://<script>alert('XSS')</script>",
			"https://<script>alert('XSS')</script>", "rm -rf /", "(){|:& };:",
			"INSERT INTO users (name, email) VALUES ('admin', 'admin@example.com');",
			"sudo rm -r /",
			// cSpell:enable
		}

		for _, url := range invalidUrls {
			_, err := NewUrl(url)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", url)
			}
		}
	})
}
