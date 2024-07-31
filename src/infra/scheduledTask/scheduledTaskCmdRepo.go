package scheduledTaskInfra

import (
	"strconv"
	"time"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type ScheduledTaskCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewScheduledTaskCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskCmdRepo {
	return &ScheduledTaskCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ScheduledTaskCmdRepo) Create(
	createDto dto.CreateScheduledTask,
) error {
	newTaskStatus, _ := valueObject.NewScheduledTaskStatus("pending")

	var runAtPtr *time.Time
	if createDto.RunAt != nil {
		runAt := time.Unix(createDto.RunAt.Int64(), 0)
		runAtPtr = &runAt
	}

	scheduledTaskModel := dbModel.NewScheduledTask(
		0, createDto.Name.String(), newTaskStatus.String(), createDto.Command.String(),
		createDto.Tags, createDto.TimeoutSecs, runAtPtr, nil, nil,
	)

	return repo.persistentDbSvc.Handler.Create(&scheduledTaskModel).Error
}

func (repo *ScheduledTaskCmdRepo) Update(
	updateDto dto.UpdateScheduledTask,
) error {
	updateMap := map[string]interface{}{}

	if updateDto.Status != nil {
		updateMap["status"] = updateDto.Status.String()
	}

	if updateDto.RunAt != nil {
		updateMap["run_at"] = updateDto.RunAt.GetAsGoTime()
	}

	if len(updateMap) == 0 {
		return nil
	}

	return repo.persistentDbSvc.Handler.
		Model(&dbModel.ScheduledTask{}).
		Where("id = ?", updateDto.Id).
		Updates(updateMap).Error
}

func (repo *ScheduledTaskCmdRepo) Run(
	pendingTask entity.ScheduledTask,
) error {
	runningStatus, _ := valueObject.NewScheduledTaskStatus("running")
	updateDto := dto.NewUpdateScheduledTask(pendingTask.Id, &runningStatus, nil)
	err := repo.Update(updateDto)
	if err != nil {
		return err
	}

	timeoutSecs := useCase.ScheduledTasksDefaultTimeoutSecs
	if pendingTask.TimeoutSecs != nil {
		timeoutSecs = *pendingTask.TimeoutSecs
	}
	timeoutStr := strconv.FormatUint(uint64(timeoutSecs), 10)

	cmdWithTimeout := "timeout --kill-after=10s " + timeoutStr + " " + pendingTask.Command.String()
	rawOutput, rawError := infraHelper.RunCmdWithSubShell(cmdWithTimeout)

	finalStatus, _ := valueObject.NewScheduledTaskStatus("completed")
	if rawError != nil {
		finalStatus, _ = valueObject.NewScheduledTaskStatus("failed")
	}

	updateMap := map[string]interface{}{
		"status": finalStatus.String(),
	}

	if len(rawOutput) > 0 {
		taskOutput, err := valueObject.NewScheduledTaskOutput(rawOutput)
		if err != nil {
			return err
		}
		updateMap["output"] = taskOutput.String()
	}

	if rawError != nil {
		taskError, err := valueObject.NewScheduledTaskOutput(rawError.Error())
		if err != nil {
			return err
		}
		updateMap["error"] = taskError.String()
	}

	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.ScheduledTask{}).
		Where("id = ?", pendingTask.Id).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *ScheduledTaskCmdRepo) Delete(id valueObject.ScheduledTaskId) error {
	return repo.persistentDbSvc.Handler.
		Where("id = ?", id).
		Delete(&dbModel.ScheduledTask{}).Error
}
