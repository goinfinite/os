package activityRecordInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
)

func TestActivityRecordQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	trailDbSvc := testHelpers.GetTrailDbSvc()
	activityRecordQueryRepo := NewActivityRecordQueryRepo(trailDbSvc)

	t.Run("ReadActivityRecordQuery", func(t *testing.T) {
		readDto := dto.ReadActivityRecords{}
		_, err := activityRecordQueryRepo.Read(readDto)
		if err != nil {
			t.Errorf("Expected no error, got: '%s'", err.Error())
		}
	})
}
