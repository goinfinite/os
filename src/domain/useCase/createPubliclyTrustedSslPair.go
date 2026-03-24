package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

func CreatePubliclyTrustedSslPair(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
	createDto dto.CreatePubliclyTrustedSslPair,
) (sslPairId valueObject.SslPairId, err error) {
	sslPairId, err = sslCmdRepo.CreatePubliclyTrusted(createDto)
	if err != nil {
		slog.Error("CreatePubliclyTrustedSslPairError", slog.String("error", err.Error()))
		return sslPairId, errors.New("CreatePubliclyTrustedSslPairInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreatePubliclyTrustedSslPair(createDto, sslPairId)

	return sslPairId, nil
}
