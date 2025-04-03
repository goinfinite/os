package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreatePubliclyTrustedSslPair(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	sslCmdRepo repository.SslCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
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
