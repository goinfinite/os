package apiController

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/labstack/echo/v4"
)

type ScheduledTaskController struct {
	scheduledTaskLiaison *liaison.ScheduledTaskLiaison
	persistentDbSvc      *internalDbInfra.PersistentDatabaseService
}

func NewScheduledTaskController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskController {
	return &ScheduledTaskController{
		scheduledTaskLiaison: liaison.NewScheduledTaskLiaison(persistentDbSvc),
		persistentDbSvc:      persistentDbSvc,
	}
}

func (controller *ScheduledTaskController) parseTaskTags(
	rawTags interface{},
) ([]valueObject.ScheduledTaskTag, error) {
	taskTags := []valueObject.ScheduledTaskTag{}

	rawTagsStr, assertOk := rawTags.(string)
	if !assertOk {
		return taskTags, errors.New("InvalidTaskTagsStructure")
	}

	for rawTagIndex, rawTag := range strings.Split(rawTagsStr, ";") {
		taskTag, err := valueObject.NewScheduledTaskTag(rawTag)
		if err != nil {
			slog.Debug("InvalidTaskTag", slog.Int("tagIndex", rawTagIndex))
			continue
		}
		taskTags = append(taskTags, taskTag)
	}

	return taskTags, nil
}

// ReadScheduledTasks	 godoc
// @Summary      ReadScheduledTasks
// @Description  List scheduled tasks.
// @Tags         scheduled-task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        taskId query  string  false  "TaskId"
// @Param        taskName query  string  false  "TaskName"
// @Param        taskStatus query  string  false  "TaskStatus"
// @Param        taskTags query  string  false  "TaskTags (semicolon separated)"
// @Param        startedBeforeAt query  string  false  "StartedBeforeAt"
// @Param        startedAfterAt query  string  false  "StartedAfterAt"
// @Param        finishedBeforeAt query  string  false  "FinishedBeforeAt"
// @Param        finishedAfterAt query  string  false  "FinishedAfterAt"
// @Param        createdBeforeAt query  string  false  "CreatedBeforeAt"
// @Param        createdAfterAt query  string  false  "CreatedAfterAt"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadScheduledTasksResponse
// @Router       /v1/scheduled-task/ [get]
func (controller *ScheduledTaskController) Read(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if _, exists := requestInputData["taskTags"]; exists {
		taskTags, err := controller.parseTaskTags(requestInputData["taskTags"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		requestInputData["taskTags"] = taskTags
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.scheduledTaskLiaison.Read(requestInputData),
	)
}

// UpdateScheduledTask godoc
// @Summary      UpdateScheduledTask
// @Description  Reschedule a task or change its status.
// @Tags         scheduled-task
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateScheduledTaskDto 	  body dto.UpdateScheduledTask  true  "UpdateScheduledTask (Only id is required.)"
// @Success      200 {object} object{} "ScheduledTaskUpdated"
// @Router       /v1/scheduled-task/ [put]
func (controller *ScheduledTaskController) Update(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.scheduledTaskLiaison.Update(requestInputData),
	)
}

func (controller *ScheduledTaskController) Run() {
	timer := time.NewTicker(
		time.Duration(int64(useCase.ScheduledTasksRunIntervalSecs)) * time.Second,
	)
	defer timer.Stop()

	scheduledTaskQueryRepo := scheduledTaskInfra.NewScheduledTaskQueryRepo(
		controller.persistentDbSvc,
	)
	scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(
		controller.persistentDbSvc,
	)
	for range timer.C {
		go useCase.RunScheduledTasks(scheduledTaskQueryRepo, scheduledTaskCmdRepo)
	}
}
