package activityRecordInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

func TestActivityRecordCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	trailDbSvc := testHelpers.GetTrailDbSvc()
	activityRecordCmdRepo := NewActivityRecordCmdRepo(trailDbSvc)
	level, _ := valueObject.NewActivityRecordLevel("SEC")
	recordCode, _ := valueObject.NewActivityRecordCode("LoginFailed")
	operatorIpAddress := valueObject.NewLocalhostIpAddress()

	t.Run("CreateActivityRecord", func(t *testing.T) {
		createDto := dto.CreateActivityRecord{
			RecordLevel:       level,
			RecordCode:        recordCode,
			OperatorIpAddress: &operatorIpAddress,
		}

		err := activityRecordCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})

	t.Run("DeleteActivityRecords", func(t *testing.T) {
		ipAddress := valueObject.NewLocalhostIpAddress()
		deleteDto := dto.NewDeleteActivityRecord(
			nil, &level, &recordCode, nil, nil, &ipAddress, nil, nil,
		)

		err := activityRecordCmdRepo.Delete(deleteDto)
		if err != nil {
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
