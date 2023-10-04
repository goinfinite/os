package useCase

import (
	"regexp"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/repository"
)

func GetSsls(
	sslQueryRepo repository.SslQueryRepo,
) ([]entity.Ssl, error) {
	sslList, err := sslQueryRepo.Get()
	if err != nil {
		matchErr, _ := regexp.MatchString("^(HttpdVhostsConfigEmpty|VhostConfigEmpty)$", err.Error())
		if !matchErr {
			return sslList, err
		}

		return sslList, nil
	}

	return sslList, err
}
