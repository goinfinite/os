package service

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	cronInfra "github.com/goinfinite/os/src/infra/cron"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type CronService struct {
	cronQueryRepo         *cronInfra.CronQueryRepo
	cronCmdRepo           *cronInfra.CronCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewCronService(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *CronService {
	return &CronService{
		cronQueryRepo:         cronInfra.NewCronQueryRepo(),
		cronCmdRepo:           cronInfra.NewCronCmdRepo(),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *CronService) Read(input map[string]interface{}) ServiceOutput {
	var idPtr *valueObject.CronId
	if input["id"] != nil {
		id, err := valueObject.NewCronId(input["id"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		idPtr = &id
	}

	var commentPtr *valueObject.CronComment
	if input["comment"] != nil {
		slug, err := valueObject.NewCronComment(input["comment"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		commentPtr = &slug
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			input["sortDirection"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadCronsRequest{
		Pagination:  paginationDto,
		CronId:      idPtr,
		CronComment: commentPtr,
	}

	cronsList, err := useCase.ReadCrons(service.cronQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, cronsList)
}

func (service *CronService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"schedule", "command"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	schedule, err := valueObject.NewCronSchedule(input["schedule"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	command, err := valueObject.NewUnixCommand(input["command"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var commentPtr *valueObject.CronComment
	if input["comment"] != nil && input["comment"] != "" {
		comment, err := valueObject.NewCronComment(input["comment"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		commentPtr = &comment
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateCron(
		schedule, command, commentPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateCron(
		service.cronCmdRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "CronCreated")
}

func (service *CronService) Update(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"id"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	id, err := valueObject.NewCronId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var schedulePtr *valueObject.CronSchedule
	if input["schedule"] != nil {
		schedule, err := valueObject.NewCronSchedule(input["schedule"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		schedulePtr = &schedule
	}

	var commandPtr *valueObject.UnixCommand
	if input["command"] != nil {
		command, err := valueObject.NewUnixCommand(input["command"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		commandPtr = &command
	}

	var commentPtr *valueObject.CronComment
	if input["comment"] != nil {
		comment, err := valueObject.NewCronComment(input["comment"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		commentPtr = &comment
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateCron(
		id, schedulePtr, commandPtr, commentPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateCron(
		service.cronQueryRepo, service.cronCmdRepo, service.activityRecordCmdRepo,
		updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "CronUpdated")
}

func (service *CronService) Delete(input map[string]interface{}) ServiceOutput {
	var idPtr *valueObject.CronId
	if input["cronId"] != nil {
		id, err := valueObject.NewCronId(input["cronId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		idPtr = &id
	}

	var commentPtr *valueObject.CronComment
	if input["comment"] != nil {
		comment, err := valueObject.NewCronComment(input["comment"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		commentPtr = &comment
	}

	var err error

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteCron(
		idPtr, commentPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteCron(
		service.cronQueryRepo, service.cronCmdRepo, service.activityRecordCmdRepo,
		deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "CronDeleted")
}
