package useCase

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func ReadAccessTokenDetails(
	authQueryRepo repository.AuthQueryRepo,
	accessToken tkValueObject.AccessTokenValue,
	trustedCidrs []tkValueObject.CidrBlock,
	ipAddress tkValueObject.IpAddress,
) (dto.AccessTokenDetails, error) {
	accessTokenDetails, err := authQueryRepo.ReadAccessTokenDetails(accessToken)
	if err != nil {
		return dto.AccessTokenDetails{}, err
	}

	if accessTokenDetails.IpAddress == nil {
		return accessTokenDetails, nil
	}

	for _, cidrBlock := range trustedCidrs {
		if cidrBlock.Contains(ipAddress) {
			return accessTokenDetails, nil
		}
	}

	if accessTokenDetails.IpAddress.String() != ipAddress.String() {
		return dto.AccessTokenDetails{}, errors.New("IpAddressChanged")
	}

	return accessTokenDetails, nil
}
