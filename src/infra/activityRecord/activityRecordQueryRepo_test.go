package activityRecordInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
)

func TestActivityRecordQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	trailDbSvc := testHelpers.GetTrailDbSvc()
	activityRecordQueryRepo := NewActivityRecordQueryRepo(trailDbSvc)

	t.Run("ReadActivityRecordQuery", func(t *testing.T) {
		readDto := tkDto.ReadActivityRecordsRequest{
			Pagination: tkDto.PaginationUnpaginated,
		}
		_, err := activityRecordQueryRepo.Read(readDto)
		if err != nil {
			t.Errorf("Expected no error, got: '%s'", err.Error())
		}
	})
}
