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
	_, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &createDto.Hostname,
	})
	if err == nil {
		return errors.New("HostnameAlreadyInUse")
	}

	if createDto.Type == valueObject.VirtualHostTypeAlias && createDto.ParentHostname == nil {
		return errors.New("MissingAliasParentHostname")
	}

	isWildcardHostname := strings.HasPrefix(createDto.Hostname.String(), "*.")
	if isWildcardHostname {
		hostnameWithoutWildcardStr := strings.Replace(createDto.Hostname.String(), "*.", "", 1)
		hostnameWithoutWildcard, err := valueObject.NewFqdn(hostnameWithoutWildcardStr)
		if err != nil {
			return errors.New("RemoveWildcardFromHostnameError")
		}

		createDto.Type = valueObject.VirtualHostTypeWildcard
		createDto.Hostname = hostnameWithoutWildcard
	}

	if createDto.Type == valueObject.VirtualHostTypeAlias {
		parentVirtualHostEntity, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
			Hostname: createDto.ParentHostname,
		})
		if err != nil {
			slog.Error("ReadAliasParentVirtualHostError", slog.String("err", err.Error()))
			return errors.New("ReadAliasParentVirtualHostError")
		}

		if parentVirtualHostEntity.Type == valueObject.VirtualHostTypeAlias {
			return errors.New("AliasParentVirtualHostCannotAlsoBeAlias")
		}
	}

	if createDto.Type == valueObject.VirtualHostTypeWildcard {
		isWildcard := true
		createDto.IsWildcard = &isWildcard
	}

	err = vhostCmdRepo.Create(createDto)
	if err != nil {
		slog.Error("CreateVirtualHostError", slog.String("err", err.Error()))
		return errors.New("CreateVirtualHostInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		CreateVirtualHost(createDto)

	return nil
}
