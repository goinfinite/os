package service

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
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
		service.sslCmdRepo, vhostQueryRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "SslPairCreated")
}

func (service *SslService) Delete(input map[string]interface{}) ServiceOutput {
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

func (service *SslService) DeleteVhosts(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"id", "virtualHosts"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	pairId, err := valueObject.NewSslPairId(input["id"])
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

	dto := dto.NewDeleteSslPairVhosts(pairId, vhosts)

	err = useCase.DeleteSslPairVhosts(service.sslQueryRepo, service.sslCmdRepo, dto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "SslPairVhostsDeleted")
}
