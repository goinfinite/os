package activityRecordInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
)

func TestActivityRecordQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	trailDbSvc := testHelpers.GetTrailDbSvc()
	activityRecordQueryRepo := NewActivityRecordQueryRepo(trailDbSvc)

	t.Run("ReadActivityRecordQuery", func(t *testing.T) {
		readDto := dto.ReadActivityRecords{}
		_, err := activityRecordQueryRepo.Read(readDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
