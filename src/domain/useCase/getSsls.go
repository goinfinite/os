package useCase

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetSsls(
	sslQueryRepo repository.SslQueryRepo,
) ([]dto.GetSsl, error) {
	sslList, err := sslQueryRepo.Get()
	if err != nil {
		return []dto.GetSsl{}, err
	}

	sslListFormatted := []dto.GetSsl{}
	for _, ssl := range sslList {
		sslFormatted := dto.NewGetSsl(
			ssl.Id,
			ssl.Hostname,
			ssl.Certificate,
			ssl.Key,
			ssl.ChainCertificates,
		)
		sslListFormatted = append(sslListFormatted, sslFormatted)
	}

	return sslListFormatted, nil
}
