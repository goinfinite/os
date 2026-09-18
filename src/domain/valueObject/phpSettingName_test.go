package valueObject

import "testing"

func TestPhpSettingName(t *testing.T) {
	t.Run("ValidPhpSettingNames", func(t *testing.T) {
		validNames := []any{
			"ioncube", "apcu", "imagick", "opcache", "mysqli",
			"date.timezone", "soap.wsdl_cache_enabled", "eaccelerator-off",
		}

		for _, name := range validNames {
			_, err := NewPhpSettingName(name)
			if err != nil {
				t.Errorf("UnexpectedError: %v, name: '%v'", err.Error(), name)
			}
		}
	})

	t.Run("InvalidPhpSettingNames", func(t *testing.T) {
		invalidNames := []any{
			"ioncube_loader.so!", "<script>alert('xss')</script>", "@blabla@",
			".leadingDot", "trailingDot.", "spaces in name",
		}

		for _, name := range invalidNames {
			_, err := NewPhpSettingName(name)
			if err == nil {
				t.Errorf("MissingExpectedError: '%v' was accepted", name)
			}
		}
	})
}
