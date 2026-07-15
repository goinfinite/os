package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

var (
	ErrIpAddressMismatch = errors.New("IpAddressMismatch")
)

func ReadAccessTokenDetails(
	authQueryRepo repository.AuthQueryRepo,
	accessToken tkValueObject.AccessTokenValue,
	trustedCidrs []tkValueObject.CidrBlock,
	ipAddress tkValueObject.IpAddress,
) (accessTokenDetails dto.AccessTokenDetails, err error) {
	accessTokenDetails, err = authQueryRepo.ReadAccessTokenDetails(accessToken)
	if err != nil {
		return accessTokenDetails, err
	}

	if accessTokenDetails.IpAddress == nil {
		return accessTokenDetails, nil
	}

	for _, cidrBlock := range trustedCidrs {
		if cidrBlock.Contains(ipAddress) {
			return accessTokenDetails, nil
		}
	}

	tokenIpAddressStr := accessTokenDetails.IpAddress.String()
	operatorIpAddressStr := ipAddress.String()
	if tokenIpAddressStr != operatorIpAddressStr {
		slog.Debug(
			ErrIpAddressMismatch.Error(),
			slog.String("tokenIpAddress", tokenIpAddressStr),
			slog.String("operatorIpAddress", operatorIpAddressStr),
		)
		return accessTokenDetails, ErrIpAddressMismatch
	}

	return accessTokenDetails, nil
}
