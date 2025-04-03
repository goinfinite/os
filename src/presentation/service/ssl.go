package service

import (
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
	sslQueryRepo          sslInfra.SslQueryRepo
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
		sslQueryRepo:          sslInfra.SslQueryRepo{},
		sslCmdRepo:            sslInfra.NewSslCmdRepo(persistentDbSvc, transientDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *SslService) Read() ServiceOutput {
	pairsList, err := useCase.ReadSslPairs(service.sslQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, pairsList)
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
	cert, err := entity.NewSslCertificate(certContent)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
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
		vhosts, cert, privateKeyContent, operatorAccountId, operatorIpAddress,
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
	if input["sslPairId"] != nil && input["id"] == nil {
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

	deleteDto := dto.NewDeleteSslPair(pairId, operatorAccountId, operatorIpAddress)

	err = useCase.DeleteSslPair(
		service.sslQueryRepo, service.sslCmdRepo, service.activityRecordCmdRepo,
		deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "SslPairDeleted")
}
