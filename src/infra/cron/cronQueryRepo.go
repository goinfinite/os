package cronInfra

import (
	"errors"
	"log/slog"
	"slices"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type CronQueryRepo struct {
}

func NewCronQueryRepo() *CronQueryRepo {
	return &CronQueryRepo{}
}

func (repo *CronQueryRepo) cronFactory(
	cronIndex int,
	cronLine string,
) (cron entity.Cron, err error) {
	cronRegex := `^(?P<frequency>(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+) ?){5,7}))(?P<command>[^#\r\n]{1,1000})(?P<comment>#(.*)){0,1000}$`
	cronNamedGroupMap := voHelper.FindNamedGroupsMatches(cronRegex, cronLine)

	rawId := cronIndex + 1
	id, err := valueObject.NewCronId(rawId)
	if err != nil {
		return cron, err
	}

	schedule, err := valueObject.NewCronSchedule(cronNamedGroupMap["frequency"])
	if err != nil {
		return cron, err
	}

	command, err := valueObject.NewUnixCommand(cronNamedGroupMap["command"])
	if err != nil {
		return cron, err
	}

	var commentPtr *valueObject.CronComment
	if cronNamedGroupMap["comment"] != "" {
		commentWithoutLeadingHash := strings.Trim(cronNamedGroupMap["comment"], "#")
		cronComment, err := valueObject.NewCronComment(commentWithoutLeadingHash)
		if err != nil {
			return cron, err
		}
		commentPtr = &cronComment
	}

	return entity.NewCron(id, schedule, command, commentPtr), nil
}

func (repo *CronQueryRepo) readCronsFromCrontab() ([]entity.Cron, error) {
	crons := []entity.Cron{}

	cronOut, err := infraHelper.RunCmd("crontab", "-l")
	if err != nil {
		if strings.Contains(err.Error(), "no crontab") {
			return crons, nil
		}
		return crons, errors.New("CrontabReadError: " + err.Error())
	}

	cronLines := strings.Split(cronOut, "\n")
	if len(cronLines) == 0 {
		return crons, nil
	}

	for cronIndex, cronLine := range cronLines {
		if cronLine == "" {
			continue
		}

		if strings.HasPrefix(cronLine, "#") {
			continue
		}

		cron, err := repo.cronFactory(cronIndex, cronLine)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("index", cronIndex))
			continue
		}
		crons = append(crons, cron)
	}

	return crons, nil
}

func (repo *CronQueryRepo) Read(
	requestDto dto.ReadCronsRequest,
) (responseDto dto.ReadCronsResponse, err error) {
	crons, err := repo.readCronsFromCrontab()
	if err != nil {
		return responseDto, err
	}

	filteredCrons := []entity.Cron{}
	for _, cron := range crons {
		if requestDto.CronId != nil && *requestDto.CronId != cron.Id {
			continue
		}

		if requestDto.CronComment != nil && requestDto.CronComment != cron.Comment {
			continue
		}

		filteredCrons = append(filteredCrons, cron)
	}

	if len(filteredCrons) > int(requestDto.Pagination.ItemsPerPage) {
		filteredCrons = filteredCrons[:requestDto.Pagination.ItemsPerPage]
	}

	sortDirectionStr := "asc"
	if requestDto.Pagination.SortDirection != nil {
		sortDirectionStr = requestDto.Pagination.SortDirection.String()
	}

	if requestDto.Pagination.SortBy != nil {
		slices.SortStableFunc(filteredCrons, func(a, b entity.Cron) int {
			firstElement := a
			secondElement := b
			if sortDirectionStr != "asc" {
				firstElement = b
				secondElement = a
			}

			switch requestDto.Pagination.SortBy.String() {
			case "id":
				if firstElement.Id.Uint64() < secondElement.Id.Uint64() {
					return -1
				}
				if firstElement.Id.Uint64() > secondElement.Id.Uint64() {
					return 1
				}
				return 0
			case "comment":
				return strings.Compare(
					firstElement.Comment.String(), secondElement.Comment.String(),
				)
			default:
				return 0
			}
		})
	}

	paginationDto := requestDto.Pagination

	itemsTotal := uint64(len(filteredCrons))
	paginationDto.ItemsTotal = &itemsTotal

	pagesTotal := uint32(itemsTotal / uint64(requestDto.Pagination.ItemsPerPage))
	paginationDto.PagesTotal = &pagesTotal

	return dto.ReadCronsResponse{
		Pagination: paginationDto,
		Crons:      filteredCrons,
	}, nil
}

func (repo *CronQueryRepo) ReadFirst(
	requestDto dto.ReadCronsRequest,
) (cron entity.Cron, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.Read(requestDto)
	if err != nil {
		return cron, err
	}

	if len(responseDto.Crons) == 0 {
		return cron, errors.New("CronNotFound")
	}

	return responseDto.Crons[0], nil
}
