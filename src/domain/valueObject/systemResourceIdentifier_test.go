package valueObject

import "testing"

func TestSystemResourceIdentifier(t *testing.T) {
	t.Run("ValidSystemResourceIdentifier", func(t *testing.T) {
		validSystemResourceIdentifier := []interface{}{
			"sri://1:account/120", "sri://10:secureAccessPublicKey/1", "sri://100:cron/1",
			"sri://1000:database/myDb", "sri://1:databaseUser/myDbUser",
			"sri://10:marketplaceCatalogItem/1", "sri://100:marketplaceCatalogItem/php",
			"sri://1000:marketplaceInstalledItem/1", "sri://1:phpRuntime/local.os",
			"sri://10:installableService/node", "sri://100:customService/node-e87qxc21",
			"sri://1000:installedService/1", "sri://1:ssl/1", "sri://10:virtualHost/local.os",
			"sri://100:mapping/1", "sri://1000:unixFile//app/.trash",
		}

		for _, identifier := range validSystemResourceIdentifier {
			_, err := NewSystemResourceIdentifier(identifier)
			if err != nil {
				t.Errorf(
					"Expected no error for '%v', got '%s'", identifier, err.Error(),
				)
			}
		}
	})

	t.Run("InvalidSystemResourceIdentifier", func(t *testing.T) {
		invalidSystemResourceIdentifier := []interface{}{
			"", "sri://0:/", true, 1000,
		}

		for _, identifier := range invalidSystemResourceIdentifier {
			_, err := NewSystemResourceIdentifier(identifier)
			if err == nil {
				t.Errorf("Expected error for '%v', got nil", identifier)
			}
		}
	})
}
