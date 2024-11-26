package useCase

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func CreateVirtualHost(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	createDto dto.CreateVirtualHost,
) error {
	_, err := vhostQueryRepo.ReadByHostname(createDto.Hostname)
	if err == nil {
		return errors.New("VirtualHostAlreadyExists")
	}

	isAlias := createDto.Type.String() == "alias"
	if isAlias && createDto.ParentHostname == nil {
		return errors.New("AliasMustHaveParentHostname")
	}

	hostnameStr := createDto.Hostname.String()
	hasWildcardInHostname := strings.HasPrefix(hostnameStr, "*.")
	if hasWildcardInHostname {
		hostnameWithoutWildcardStr := strings.Replace(hostnameStr, "*.", "", 1)
		hostnameWithoutWildcard, err := valueObject.NewFqdn(hostnameWithoutWildcardStr)
		if err != nil {
			return errors.New("FailedToRemoveWildcardFromHostname: " + err.Error())
		}

		createDto.Hostname = hostnameWithoutWildcard
	}

	err = vhostCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateVirtualHostError", slog.Any("err", err))
		return errors.New("CreateVirtualHostInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateVirtualHost(createDto)

	return nil
}
