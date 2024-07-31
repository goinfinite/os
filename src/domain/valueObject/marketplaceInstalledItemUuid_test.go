package valueObject

import "testing"

func TestMarketplaceInstalledItemUuid(t *testing.T) {
	t.Run("ValidMarketplaceInstalledItemUuid", func(t *testing.T) {
		validUuids := []interface{}{
			"abc123def4",
			"1234567890ab",
			"abcdef123456",
			"9876543210ab",
			"1234abcd5678",
		}

		for _, uuid := range validUuids {
			_, err := NewMarketplaceInstalledItemUuid(uuid)
			if err != nil {
				t.Errorf("Expected no error for %s, got %s", uuid, err.Error())
			}
		}
	})

	t.Run("InvalidMarketplaceInstalledItemUuid", func(t *testing.T) {
		invalidUuids := []interface{}{
			"abc123",
			"tolongmarketplaceinstalleditemuuid",
			"12345678!@#",
			"short12",
		}

		for _, uuid := range invalidUuids {
			_, err := NewMarketplaceInstalledItemUuid(uuid)
			if err == nil {
				t.Errorf("Expected error for %s, got nil", uuid)
			}
		}
	})
}
