package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/speedianet/os/src/infra/scheduledTask"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
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

func (service *ScheduledTaskService) Read() ServiceOutput {
	scheduledTaskQueryRepo := scheduledTaskInfra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTasksList, err := useCase.ReadScheduledTasks(scheduledTaskQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, scheduledTasksList)
}

func (service *ScheduledTaskService) Update(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"id"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	taskId, err := valueObject.NewScheduledTaskId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var taskStatusPtr *valueObject.ScheduledTaskStatus
	if _, exists := input["status"]; exists {
		taskStatus, err := valueObject.NewScheduledTaskStatus(input["status"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		taskStatusPtr = &taskStatus
	}

	var runAtPtr *valueObject.UnixTime
	if _, exists := input["runAt"]; exists {
		runAt, err := valueObject.NewUnixTime(input["runAt"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		runAtPtr = &runAt
	}

	updateDto := dto.NewUpdateScheduledTask(
		taskId,
		taskStatusPtr,
		runAtPtr,
	)

	scheduledTaskQueryRepo := scheduledTaskInfra.NewScheduledTaskQueryRepo(service.persistentDbSvc)
	scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(service.persistentDbSvc)

	err = useCase.UpdateScheduledTask(
		scheduledTaskQueryRepo,
		scheduledTaskCmdRepo,
		updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "ScheduledTaskUpdated")
}
