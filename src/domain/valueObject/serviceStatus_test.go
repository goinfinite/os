package valueObject

import (
	"testing"
)

func TestServiceStatus(t *testing.T) {
	t.Run("ValidServiceStatuses", func(t *testing.T) {
		validStatusesAndAliases := []string{
			"running",
			"run",
			"up",
			"true",
			"false",
			"off",
			"no",
			"stop",
			"stopped",
			"halt",
			"uninstall",
			"uninstalled",
			"remove",
			"purge",
		}
		for _, status := range validStatusesAndAliases {
			_, err := NewServiceStatus(status)
			if err != nil {
				t.Errorf("(%s) ExpectedNoErrorButGot: %s", status, err.Error())
			}
		}
	})

	t.Run("InvalidServiceStatuses", func(t *testing.T) {
		invalidStatusesAndAliases := []string{
			"runningg",
			"runn",
			"upp",
			"truee",
			"falsee",
			"offf",
			"runn1ng",
			"un11install",
		}
		for _, status := range invalidStatusesAndAliases {
			_, err := NewServiceStatus(status)
			if err == nil {
				t.Errorf("(%s) ExpectedErrorButGotNil", status)
			}
		}
	})
}
