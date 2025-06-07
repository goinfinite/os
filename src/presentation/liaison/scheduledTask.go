package liaison

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ScheduledTaskLiaison struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewScheduledTaskLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskLiaison {
	return &ScheduledTaskLiaison{
		persistentDbSvc: persistentDbSvc,
	}
}

func (liaison *ScheduledTaskLiaison) Read(untrustedInput map[string]any) LiaisonOutput {
	var taskIdPtr *valueObject.ScheduledTaskId
	if untrustedInput["id"] != nil {
		untrustedInput["taskId"] = untrustedInput["id"]
	}
	if untrustedInput["taskId"] != nil {
		taskId, err := valueObject.NewScheduledTaskId(untrustedInput["taskId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		taskIdPtr = &taskId
	}

	var taskNamePtr *valueObject.ScheduledTaskName
	if untrustedInput["taskName"] != nil {
		taskName, err := valueObject.NewScheduledTaskName(untrustedInput["taskName"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		taskNamePtr = &taskName
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if untrustedInput["taskStatus"] != nil {
		taskStatus, err := valueObject.NewScheduledTaskStatus(untrustedInput["taskStatus"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		taskStatusPtr = &taskStatus
	}

	taskTags := []valueObject.ScheduledTaskTag{}
	if untrustedInput["taskTags"] != nil {
		var assertOk bool
		taskTags, assertOk = untrustedInput["taskTags"].([]valueObject.ScheduledTaskTag)
		if !assertOk {
			return NewLiaisonOutput(UserError, errors.New("InvalidTaskTags"))
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
		if untrustedInput[timeParamName] == nil {
			continue
		}

		timeParam, err := valueObject.NewUnixTime(untrustedInput[timeParamName])
		if err != nil {
			capitalParamName := cases.Title(language.English).String(timeParamName)
			return NewLiaisonOutput(UserError, errors.New("Invalid"+capitalParamName))
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
	if untrustedInput["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(untrustedInput["pageNumber"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if untrustedInput["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(untrustedInput["itemsPerPage"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if untrustedInput["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(untrustedInput["sortBy"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if untrustedInput["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(untrustedInput["sortDirection"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if untrustedInput["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(untrustedInput["lastSeenId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
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

	scheduledTaskQueryRepo := scheduledTaskInfra.NewScheduledTaskQueryRepo(liaison.persistentDbSvc)
	scheduledTasksList, err := useCase.ReadScheduledTasks(scheduledTaskQueryRepo, readDto)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, scheduledTasksList)
}

func (liaison *ScheduledTaskLiaison) Update(untrustedInput map[string]any) LiaisonOutput {
	if untrustedInput["id"] != nil {
		untrustedInput["taskId"] = untrustedInput["id"]
	}

	requiredParams := []string{"taskId"}

	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	taskId, err := valueObject.NewScheduledTaskId(untrustedInput["taskId"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if untrustedInput["status"] != nil {
		taskStatus, err := valueObject.NewScheduledTaskStatus(untrustedInput["status"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	var runAtPtr *valueObject.UnixTime
	if untrustedInput["runAt"] != nil {
		runAt, err := valueObject.NewUnixTime(untrustedInput["runAt"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		runAtPtr = &runAt
	}

	updateDto := dto.NewUpdateScheduledTask(
		taskId, taskStatusPtr, runAtPtr,
	)

	scheduledTaskQueryRepo := scheduledTaskInfra.NewScheduledTaskQueryRepo(liaison.persistentDbSvc)
	scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(liaison.persistentDbSvc)

	err = useCase.UpdateScheduledTask(
		scheduledTaskQueryRepo, scheduledTaskCmdRepo, updateDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "ScheduledTaskUpdated")
}
