package useCase

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetSslPairs(
	sslQueryRepo repository.SslQueryRepo,
) ([]dto.GetSslPair, error) {
	sslPairList, err := sslQueryRepo.GetSslPairs()
	if err != nil {
		return []dto.GetSslPair{}, err
	}

	sslPairListFormatted := []dto.GetSslPair{}
	for _, ssl := range sslPairList {
		sslFormatted := dto.NewGetSslPair(
			ssl.SerialNumber,
			ssl.Hostname,
			ssl.Certificate,
			ssl.Key,
			ssl.ChainCertificates,
		)
		sslPairListFormatted = append(sslPairListFormatted, sslFormatted)
	}

	return sslPairListFormatted, nil
}
