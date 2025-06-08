package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/spf13/cobra"
)

type CronController struct {
	cronLiaison *liaison.CronLiaison
}

func NewCronController(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *CronController {
	return &CronController{
		cronLiaison: liaison.NewCronLiaison(trailDbSvc),
	}
}

func (controller *CronController) Read() *cobra.Command {
	var idUint uint64
	var commentStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadCrons",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if idUint != 0 {
				requestBody["id"] = idUint
			}

			if commentStr != "" {
				requestBody["comment"] = commentStr
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

			cliHelper.LiaisonResponseWrapper(controller.cronLiaison.Read(requestBody))
		},
	}

	cmd.Flags().Uint64VarP(&idUint, "cron-id", "i", 0, "CronId")
	cmd.Flags().StringVarP(&commentStr, "cron-comment", "s", "", "CronComment")
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "p", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "m", 0,
		"ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "r", "",
		"SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *CronController) Create() *cobra.Command {
	var scheduleStr, commandStr, commentStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewCron",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"schedule": scheduleStr,
				"command":  commandStr,
			}

			if commentStr != "" {
				requestBody["comment"] = commentStr
			}

			cliHelper.LiaisonResponseWrapper(
				controller.cronLiaison.Create(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&scheduleStr, "schedule", "s", "", "Schedule")
	cmd.MarkFlagRequired("schedule")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command")
	cmd.MarkFlagRequired("command")
	cmd.Flags().StringVarP(&commentStr, "comment", "d", "", "Comment")
	return cmd
}

func (controller *CronController) Update() *cobra.Command {
	var idStr, scheduleStr, commandStr, commentStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateCron",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": idStr,
			}

			if scheduleStr != "" {
				requestBody["schedule"] = scheduleStr
			}

			if commandStr != "" {
				requestBody["command"] = commandStr
			}

			if commentStr != "" {
				requestBody["comment"] = commentStr
			}

			cliHelper.LiaisonResponseWrapper(
				controller.cronLiaison.Update(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&idStr, "id", "i", "", "CronId")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&scheduleStr, "schedule", "s", "", "Schedule")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command")
	cmd.Flags().StringVarP(&commentStr, "comment", "d", "", "Comment")
	return cmd
}

func (controller *CronController) Delete() *cobra.Command {
	var idStr, commentStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteCron",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if idStr != "" {
				requestBody["id"] = idStr
			}

			if commentStr != "" {
				requestBody["comment"] = commentStr
			}

			cliHelper.LiaisonResponseWrapper(
				controller.cronLiaison.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&idStr, "id", "i", "", "CronId")
	cmd.Flags().StringVarP(&commentStr, "comment", "d", "", "CronComment")
	return cmd
}
