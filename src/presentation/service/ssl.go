package service

import (
	"github.com/speedianet/os/src/domain/useCase"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
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
