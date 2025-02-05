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

	var schedule valueObject.CronSchedule
	if _, exists := cronNamedGroupMap["frequency"]; exists {
		schedule, err = valueObject.NewCronSchedule(cronNamedGroupMap["frequency"])
		if err != nil {
			return cron, err
		}
	}

	var command valueObject.UnixCommand
	if _, exists := cronNamedGroupMap["command"]; exists {
		command, err = valueObject.NewUnixCommand(cronNamedGroupMap["command"])
		if err != nil {
			return cron, err
		}
	}

	var commentPtr *valueObject.CronComment
	if _, exists := cronNamedGroupMap["comment"]; exists {
		commentWithoutLeadingHash := strings.Trim(cronNamedGroupMap["comment"], "#")
		comment, err := valueObject.NewCronComment(commentWithoutLeadingHash)
		if err != nil {
			return cron, err
		}
		commentPtr = &comment
	}

	return entity.NewCron(id, schedule, command, commentPtr), nil
}

func (repo *CronQueryRepo) readCronsFromCrontab() ([]entity.Cron, error) {
	crons := []entity.Cron{}

	rawCronOutput, err := infraHelper.RunCmd("crontab", "-l")
	if err != nil {
		if strings.Contains(err.Error(), "no crontab") {
			return crons, nil
		}
		return crons, errors.New("CrontabReadError: " + err.Error())
	}

	cronLines := strings.Split(rawCronOutput, "\n")
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

		cronEntity, err := repo.cronFactory(cronIndex, cronLine)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("index", cronIndex))
			continue
		}
		crons = append(crons, cronEntity)
	}

	return crons, nil
}

func (repo *CronQueryRepo) Read(
	requestDto dto.ReadCronsRequest,
) (responseDto dto.ReadCronsResponse, err error) {
	originalCronEntities, err := repo.readCronsFromCrontab()
	if err != nil {
		return responseDto, err
	}

	filteredCrons := []entity.Cron{}
	for _, cron := range originalCronEntities {
		if requestDto.CronId != nil && requestDto.CronId.Uint64() != cron.Id.Uint64() {
			continue
		}

		hasCronCommentToFilter := requestDto.CronComment != nil && cron.Comment != nil
		if hasCronCommentToFilter && requestDto.CronComment.String() != cron.Comment.String() {
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
