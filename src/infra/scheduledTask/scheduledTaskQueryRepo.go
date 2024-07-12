package scheduledTaskInfra

import (
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type ScheduledTaskQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewScheduledTaskQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskQueryRepo {
	return &ScheduledTaskQueryRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ScheduledTaskQueryRepo) Read() ([]entity.ScheduledTask, error) {
	scheduledTasks := []entity.ScheduledTask{}

	scheduledTaskModels := []dbModel.ScheduledTask{}
	err := repo.persistentDbSvc.Handler.
		Find(&scheduledTaskModels).Error
	if err != nil {
		return scheduledTasks, err
	}

	for _, scheduledTaskModel := range scheduledTaskModels {
		scheduledTaskEntity, err := scheduledTaskModel.ToEntity()
		if err != nil {
			log.Printf("[%d] %s", scheduledTaskModel.ID, err.Error())
			continue
		}
		scheduledTasks = append(scheduledTasks, scheduledTaskEntity)
	}

	return scheduledTasks, nil
}

func (repo *ScheduledTaskQueryRepo) ReadById(
	id valueObject.ScheduledTaskId,
) (taskEntity entity.ScheduledTask, err error) {
	var scheduledTaskModel dbModel.ScheduledTask
	err = repo.persistentDbSvc.Handler.
		Where("id = ?", id).
		First(&scheduledTaskModel).Error
	if err != nil {
		return taskEntity, err
	}

	return scheduledTaskModel.ToEntity()
}

func (repo *ScheduledTaskQueryRepo) ReadByStatus(
	status valueObject.ScheduledTaskStatus,
) ([]entity.ScheduledTask, error) {
	scheduledTasks := []entity.ScheduledTask{}

	scheduledTaskModels := []dbModel.ScheduledTask{}
	err := repo.persistentDbSvc.Handler.
		Where("status = ?", status.String()).
		Find(&scheduledTaskModels).Error
	if err != nil {
		return scheduledTasks, err
	}

	for _, scheduledTaskModel := range scheduledTaskModels {
		scheduledTaskEntity, err := scheduledTaskModel.ToEntity()
		if err != nil {
			log.Printf("[%d] ModelToEntityError: %s", scheduledTaskModel.ID, err.Error())
			continue
		}
		scheduledTasks = append(scheduledTasks, scheduledTaskEntity)
	}

	return scheduledTasks, nil
}
