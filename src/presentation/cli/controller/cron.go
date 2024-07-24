package cliController

import (
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type CronController struct {
	cronService *service.CronService
}

func NewCronController() *CronController {
	return &CronController{
		cronService: service.NewCronService(),
	}
}

func (controller *CronController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadCrons",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.cronService.Read())
		},
	}
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

			cliHelper.ServiceResponseWrapper(
				controller.cronService.Create(requestBody),
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

			cliHelper.ServiceResponseWrapper(
				controller.cronService.Update(requestBody),
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

			cliHelper.ServiceResponseWrapper(
				controller.cronService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&idStr, "id", "i", "", "CronId")
	cmd.Flags().StringVarP(&commentStr, "comment", "d", "", "CronComment")
	return cmd
}
