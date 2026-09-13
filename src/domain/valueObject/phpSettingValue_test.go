package valueObject

import "testing"

func TestPhpSettingValue(t *testing.T) {
	t.Run("ValidPhpSettingValues", func(t *testing.T) {
		validValues := []any{
			"on", "off", "ON", "OFF", "true", "false", "TRUE", "FALSE", true, false,
			0, 1, 2, "test", "dev", "prod",
			"/usr/sbin/sendmail -t -i", "E_ALL & ~E_DEPRECATED & ~E_STRICT",
		}

		for _, value := range validValues {
			_, err := NewPhpSettingValue(value)
			if err != nil {
				t.Errorf("UnexpectedError: %v, value: '%v'", err.Error(), value)
			}
		}
	})

	t.Run("InvalidPhpSettingValues", func(t *testing.T) {
		invalidValues := []any{
			"", "display_errors\nallow_url_fopen On", "date.timezone\tUTC", "log\x00",
		}

		for _, value := range invalidValues {
			_, err := NewPhpSettingValue(value)
			if err == nil {
				t.Errorf("MissingExpectedError: '%v' was accepted", value)
			}
		}
	})

	t.Run("UnsetPhpSettingValueIsTypedAsString", func(t *testing.T) {
		unsetValue := PhpSettingValue("")
		if unsetValue.ReadType() != PhpSettingValueTypeString {
			t.Errorf(
				"ExpectedStringTypeForUnsetValue: got '%s'",
				unsetValue.ReadType(),
			)
		}
	})

	t.Run("ByteSizeUnitsAreNormalizedToUppercase", func(t *testing.T) {
		byteSizes := map[string]string{
			"128m":  "128M",
			"512k":  "512K",
			"2g":    "2G",
			"4096K": "4096K",
		}

		for rawValue, expectedValue := range byteSizes {
			settingValue, err := NewPhpSettingValue(rawValue)
			if err != nil {
				t.Errorf("UnexpectedError: %v, value: '%s'", err.Error(), rawValue)
				continue
			}
			if settingValue.String() != expectedValue {
				t.Errorf(
					"ByteSizeNormalizationMismatch: expected '%s' for '%s', got '%s'",
					expectedValue,
					rawValue,
					settingValue.String(),
				)
			}
			if settingValue.ReadType() != PhpSettingValueTypeByteSize {
				t.Errorf(
					"ExpectedByteSizeType: value '%s', got type '%s'",
					expectedValue,
					settingValue.ReadType(),
				)
			}
		}
	})

	t.Run("ValuesWithByteSizeLookalikesAreNotByteSizes", func(t *testing.T) {
		lookalikes := []string{"OK", "128X", "M", "12.8M", "-1M", "128MB"}

		for _, lookalike := range lookalikes {
			settingValue, err := NewPhpSettingValue(lookalike)
			if err != nil {
				t.Errorf("UnexpectedError: %v, value: '%s'", err.Error(), lookalike)
				continue
			}
			if settingValue.ReadType() != PhpSettingValueTypeString {
				t.Errorf(
					"ExpectedStringType: value '%s', got type '%s'",
					lookalike,
					settingValue.ReadType(),
				)
			}
		}
	})
}
