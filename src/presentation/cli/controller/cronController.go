package cliController

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetCronsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetCrons",
		Run: func(cmd *cobra.Command, args []string) {
			cronQueryRepo := infra.CronQueryRepo{}

			cronsList, err := useCase.GetCrons(cronQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, cronsList)
		},
	}
	return cmd
}

func AddCronControler() *cobra.Command {
	var scheduleStr string
	var commandStr string
	var commentStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewCron",
		Run: func(cmd *cobra.Command, args []string) {
			var commentPtr *valueObject.CronComment
			if commentStr != "" {
				comment := valueObject.NewCronCommentPanic(commentStr)
				commentPtr = &comment
			}

			addCronDto := dto.NewAddCron(
				valueObject.NewCronSchedulePanic(scheduleStr),
				valueObject.NewUnixCommandPanic(commandStr),
				commentPtr,
			)

			cronCmdRepo := infra.CronCmdRepo{}

			err := useCase.AddCron(
				cronCmdRepo,
				addCronDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "CronAdded")
		},
	}

	cmd.Flags().StringVarP(&scheduleStr, "schedule", "s", "", "Schedule")
	cmd.MarkFlagRequired("schedule")
	cmd.Flags().StringVarP(&commandStr, "command", "c", "", "Command")
	cmd.MarkFlagRequired("command")
	cmd.Flags().StringVarP(&commentStr, "comment", "d", "", "Comment")
	return cmd
}
