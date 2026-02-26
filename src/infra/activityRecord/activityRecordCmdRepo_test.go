package activityRecordInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestActivityRecordCmdRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	trailDbSvc := testHelpers.GetTrailDbSvc()
	activityRecordCmdRepo := NewActivityRecordCmdRepo(trailDbSvc)
	level, _ := tkValueObject.NewActivityRecordLevel("SEC")
	recordCode, _ := tkValueObject.NewActivityRecordCode("LoginFailed")
	operatorIpAddress := tkValueObject.IpAddressLocal

	t.Run("CreateActivityRecord", func(t *testing.T) {
		createDto := tkDto.CreateActivityRecord{
			RecordLevel:       level,
			RecordCode:        recordCode,
			OperatorIpAddress: &operatorIpAddress,
		}

		err := activityRecordCmdRepo.Create(createDto)
		if err != nil {
			t.Errorf("Expected no error, got: '%s'", err.Error())
		}
	})

	t.Run("DeleteActivityRecords", func(t *testing.T) {
		ipAddress := tkValueObject.IpAddressLocal
		deleteDto := tkDto.NewDeleteActivityRecord(
			nil, &level, &recordCode, nil, nil, &ipAddress, nil, nil,
		)

		err := activityRecordCmdRepo.Delete(deleteDto)
		if err != nil {
			t.Errorf("Expected no error, got: '%s'", err.Error())
		}
	})
}
