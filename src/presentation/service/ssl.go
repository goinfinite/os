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
}

func NewSslService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *SslService {
	return &SslService{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
	}
}

func (service *SslService) Read() ServiceOutput {
	sslQueryRepo := sslInfra.SslQueryRepo{}
	pairsList, err := useCase.ReadSslPairs(sslQueryRepo)
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

	sslCmdRepo := sslInfra.NewSslCmdRepo(service.persistentDbSvc, service.transientDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)

	err = useCase.CreateSslPair(sslCmdRepo, vhostQueryRepo, dto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "SslPairCreated")
}
