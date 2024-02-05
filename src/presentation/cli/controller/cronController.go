package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	cronInfra "github.com/speedianet/os/src/infra/cron"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetCronsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetCrons",
		Run: func(cmd *cobra.Command, args []string) {
			cronQueryRepo := cronInfra.CronQueryRepo{}

			cronsList, err := useCase.GetCrons(cronQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, cronsList)
		},
	}
	return cmd
}

func CreateCronControler() *cobra.Command {
	var scheduleStr string
	var commandStr string
	var commentStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "CreateNewCron",
		Run: func(cmd *cobra.Command, args []string) {
			var commentPtr *valueObject.CronComment
			if commentStr != "" {
				comment := valueObject.NewCronCommentPanic(commentStr)
				commentPtr = &comment
			}

			createCronDto := dto.NewCreateCron(
				valueObject.NewCronSchedulePanic(scheduleStr),
				valueObject.NewUnixCommandPanic(commandStr),
				commentPtr,
			)

			cronCmdRepo, err := cronInfra.NewCronCmdRepo()
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			err = useCase.CreateCron(
				cronCmdRepo,
				createCronDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "CronCreated")
		},
	}

	cmd.Flags().StringVarP(&scheduleStr, "schedule", "s", "", "Schedule")
	cmd.MarkFlagRequired("schedule")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command")
	cmd.MarkFlagRequired("command")
	cmd.Flags().StringVarP(&commentStr, "comment", "d", "", "Comment")
	return cmd
}

func UpdateCronController() *cobra.Command {
	var idStr string
	var scheduleStr string
	var commandStr string
	var commentStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateCron",
		Run: func(cmd *cobra.Command, args []string) {
			var schedulePtr *valueObject.CronSchedule
			if scheduleStr != "" {
				schedule := valueObject.NewCronSchedulePanic(scheduleStr)
				schedulePtr = &schedule
			}

			var commandPtr *valueObject.UnixCommand
			if commandStr != "" {
				command := valueObject.NewUnixCommandPanic(commandStr)
				commandPtr = &command
			}

			var commentPtr *valueObject.CronComment
			if commentStr != "" {
				comment := valueObject.NewCronCommentPanic(commentStr)
				commentPtr = &comment
			}

			updateCronDto := dto.NewUpdateCron(
				valueObject.NewCronIdPanic(idStr),
				schedulePtr,
				commandPtr,
				commentPtr,
			)

			cronQueryRepo := cronInfra.CronQueryRepo{}
			cronCmdRepo, err := cronInfra.NewCronCmdRepo()
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			err = useCase.UpdateCron(
				cronQueryRepo,
				cronCmdRepo,
				updateCronDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "CronUpdated")
		},
	}

	cmd.Flags().StringVarP(&idStr, "id", "i", "", "CronId")
	cmd.MarkFlagRequired("id")
	cmd.Flags().StringVarP(&scheduleStr, "schedule", "s", "", "Schedule")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command")
	cmd.Flags().StringVarP(&commentStr, "comment", "d", "", "Comment")
	return cmd
}

func DeleteCronController() *cobra.Command {
	var cronIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteCron",
		Run: func(cmd *cobra.Command, args []string) {
			cronId := valueObject.NewCronIdPanic(cronIdStr)

			cronQueryRepo := cronInfra.CronQueryRepo{}
			cronCmdRepo, err := cronInfra.NewCronCmdRepo()
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			err = useCase.DeleteCron(
				cronQueryRepo,
				cronCmdRepo,
				cronId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "CronDeleted")
		},
	}

	cmd.Flags().StringVarP(&cronIdStr, "id", "i", "", "CronId")
	cmd.MarkFlagRequired("id")
	return cmd
}
