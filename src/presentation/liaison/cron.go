package liaison

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	cronInfra "github.com/goinfinite/os/src/infra/cron"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
)

type CronLiaison struct {
	cronQueryRepo         *cronInfra.CronQueryRepo
	cronCmdRepo           *cronInfra.CronCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewCronLiaison(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *CronLiaison {
	return &CronLiaison{
		cronQueryRepo:         cronInfra.NewCronQueryRepo(),
		cronCmdRepo:           cronInfra.NewCronCmdRepo(),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (liaison *CronLiaison) Read(untrustedInput map[string]any) LiaisonOutput {
	var idPtr *valueObject.CronId
	if untrustedInput["id"] != nil {
		id, err := valueObject.NewCronId(untrustedInput["id"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		idPtr = &id
	}

	var commentPtr *valueObject.CronComment
	if untrustedInput["comment"] != nil {
		slug, err := valueObject.NewCronComment(untrustedInput["comment"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		commentPtr = &slug
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if untrustedInput["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(untrustedInput["pageNumber"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if untrustedInput["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(untrustedInput["itemsPerPage"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if untrustedInput["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(untrustedInput["sortBy"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if untrustedInput["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			untrustedInput["sortDirection"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if untrustedInput["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(untrustedInput["lastSeenId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadCronsRequest{
		Pagination:  paginationDto,
		CronId:      idPtr,
		CronComment: commentPtr,
	}

	cronsList, err := useCase.ReadCrons(liaison.cronQueryRepo, readDto)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, cronsList)
}

func (liaison *CronLiaison) Create(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"schedule", "command"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	schedule, err := valueObject.NewCronSchedule(untrustedInput["schedule"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	command, err := valueObject.NewUnixCommand(untrustedInput["command"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	var commentPtr *valueObject.CronComment
	if untrustedInput["comment"] != nil && untrustedInput["comment"] != "" {
		comment, err := valueObject.NewCronComment(untrustedInput["comment"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		commentPtr = &comment
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateCron(
		schedule, command, commentPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateCron(
		liaison.cronCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "CronCreated")
}

func (liaison *CronLiaison) Update(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"id"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	id, err := valueObject.NewCronId(untrustedInput["id"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	var schedulePtr *valueObject.CronSchedule
	if untrustedInput["schedule"] != nil {
		schedule, err := valueObject.NewCronSchedule(untrustedInput["schedule"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		schedulePtr = &schedule
	}

	var commandPtr *valueObject.UnixCommand
	if untrustedInput["command"] != nil {
		command, err := valueObject.NewUnixCommand(untrustedInput["command"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		commandPtr = &command
	}

	clearableFields := []string{}

	var commentPtr *valueObject.CronComment
	switch commentValue := untrustedInput["comment"]; {
	case commentValue == nil:
	case commentValue == "" || commentValue == " ":
		clearableFields = append(clearableFields, "comment")
	default:
		comment, err := valueObject.NewCronComment(commentValue)
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		commentPtr = &comment
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdateCron(
		id, schedulePtr, commandPtr, commentPtr, clearableFields,
		operatorAccountId, operatorIpAddress,
	)

	err = useCase.UpdateCron(
		liaison.cronQueryRepo, liaison.cronCmdRepo, liaison.activityRecordCmdRepo,
		updateDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "CronUpdated")
}

func (liaison *CronLiaison) Delete(untrustedInput map[string]any) LiaisonOutput {
	var idPtr *valueObject.CronId
	if untrustedInput["cronId"] != nil {
		id, err := valueObject.NewCronId(untrustedInput["cronId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		idPtr = &id
	}

	var commentPtr *valueObject.CronComment
	if untrustedInput["comment"] != nil {
		comment, err := valueObject.NewCronComment(untrustedInput["comment"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		commentPtr = &comment
	}

	var err error

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteCron(
		idPtr, commentPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteCron(
		liaison.cronQueryRepo, liaison.cronCmdRepo, liaison.activityRecordCmdRepo,
		deleteDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "CronDeleted")
}
