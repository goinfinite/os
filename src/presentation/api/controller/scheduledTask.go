package apiController

import (
	"time"

	"github.com/goinfinite/os/src/domain/useCase"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type ScheduledTaskController struct {
	scheduledTaskService *service.ScheduledTaskService
	persistentDbSvc      *internalDbInfra.PersistentDatabaseService
}

func NewScheduledTaskController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskController {
	return &ScheduledTaskController{
		scheduledTaskService: service.NewScheduledTaskService(persistentDbSvc),
		persistentDbSvc:      persistentDbSvc,
	}
}

// ReadScheduledTasks	 godoc
// @Summary      ReadScheduledTasks
// @Description  List scheduled tasks.
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {array} entity.ScheduledTask
// @Router       /v1/task/ [get]
func (controller *ScheduledTaskController) Read(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(c, controller.scheduledTaskService.Read())
}

// UpdateScheduledTask godoc
// @Summary      UpdateScheduledTask
// @Description  Reschedule a task or change its status.
// @Tags         task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateScheduledTaskDto 	  body dto.UpdateScheduledTask  true  "UpdateScheduledTask (Only id is required.)"
// @Success      200 {object} object{} "ScheduledTaskUpdated"
// @Router       /v1/task/ [put]
func (controller *ScheduledTaskController) Update(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.scheduledTaskService.Update(requestBody),
	)
}

func (controller *ScheduledTaskController) Run() {
	timer := time.NewTicker(
		time.Duration(int64(useCase.ScheduledTasksRunIntervalSecs)) * time.Second,
	)
	defer timer.Stop()

	scheduledTaskQueryRepo := scheduledTaskInfra.NewScheduledTaskQueryRepo(controller.persistentDbSvc)
	scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(controller.persistentDbSvc)
	for range timer.C {
		go useCase.RunScheduledTasks(scheduledTaskQueryRepo, scheduledTaskCmdRepo)
	}
}
