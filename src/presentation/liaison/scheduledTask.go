package liaison

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
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

func (liaison *ScheduledTaskLiaison) Read(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	var taskIdPtr *valueObject.ScheduledTaskId
	if untrustedInput["id"] != nil {
		untrustedInput["taskId"] = untrustedInput["id"]
	}
	if untrustedInput["taskId"] != nil {
		taskId, err := valueObject.NewScheduledTaskId(untrustedInput["taskId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err)
		}
		taskIdPtr = &taskId
	}

	var taskNamePtr *valueObject.ScheduledTaskName
	if untrustedInput["taskName"] != nil {
		taskName, err := valueObject.NewScheduledTaskName(untrustedInput["taskName"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err)
		}
		taskNamePtr = &taskName
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if untrustedInput["taskStatus"] != nil {
		taskStatus, err := valueObject.NewScheduledTaskStatus(untrustedInput["taskStatus"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err)
		}
		taskStatusPtr = &taskStatus
	}

	taskTags := []valueObject.ScheduledTaskTag{}
	if untrustedInput["taskTags"] != nil {
		var assertOk bool
		taskTags, assertOk = untrustedInput["taskTags"].([]valueObject.ScheduledTaskTag)
		if !assertOk {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, errors.New("InvalidTaskTags"))
		}
	}

	var startedBeforeAtPtr, startedAfterAtPtr *tkValueObject.UnixTime
	var finishedBeforeAtPtr, finishedAfterAtPtr *tkValueObject.UnixTime
	var createdBeforeAtPtr, createdAfterAtPtr *tkValueObject.UnixTime

	timeParamNames := []string{
		"startedBeforeAt", "startedAfterAt",
		"finishedBeforeAt", "finishedAfterAt",
		"createdBeforeAt", "createdAfterAt",
	}
	for _, timeParamName := range timeParamNames {
		if untrustedInput[timeParamName] == nil {
			continue
		}

		timeParam, err := tkValueObject.NewUnixTime(untrustedInput[timeParamName])
		if err != nil {
			capitalParamName := cases.Title(language.English).String(timeParamName)
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, errors.New("Invalid"+capitalParamName))
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

	requestPagination, err := tkPresentation.PaginationParser(
		useCase.ScheduledTasksDefaultPagination, untrustedInput,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	readDto := dto.ReadScheduledTasksRequest{
		Pagination:       requestPagination,
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
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, scheduledTasksList)
}

func (liaison *ScheduledTaskLiaison) Update(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	if untrustedInput["id"] != nil {
		untrustedInput["taskId"] = untrustedInput["id"]
	}

	requiredParams := []string{"taskId"}

	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	taskId, err := valueObject.NewScheduledTaskId(untrustedInput["taskId"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if untrustedInput["status"] != nil {
		taskStatus, err := valueObject.NewScheduledTaskStatus(untrustedInput["status"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	var runAtPtr *tkValueObject.UnixTime
	if untrustedInput["runAt"] != nil {
		runAt, err := tkValueObject.NewUnixTime(untrustedInput["runAt"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, "ScheduledTaskUpdated")
}
