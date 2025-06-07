package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/spf13/cobra"
)

type ScheduledTaskController struct {
	persistentDbSvc      *internalDbInfra.PersistentDatabaseService
	scheduledTaskLiaison *liaison.ScheduledTaskLiaison
}

func NewScheduledTaskController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ScheduledTaskController {
	return &ScheduledTaskController{
		persistentDbSvc:      persistentDbSvc,
		scheduledTaskLiaison: liaison.NewScheduledTaskLiaison(persistentDbSvc),
	}
}

func (controller *ScheduledTaskController) Read() *cobra.Command {
	var taskIdUint uint64
	var taskNameStr, taskStatusStr string
	var taskTagsStrSlice []string
	var startedBeforeAtInt64, startedAfterAtInt64 int64
	var finishedBeforeAtInt64, finishedAfterAtInt64 int64
	var createdBeforeAtInt64, createdAfterAtInt64 int64
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr string
	var paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadScheduledTasks",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if taskIdUint != 0 {
				requestBody["taskId"] = taskIdUint
			}
			if taskNameStr != "" {
				requestBody["taskName"] = taskNameStr
			}
			if taskStatusStr != "" {
				requestBody["taskStatus"] = taskStatusStr
			}
			if len(taskTagsStrSlice) > 0 {
				requestBody["taskTags"] = taskTagsStrSlice
			}
			if startedBeforeAtInt64 != 0 {
				requestBody["startedBeforeAt"] = startedBeforeAtInt64
			}
			if startedAfterAtInt64 != 0 {
				requestBody["startedAfterAt"] = startedAfterAtInt64
			}
			if finishedBeforeAtInt64 != 0 {
				requestBody["finishedBeforeAt"] = finishedBeforeAtInt64
			}
			if finishedAfterAtInt64 != 0 {
				requestBody["finishedAfterAt"] = finishedAfterAtInt64
			}
			if createdBeforeAtInt64 != 0 {
				requestBody["createdBeforeAt"] = createdBeforeAtInt64
			}
			if createdAfterAtInt64 != 0 {
				requestBody["createdAfterAt"] = createdAfterAtInt64
			}
			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}
			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}
			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}
			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}
			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.LiaisonResponseWrapper(
				controller.scheduledTaskLiaison.Read(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&taskIdUint, "task-id", "i", 0, "TaskId")
	cmd.Flags().StringVarP(&taskNameStr, "task-name", "n", "", "TaskName")
	cmd.Flags().StringVarP(&taskStatusStr, "task-status", "s", "", "TaskStatus")
	cmd.Flags().StringSliceVarP(&taskTagsStrSlice, "task-tags", "t", []string{}, "TaskTags")
	cmd.Flags().Int64VarP(
		&startedBeforeAtInt64, "started-before-at", "b", 0, "StartedBeforeAt (UnixTime)",
	)
	cmd.Flags().Int64VarP(
		&startedAfterAtInt64, "started-after-at", "a", 0, "StartedAfterAt (UnixTime)",
	)
	cmd.Flags().Int64VarP(
		&finishedBeforeAtInt64, "finished-before-at", "f", 0, "FinishedBeforeAt (UnixTime)",
	)
	cmd.Flags().Int64VarP(
		&finishedAfterAtInt64, "finished-after-at", "e", 0, "FinishedAfterAt (UnixTime)",
	)
	cmd.Flags().Int64VarP(
		&createdBeforeAtInt64, "created-before-at", "c", 0, "CreatedBeforeAt (UnixTime)",
	)
	cmd.Flags().Int64VarP(
		&createdAfterAtInt64, "created-after-at", "d", 0, "CreatedAfterAt (UnixTime)",
	)
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "p", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "m", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "r", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *ScheduledTaskController) Update() *cobra.Command {
	var taskIdUint64 uint64
	var statusStr string
	var runAtInt64 int64

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateScheduledTask",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"taskId": taskIdUint64,
			}

			if statusStr != "" {
				requestBody["status"] = statusStr
			}
			if runAtInt64 != 0 {
				requestBody["runAt"] = runAtInt64
			}

			cliHelper.LiaisonResponseWrapper(
				controller.scheduledTaskLiaison.Update(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&taskIdUint64, "task-id", "i", 0, "TaskId")
	cmd.MarkFlagRequired("task-id")
	cmd.Flags().StringVarP(&statusStr, "status", "s", "", "Status (pending/cancelled)")
	cmd.Flags().Int64VarP(&runAtInt64, "run-at", "r", 0, "RunAt (UnixTime)")
	return cmd
}
