package liaison

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	cronInfra "github.com/goinfinite/os/src/infra/cron"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
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

func (liaison *CronLiaison) Read(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	var idPtr *valueObject.CronId
	if untrustedInput["id"] != nil {
		id, err := valueObject.NewCronId(untrustedInput["id"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err)
		}
		idPtr = &id
	}

	var commentPtr *valueObject.CronComment
	if untrustedInput["comment"] != nil {
		slug, err := valueObject.NewCronComment(untrustedInput["comment"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err)
		}
		commentPtr = &slug
	}

	requestPagination, err := tkPresentation.PaginationParser(
		useCase.CronsDefaultPagination, untrustedInput,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	readDto := dto.ReadCronsRequest{
		Pagination:  requestPagination,
		CronId:      idPtr,
		CronComment: commentPtr,
	}

	cronsList, err := useCase.ReadCrons(liaison.cronQueryRepo, readDto)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, cronsList)
}

func (liaison *CronLiaison) Create(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"schedule", "command"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	schedule, err := valueObject.NewCronSchedule(untrustedInput["schedule"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	command, err := tkValueObject.NewUnixCommand(untrustedInput["command"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	var commentPtr *valueObject.CronComment
	if untrustedInput["comment"] != nil && untrustedInput["comment"] != "" {
		comment, err := valueObject.NewCronComment(untrustedInput["comment"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
		commentPtr = &comment
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	createDto := dto.NewCreateCron(
		schedule, command, commentPtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateCron(
		liaison.cronCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusCreated, "CronCreated")
}

func (liaison *CronLiaison) Update(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"id"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	id, err := valueObject.NewCronId(untrustedInput["id"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	var schedulePtr *valueObject.CronSchedule
	if untrustedInput["schedule"] != nil {
		schedule, err := valueObject.NewCronSchedule(untrustedInput["schedule"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
		schedulePtr = &schedule
	}

	var commandPtr *tkValueObject.UnixCommand
	if untrustedInput["command"] != nil {
		command, err := tkValueObject.NewUnixCommand(untrustedInput["command"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
		commentPtr = &comment
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, "CronUpdated")
}

func (liaison *CronLiaison) Delete(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	var idPtr *valueObject.CronId
	if untrustedInput["cronId"] != nil {
		id, err := valueObject.NewCronId(untrustedInput["cronId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
		idPtr = &id
	}

	var commentPtr *valueObject.CronComment
	if untrustedInput["comment"] != nil {
		comment, err := valueObject.NewCronComment(untrustedInput["comment"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
		commentPtr = &comment
	}

	var err error

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, "CronDeleted")
}
