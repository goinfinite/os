package cronInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func TestCronQueryRepo(t *testing.T) {
	cronQueryRepo := NewCronQueryRepo()
	testHelpers.LoadEnvVars()

	schedule, _ := valueObject.NewCronSchedule("* * * * *")
	command, _ := valueObject.NewUnixCommand("echo \"cronTest\" >> crontab_log.txt")
	comment, _ := valueObject.NewCronComment("Test cron job")
	ipAddress := valueObject.IpAddressSystem
	operatorAccountId, _ := valueObject.NewAccountId(0)

	createDto := dto.NewCreateCron(
		schedule, command, &comment, operatorAccountId, ipAddress,
	)

	cronCmdRepo := NewCronCmdRepo()
	cronId, err := cronCmdRepo.Create(createDto)
	if err != nil {
		t.Fatalf("Expected no error but got '%s'", err.Error())
	}

	t.Run("Read", func(t *testing.T) {
		paginationDto := useCase.CronsDefaultPagination

		sortBy, _ := valueObject.NewPaginationSortBy("id")
		paginationDto.SortBy = &sortBy

		sortDirection, _ := valueObject.NewPaginationSortDirection("desc")
		paginationDto.SortDirection = &sortDirection

		readRequestDto := dto.ReadCronsRequest{
			Pagination: paginationDto,
			CronId:     &cronId,
		}

		responseDto, err := cronQueryRepo.Read(readRequestDto)
		if err != nil {
			t.Fatalf("Expected no error, but got '%s'", err.Error())
		}

		if len(responseDto.Crons) == 0 {
			t.Error("Expected one cronjob at least, but got 0")
		}
	})

	t.Run("ReadFirst", func(t *testing.T) {
		readRequestDto := dto.ReadCronsRequest{
			CronId: &cronId,
		}

		_, err := cronQueryRepo.ReadFirst(readRequestDto)
		if err != nil {
			t.Errorf("Expected no error, but got '%s'", err.Error())
		}
	})
}
