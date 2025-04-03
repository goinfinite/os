package service

import (
	"errors"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	sslInfra "github.com/goinfinite/os/src/infra/ssl"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type SslService struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	sslQueryRepo          *sslInfra.SslQueryRepo
	sslCmdRepo            *sslInfra.SslCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewSslService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SslService {
	return &SslService{
		persistentDbSvc:       persistentDbSvc,
		sslQueryRepo:          sslInfra.NewSslQueryRepo(),
		sslCmdRepo:            sslInfra.NewSslCmdRepo(persistentDbSvc, transientDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *SslService) SslPairReadRequestFactory(
	serviceInput map[string]interface{},
	withMappings bool,
) (readRequestDto dto.ReadSslPairsRequest, err error) {
	if serviceInput["sslPairId"] == nil && serviceInput["id"] != nil {
		serviceInput["sslPairId"] = serviceInput["id"]
	}

	var sslPairIdPtr *valueObject.SslPairId
	if serviceInput["sslPairId"] != nil {
		sslPairId, err := valueObject.NewSslPairId(serviceInput["sslPairId"])
		if err != nil {
			return readRequestDto, err
		}
		sslPairIdPtr = &sslPairId
	}

	if serviceInput["virtualHostHostname"] == nil && serviceInput["hostname"] != nil {
		serviceInput["virtualHostHostname"] = serviceInput["hostname"]
	}

	var vhostHostnamePtr *valueObject.Fqdn
	if serviceInput["virtualHostHostname"] != nil {
		vhostHostname, err := valueObject.NewFqdn(serviceInput["virtualHostHostname"])
		if err != nil {
			return readRequestDto, err
		}
		vhostHostnamePtr = &vhostHostname
	}

	altNames := []valueObject.SslHostname{}
	if serviceInput["altNames"] != nil {
		var assertOk bool
		altNames, assertOk = serviceInput["altNames"].([]valueObject.SslHostname)
		if !assertOk {
			return readRequestDto, errors.New("InvalidAltNamesStructure")
		}
	}

	timeParamNames := []string{
		"issuedBeforeAt", "issuedAfterAt", "expiresBeforeAt", "expiresAfterAt",
	}
	timeParamPtrs := serviceHelper.TimeParamsParser(timeParamNames, serviceInput)

	requestPagination, err := serviceHelper.PaginationParser(
		serviceInput, useCase.SslPairsDefaultPagination,
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

func (service *SslService) Read(
	serviceInput map[string]interface{},
) ServiceOutput {
	readRequestDto, err := service.SslPairReadRequestFactory(serviceInput, false)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	readResponseDto, err := useCase.ReadSslPairs(service.sslQueryRepo, readRequestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, readResponseDto)
}

func (service *SslService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"virtualHosts", "certificate", "key"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	vhosts := []valueObject.Fqdn{}
	for _, rawVhost := range input["virtualHosts"].([]string) {
		vhost, err := valueObject.NewFqdn(rawVhost)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		vhosts = append(vhosts, vhost)
	}

	certContent, err := valueObject.NewSslCertificateContent(input["certificate"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}
	certEntity, err := entity.NewSslCertificate(certContent)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var chainCertsPtr *entity.SslCertificate
	if input["chainCertificates"] != nil {
		chainCertContent, err := valueObject.NewSslCertificateContent(input["chainCertificates"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("SslCertificateChainContentError"))
		}
		chainCertEntity, err := entity.NewSslCertificate(chainCertContent)
		if err != nil {
			return NewServiceOutput(UserError, errors.New("SslCertificateChainParseError"))
		}
		chainCertsPtr = &chainCertEntity
	}

	privateKeyContent, err := valueObject.NewSslPrivateKey(input["key"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
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

	createDto := dto.NewCreateSslPair(
		vhosts, certEntity, chainCertsPtr, privateKeyContent,
		operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)

	err = useCase.CreateSslPair(
		vhostQueryRepo, service.sslCmdRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "SslPairCreated")
}

func (service *SslService) CreatePubliclyTrusted(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
	if input["hostname"] != nil && input["virtualHostHostname"] == nil {
		input["virtualHostHostname"] = input["hostname"]
	}

	if input["vhostHostname"] != nil && input["virtualHostHostname"] == nil {
		input["virtualHostHostname"] = input["vhostHostname"]
	}

	requiredParams := []string{"virtualHostHostname"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	vhostHostname, err := valueObject.NewFqdn(input["virtualHostHostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
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

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteOsBinary + " ssl create-trusted"
		installParams := []string{
			"--hostname", vhostHostname.String(),
		}
		cliCmd += " " + strings.Join(installParams, " ")

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("CreatePubliclyTrustedSslPair")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("ssl")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSecs := uint16(1800)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "PubliclyTrustedSslPairCreationScheduled")
	}

	createDto := dto.NewCreatePubliclyTrustedSslPair(
		vhostHostname, operatorAccountId, operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)

	_, err = useCase.CreatePubliclyTrustedSslPair(
		vhostQueryRepo, service.sslCmdRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "PubliclyTrustedSslPairCreated")
}

func (service *SslService) Delete(input map[string]interface{}) ServiceOutput {
	if input["id"] == nil && input["sslPairId"] != nil {
		input["id"] = input["sslPairId"]
	}

	requiredParams := []string{"id"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	pairId, err := valueObject.NewSslPairId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
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

	err = useCase.DeleteSslPair(
		service.sslQueryRepo, service.sslCmdRepo, service.activityRecordCmdRepo,
		dto.NewDeleteSslPair(pairId, operatorAccountId, operatorIpAddress),
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "SslPairDeleted")
}
