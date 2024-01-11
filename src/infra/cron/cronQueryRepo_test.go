package cronInfra

import (
	"testing"
)

func TestCronQueryRepo(t *testing.T) {
	t.Run("GetCrons", func(t *testing.T) {
		cronQueryRepo := CronQueryRepo{}

		_, err := cronQueryRepo.Get()
		if err != nil {
			t.Errorf("GetCrons should not return: %s", err)
		}
	})
}
