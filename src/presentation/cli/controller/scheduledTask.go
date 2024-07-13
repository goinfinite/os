package cliController

import (
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type ScheduledTaskController struct {
	persistentDbSvc      *internalDbInfra.PersistentDatabaseService
	scheduledTaskService *service.ScheduledTaskService
}

func NewScheduledTaskController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskController {
	return &ScheduledTaskController{
		persistentDbSvc:      persistentDbSvc,
		scheduledTaskService: service.NewScheduledTaskService(persistentDbSvc),
	}
}

func (controller *ScheduledTaskController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadScheduledTasks",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.scheduledTaskService.Read())
		},
	}

	return cmd
}

func (controller *ScheduledTaskController) Update() *cobra.Command {
	var taskIdUint uint
	var statusStr string
	var runAtInt64 int64

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateScheduledTask",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": taskIdUint,
			}

			if statusStr != "" {
				requestBody["status"] = statusStr
			}

			if runAtInt64 != 0 {
				requestBody["runAt"] = runAtInt64
			}

			cliHelper.ServiceResponseWrapper(
				controller.scheduledTaskService.Update(requestBody),
			)
		},
	}

	cmd.Flags().UintVarP(&taskIdUint, "task-id", "t", 0, "TaskId")
	cmd.MarkFlagRequired("task-id")
	cmd.Flags().StringVarP(&statusStr, "status", "s", "", "Status (pending/cancelled)")
	cmd.Flags().Int64VarP(&runAtInt64, "run-at", "r", 0, "RunAt (in unix epoch)")
	return cmd
}
