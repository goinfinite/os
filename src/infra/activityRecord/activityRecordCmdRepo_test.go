package infra

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
			t.Errorf("ExpectedNoErrorButGot: %v", err)
		}
	})
}
