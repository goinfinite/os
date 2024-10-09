package activityRecordInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestActivityRecordCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	trailDbSvc := testHelpers.GetTrailDbSvc()
	activityRecordCmdRepo := NewActivityRecordCmdRepo(trailDbSvc)

	level, _ := valueObject.NewActivityRecordLevel("SEC")
	recordCode, _ := valueObject.NewActivityRecordCode("LoginFailed")
	ipAddress := valueObject.NewLocalhostIpAddress()

	t.Run("CreateActivityRecord", func(t *testing.T) {
		username, _ := valueObject.NewUsername("test")
		createDto := dto.CreateActivityRecord{
			Level:     level,
			Code:      &recordCode,
			IpAddress: &ipAddress,
			Username:  &username,
		}

		err := activityRecordCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("Expected no error, got: '%s'", err.Error())
		}
	})

	t.Run("DeleteActivityRecords", func(t *testing.T) {
		ipAddress := valueObject.NewLocalhostIpAddress()
		deleteDto := dto.NewDeleteActivityRecords(
			nil, &level, &recordCode, nil, &ipAddress, nil, nil, nil, nil, nil,
		)

		err := activityRecordCmdRepo.Delete(deleteDto)
		if err != nil {
			t.Errorf("Expected no error, got: '%s'", err.Error())
		}
	})
}
