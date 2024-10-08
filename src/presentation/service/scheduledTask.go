package service

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ScheduledTaskService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewScheduledTaskService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskService {
	return &ScheduledTaskService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *ScheduledTaskService) Read(input map[string]interface{}) ServiceOutput {
	var taskIdPtr *valueObject.ScheduledTaskId
	if input["id"] != nil {
		input["taskId"] = input["id"]
	}
	if input["taskId"] != nil {
		taskId, err := valueObject.NewScheduledTaskId(input["taskId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		taskIdPtr = &taskId
	}

	var taskNamePtr *valueObject.ScheduledTaskName
	if input["taskName"] != nil {
		taskName, err := valueObject.NewScheduledTaskName(input["taskName"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		taskNamePtr = &taskName
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if input["taskStatus"] != nil {
		taskStatus, err := valueObject.NewScheduledTaskStatus(input["taskStatus"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		taskStatusPtr = &taskStatus
	}

	taskTags := []valueObject.ScheduledTaskTag{}
	if input["taskTags"] != nil {
		var assertOk bool
		taskTags, assertOk = input["taskTags"].([]valueObject.ScheduledTaskTag)
		if !assertOk {
			return NewServiceOutput(UserError, errors.New("InvalidTaskTags"))
		}
	}

	var startedBeforeAtPtr, startedAfterAtPtr *valueObject.UnixTime
	var finishedBeforeAtPtr, finishedAfterAtPtr *valueObject.UnixTime
	var createdBeforeAtPtr, createdAfterAtPtr *valueObject.UnixTime

	timeParamNames := []string{
		"startedBeforeAt", "startedAfterAt",
		"finishedBeforeAt", "finishedAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	for _, timeParamName := range timeParamNames {
		if input[timeParamName] == nil {
			continue
		}

		timeParam, err := valueObject.NewUnixTime(input[timeParamName])
		if err != nil {
			capitalParamName := cases.Title(language.English).String(timeParamName)
			return NewServiceOutput(UserError, errors.New("Invalid"+capitalParamName))
		}

		switch timeParamName {
		case "startedBeforeAt":
			startedBeforeAtPtr = &timeParam
		case "startedAfterAt":
			startedAfterAtPtr = &timeParam
		case "finishedBeforeAt":
			finishedBeforeAtPtr = &timeParam
		case "finishedAfterAt":
			finishedAfterAtPtr = &timeParam
		case "createdBeforeAt":
			createdBeforeAtPtr = &timeParam
		case "createdAfterAt":
			createdAfterAtPtr = &timeParam
		}
	}

	paginationDto := useCase.ScheduledTasksDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(input["sortDirection"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadScheduledTasksRequest{
		Pagination:       paginationDto,
		TaskId:           taskIdPtr,
		TaskName:         taskNamePtr,
		TaskStatus:       taskStatusPtr,
		TaskTags:         taskTags,
		StartedBeforeAt:  startedBeforeAtPtr,
		StartedAfterAt:   startedAfterAtPtr,
		FinishedBeforeAt: finishedBeforeAtPtr,
		FinishedAfterAt:  finishedAfterAtPtr,
		CreatedBeforeAt:  createdBeforeAtPtr,
		CreatedAfterAt:   createdAfterAtPtr,
	}

	scheduledTaskQueryRepo := scheduledTaskInfra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTasksList, err := useCase.ReadScheduledTasks(scheduledTaskQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, scheduledTasksList)
}

func (service *ScheduledTaskService) Update(input map[string]interface{}) ServiceOutput {
	if input["id"] != nil {
		input["taskId"] = input["id"]
	}

	requiredParams := []string{"taskId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	taskId, err := valueObject.NewScheduledTaskId(input["taskId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if input["status"] != nil {
		taskStatus, err := valueObject.NewScheduledTaskStatus(input["status"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	var runAtPtr *valueObject.UnixTime
	if input["runAt"] != nil {
		runAt, err := valueObject.NewUnixTime(input["runAt"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		runAtPtr = &runAt
	}

	updateDto := dto.NewUpdateScheduledTask(
		taskId, taskStatusPtr, runAtPtr,
	)

	scheduledTaskQueryRepo := scheduledTaskInfra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(service.persistentDbSvc)

	err = useCase.UpdateScheduledTask(
		scheduledTaskQueryRepo, scheduledTaskCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ScheduledTaskUpdated")
}
