package valueObject

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestMarketplaceItemCmdStep(t *testing.T) {
	t.Run("ValidMarketplaceItemCmdStep", func(t *testing.T) {
		validMarketplaceItemCmdSteps := []string{
			"ls -l",
			"cat file.txt | grep \"pattern\" | sort",
			"echo \"Today is $(date +%A)\"",
			"mkdir test_directory && cd test_directory && touch file1.txt file2.txt && ls",
			"certbot certonly --webroot --webroot-path /app/html --agree-tos --register-unsafely-without-email --cert-name speedia.net -d speedia.net",
			"wget https://github.com/speedianet/os -O $PATH",
		}

		for _, miis := range validMarketplaceItemCmdSteps {
			_, err := NewMarketplaceItemCmdStep(miis)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", miis, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceItemCmdStep", func(t *testing.T) {
		invalidLength := 4100
		invalidMarketplaceItemCmdSteps := []string{
			"",
			testHelpers.GenerateString(invalidLength),
		}

		for _, miis := range invalidMarketplaceItemCmdSteps {
			_, err := NewMarketplaceItemCmdStep(miis)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", miis)
			}
		}
	})
}
