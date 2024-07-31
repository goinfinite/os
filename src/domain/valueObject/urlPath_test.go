package valueObject

import "testing"

func TestUrlPath(t *testing.T) {
	t.Run("ValidUrlPath", func(t *testing.T) {
		validUrlPath := []interface{}{
			"", "/", "blog",
			"news/new-product-from-Speedia-revolutionizes-the-market",
			"/app/html", "/info.php", "/app/html/speedia.net",
			"/v1/ticket/253/attachment/b8680d5cc332672c649f4ff8d9e3b77f.svg",
			"/politics/live-news/house-speaker-vote-10-20-23/index.html",
			"/2023/10/vulnerabilidades-top-10-da-owasp-parte-1/",
			"/wikipedia/commons/thumb/9/98/WordPress_blue_logo.svg/1200px-WordPress_blue_logo.svg.png",
		}

		for _, urlPath := range validUrlPath {
			_, err := NewUrlPath(urlPath)
			if err != nil {
				t.Errorf("Expected no error for '%v', got '%s'", urlPath, err.Error())
			}
		}
	})

	t.Run("InvalidUrlPath", func(t *testing.T) {
		invalidUrlPath := []interface{}{
			"/app/html@", "/info.php?id=1", "/path to download", "index.js=",
			"/how-to-get-habbo-coins?/2011",
		}

		for _, urlPath := range invalidUrlPath {
			_, err := NewUrlPath(urlPath)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", urlPath)
			}
		}
	})
}
