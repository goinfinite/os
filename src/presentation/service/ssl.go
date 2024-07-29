package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
)

type SslService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	sslQueryRepo    sslInfra.SslQueryRepo
	sslCmdRepo      *sslInfra.SslCmdRepo
}

func NewSslService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *SslService {
	return &SslService{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		sslQueryRepo:    sslInfra.SslQueryRepo{},
		sslCmdRepo:      sslInfra.NewSslCmdRepo(persistentDbSvc, transientDbSvc),
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

	dto := dto.NewCreateSslPair(vhosts, cert, privateKeyContent)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)

	err = useCase.CreateSslPair(service.sslCmdRepo, vhostQueryRepo, dto)
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

	pairId, err := valueObject.NewSslId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	err = useCase.DeleteSslPair(service.sslQueryRepo, service.sslCmdRepo, pairId)
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

	pairId, err := valueObject.NewSslId(input["id"])
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
