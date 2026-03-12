package liaison

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"errors"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	sslInfra "github.com/goinfinite/os/src/infra/ssl"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
)

type SslLiaison struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	sslQueryRepo          *sslInfra.SslQueryRepo
	sslCmdRepo            *sslInfra.SslCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewSslLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SslLiaison {
	return &SslLiaison{
		persistentDbSvc:       persistentDbSvc,
		sslQueryRepo:          sslInfra.NewSslQueryRepo(),
		sslCmdRepo:            sslInfra.NewSslCmdRepo(persistentDbSvc, transientDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (liaison *SslLiaison) SslPairReadRequestFactory(
	untrustedInput map[string]any,
	withMappings bool,
) (readRequestDto dto.ReadSslPairsRequest, err error) {
	if untrustedInput["sslPairId"] == nil && untrustedInput["id"] != nil {
		untrustedInput["sslPairId"] = untrustedInput["id"]
	}

	var sslPairIdPtr *valueObject.SslPairId
	if untrustedInput["sslPairId"] != nil {
		sslPairId, err := valueObject.NewSslPairId(untrustedInput["sslPairId"])
		if err != nil {
			return readRequestDto, err
		}
		sslPairIdPtr = &sslPairId
	}

	if untrustedInput["virtualHostHostname"] == nil && untrustedInput["hostname"] != nil {
		untrustedInput["virtualHostHostname"] = untrustedInput["hostname"]
	}

	var vhostHostnamePtr *tkValueObject.Fqdn
	if untrustedInput["virtualHostHostname"] != nil {
		vhostHostname, err := tkValueObject.NewFqdn(untrustedInput["virtualHostHostname"])
		if err != nil {
			return readRequestDto, err
		}
		vhostHostnamePtr = &vhostHostname
	}

	altNames := []valueObject.SslHostname{}
	if untrustedInput["altNames"] != nil {
		var assertOk bool
		altNames, assertOk = untrustedInput["altNames"].([]valueObject.SslHostname)
		if !assertOk {
			return readRequestDto, errors.New("InvalidAltNamesStructure")
		}
	}

	timeParamNames := []string{
		"issuedBeforeAt", "issuedAfterAt", "expiresBeforeAt", "expiresAfterAt",
	}
	timeParamPtrs := tkPresentation.TimeParamsParser(timeParamNames, untrustedInput)

	requestPagination, err := tkPresentation.PaginationParser(
		useCase.SslPairsDefaultPagination, untrustedInput,
	)
	if err != nil {
		return readRequestDto, err
	}

	return dto.ReadSslPairsRequest{
		Pagination:          requestPagination,
		SslPairId:           sslPairIdPtr,
		VirtualHostHostname: vhostHostnamePtr,
		AltNames:            altNames,
		IssuedBeforeAt:      timeParamPtrs["issuedBeforeAt"],
		IssuedAfterAt:       timeParamPtrs["issuedAfterAt"],
		ExpiresBeforeAt:     timeParamPtrs["expiresBeforeAt"],
		ExpiresAfterAt:      timeParamPtrs["expiresAfterAt"],
	}, nil
}

func (liaison *SslLiaison) Read(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	readRequestDto, err := liaison.SslPairReadRequestFactory(untrustedInput, false)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	readResponseDto, err := useCase.ReadSslPairs(liaison.sslQueryRepo, readRequestDto)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, readResponseDto)
}

func (liaison *SslLiaison) Create(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"virtualHostsHostnames", "certificate", "key"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	vhostHostnames, assertOk := untrustedInput["virtualHostsHostnames"].([]tkValueObject.Fqdn)
	if !assertOk {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, errors.New("InvalidVirtualHostsStructure"))
	}

	certContent, err := valueObject.NewSslCertificateContent(untrustedInput["certificate"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}
	certEntity, err := entity.NewSslCertificate(certContent)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	var chainCertsPtr *entity.SslCertificate
	if untrustedInput["chainCertificates"] != nil {
		chainCertContent, err := valueObject.NewSslCertificateContent(untrustedInput["chainCertificates"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, errors.New("SslCertificateChainContentError"))
		}
		chainCertEntity, err := entity.NewSslCertificate(chainCertContent)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, errors.New("SslCertificateChainParseError"))
		}
		chainCertsPtr = &chainCertEntity
	}

	privateKeyContent, err := valueObject.NewSslPrivateKey(untrustedInput["key"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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

	createDto := dto.NewCreateSslPair(
		vhostHostnames, certEntity, chainCertsPtr, privateKeyContent,
		operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(liaison.persistentDbSvc)

	err = useCase.CreateSslPair(
		vhostQueryRepo, liaison.sslCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusCreated, "SslPairCreated")
}

func (liaison *SslLiaison) CreatePubliclyTrusted(
	untrustedInput map[string]any,
	shouldSchedule bool,
) tkPresentation.LiaisonResponse {
	if untrustedInput["hostname"] != nil && untrustedInput["virtualHostHostname"] == nil {
		untrustedInput["virtualHostHostname"] = untrustedInput["hostname"]
	}

	if untrustedInput["vhostHostname"] != nil && untrustedInput["virtualHostHostname"] == nil {
		untrustedInput["virtualHostHostname"] = untrustedInput["vhostHostname"]
	}

	requiredParams := []string{"virtualHostHostname"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	vhostHostname, err := tkValueObject.NewFqdn(untrustedInput["virtualHostHostname"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteOsBinary + " ssl create-trusted"
		installParams := []string{
			"--hostname", vhostHostname.String(),
		}
		cliCmd += " " + strings.Join(installParams, " ")

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(liaison.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("CreatePubliclyTrustedSslPair")
		taskCmd, _ := tkValueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("ssl")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSecs := uint16(1800)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
		}

		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusCreated, "PubliclyTrustedSslPairCreationScheduled")
	}

	createDto := dto.NewCreatePubliclyTrustedSslPair(
		vhostHostname, operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(liaison.persistentDbSvc)

	_, err = useCase.CreatePubliclyTrustedSslPair(
		vhostQueryRepo, liaison.sslCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusCreated, "PubliclyTrustedSslPairCreated")
}

func (liaison *SslLiaison) Delete(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	if untrustedInput["id"] == nil && untrustedInput["sslPairId"] != nil {
		untrustedInput["id"] = untrustedInput["sslPairId"]
	}

	requiredParams := []string{"id"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	pairId, err := valueObject.NewSslPairId(untrustedInput["id"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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

	err = useCase.DeleteSslPair(
		liaison.sslQueryRepo, liaison.sslCmdRepo, liaison.activityRecordCmdRepo,
		dto.NewDeleteSslPair(pairId, operatorAccountId, operatorIpAddress),
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, "SslPairDeleted")
}
