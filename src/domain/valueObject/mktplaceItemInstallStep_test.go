package valueObject

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestMktplaceItemInstallStep(t *testing.T) {
	t.Run("ValidMktplaceItemInstallStep", func(t *testing.T) {
		validMktplaceItemInstallSteps := []string{
			"ls -l",
			"cat file.txt | grep \"pattern\" | sort",
			"echo \"Today is $(date +%A)\"",
			"mkdir test_directory && cd test_directory && touch file1.txt file2.txt && ls",
			"certbot certonly --webroot --webroot-path /app/html --agree-tos --register-unsafely-without-email --cert-name speedia.net -d speedia.net",
			"wget https://github.com/speedianet/os -O $PATH",
		}

		for _, mktplaceItemInstallStep := range validMktplaceItemInstallSteps {
			_, err := NewMktplaceItemInstallStep(mktplaceItemInstallStep)
			if err != nil {
				t.Errorf(
					"Expected no error for %s, got %s",
					mktplaceItemInstallStep,
					err.Error(),
				)
			}
		}
	})

	t.Run("ValidMktplaceItemInstallStep", func(t *testing.T) {
		invalidLength := 700
		invalidMktplaceItemInstallSteps := []string{
			"",
			testHelpers.GenerateString(invalidLength),
		}

		for _, mktplaceItemInstallStep := range invalidMktplaceItemInstallSteps {
			_, err := NewMktplaceItemInstallStep(mktplaceItemInstallStep)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", mktplaceItemInstallStep)
			}
		}
	})
}
