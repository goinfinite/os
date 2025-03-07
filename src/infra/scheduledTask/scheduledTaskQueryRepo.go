package scheduledTaskInfra

import (
	"errors"
	"log/slog"
	"math"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	"github.com/iancoleman/strcase"
)

type ScheduledTaskQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewScheduledTaskQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskQueryRepo {
	return &ScheduledTaskQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ScheduledTaskQueryRepo) Read(
	readDto dto.ReadScheduledTasksRequest,
) (responseDto dto.ReadScheduledTasksResponse, err error) {
	scheduledTaskEntities := []entity.ScheduledTask{}

	scheduledTaskModel := dbModel.ScheduledTask{}
	if readDto.TaskId != nil {
		scheduledTaskModel.ID = readDto.TaskId.Uint64()
	}
	if readDto.TaskName != nil {
		scheduledTaskModel.Name = readDto.TaskName.String()
	}
	if readDto.TaskStatus != nil {
		scheduledTaskModel.Status = readDto.TaskStatus.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Model(&scheduledTaskModel).
		Where(&scheduledTaskModel)
	if len(readDto.TaskTags) == 0 {
		dbQuery = dbQuery.Preload("Tags")
	} else {
		tagsStrSlice := []string{}
		for _, taskTag := range readDto.TaskTags {
			tagsStrSlice = append(tagsStrSlice, taskTag.String())
		}
		dbQuery = dbQuery.
			Joins("JOIN scheduled_tasks_tags ON scheduled_tasks_tags.scheduled_task_id = scheduled_tasks.id").
			Where("scheduled_tasks_tags.tag IN (?)", tagsStrSlice)
	}
	if readDto.StartedBeforeAt != nil {
		dbQuery = dbQuery.Where("started_at < ?", readDto.StartedBeforeAt.ReadAsGoTime())
	}
	if readDto.StartedAfterAt != nil {
		dbQuery = dbQuery.Where("started_at > ?", readDto.StartedAfterAt.ReadAsGoTime())
	}
	if readDto.FinishedBeforeAt != nil {
		dbQuery = dbQuery.Where("finished_at < ?", readDto.FinishedBeforeAt.ReadAsGoTime())
	}
	if readDto.FinishedAfterAt != nil {
		dbQuery = dbQuery.Where("finished_at > ?", readDto.FinishedAfterAt.ReadAsGoTime())
	}
	if readDto.CreatedBeforeAt != nil {
		dbQuery = dbQuery.Where("created_at < ?", readDto.CreatedBeforeAt.ReadAsGoTime())
	}
	if readDto.CreatedAfterAt != nil {
		dbQuery = dbQuery.Where("created_at > ?", readDto.CreatedAfterAt.ReadAsGoTime())
	}

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return responseDto, errors.New("CountItemsTotalError: " + err.Error())
	}

	dbQuery = dbQuery.Limit(int(readDto.Pagination.ItemsPerPage))
	if readDto.Pagination.LastSeenId == nil {
		offset := int(readDto.Pagination.PageNumber) * int(readDto.Pagination.ItemsPerPage)
		dbQuery = dbQuery.Offset(offset)
	} else {
		dbQuery = dbQuery.Where("id > ?", readDto.Pagination.LastSeenId.String())
	}
	if readDto.Pagination.SortBy != nil {
		orderStatement := readDto.Pagination.SortBy.String()
		orderStatement = strcase.ToSnake(orderStatement)
		if orderStatement == "id" {
			orderStatement = "ID"
		}

		if readDto.Pagination.SortDirection != nil {
			orderStatement += " " + readDto.Pagination.SortDirection.String()
		}

		dbQuery = dbQuery.Order(orderStatement)
	}

	scheduledTaskModels := []dbModel.ScheduledTask{}
	err = dbQuery.Find(&scheduledTaskModels).Error
	if err != nil {
		return responseDto, errors.New("FindScheduledTasksError: " + err.Error())
	}

	for _, scheduledTaskModel := range scheduledTaskModels {
		scheduledTaskEntity, err := scheduledTaskModel.ToEntity()
		if err != nil {
			slog.Debug(
				"ModelToEntityError",
				slog.Uint64("id", scheduledTaskModel.ID),
				slog.String("err", err.Error()),
			)
			continue
		}
		scheduledTaskEntities = append(scheduledTaskEntities, scheduledTaskEntity)
	}

	itemsTotalUint := uint64(itemsTotal)
	pagesTotal := uint32(
		math.Ceil(float64(itemsTotal) / float64(readDto.Pagination.ItemsPerPage)),
	)
	responsePagination := dto.Pagination{
		PageNumber:    readDto.Pagination.PageNumber,
		ItemsPerPage:  readDto.Pagination.ItemsPerPage,
		SortBy:        readDto.Pagination.SortBy,
		SortDirection: readDto.Pagination.SortDirection,
		PagesTotal:    &pagesTotal,
		ItemsTotal:    &itemsTotalUint,
	}

	return dto.ReadScheduledTasksResponse{
		Pagination: responsePagination,
		Tasks:      scheduledTaskEntities,
	}, nil
}
