package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	cronInfra "github.com/speedianet/os/src/infra/cron"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
)

type CronService struct {
}

func NewCronService() CronService {
	return CronService{}
}

func (service CronService) Read() ServiceOutput {
	cronQueryRepo := cronInfra.CronQueryRepo{}
	cronsList, err := useCase.GetCrons(cronQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, cronsList)
}

func (service CronService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"schedule", "command"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	cronSchedule, err := valueObject.NewCronSchedule(input["schedule"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	cronCommand, err := valueObject.NewUnixCommand(input["command"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var cronCommentPtr *valueObject.CronComment
	if input["comment"] != nil {
		cronComment, err := valueObject.NewCronComment(input["comment"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		cronCommentPtr = &cronComment
	}

	dto := dto.NewCreateCron(cronSchedule, cronCommand, cronCommentPtr)

	cronCmdRepo, err := cronInfra.NewCronCmdRepo()
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	err = useCase.CreateCron(cronCmdRepo, dto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "CronCreated")
}

func (service CronService) Update(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"id"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	cronId, err := valueObject.NewCronId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var cronSchedulePtr *valueObject.CronSchedule
	if input["schedule"] != nil {
		cronSchedule, err := valueObject.NewCronSchedule(input["schedule"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		cronSchedulePtr = &cronSchedule
	}

	var cronCommandPtr *valueObject.UnixCommand
	if input["command"] != nil {
		cronCommand, err := valueObject.NewUnixCommand(input["command"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		cronCommandPtr = &cronCommand
	}

	var cronCommentPtr *valueObject.CronComment
	if input["comment"] != nil {
		cronComment, err := valueObject.NewCronComment(input["comment"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		cronCommentPtr = &cronComment
	}

	updateCronDto := dto.NewUpdateCron(
		cronId, cronSchedulePtr, cronCommandPtr, cronCommentPtr,
	)

	cronQueryRepo := cronInfra.CronQueryRepo{}
	cronCmdRepo, err := cronInfra.NewCronCmdRepo()
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	err = useCase.UpdateCron(
		cronQueryRepo,
		cronCmdRepo,
		updateCronDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "CronUpdated")
}

func (service CronService) Delete(input map[string]interface{}) ServiceOutput {
	cronId, err := valueObject.NewCronId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	cronQueryRepo := cronInfra.CronQueryRepo{}
	cronCmdRepo, err := cronInfra.NewCronCmdRepo()
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	deleteDto := dto.NewDeleteCron(&cronId, nil)

	err = useCase.DeleteCron(cronQueryRepo, cronCmdRepo, deleteDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "CronDeleted")
}
