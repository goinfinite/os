package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
)

var VirtualHostsDefaultPagination dto.Pagination = dto.Pagination{
	PageNumber:   0,
	ItemsPerPage: 10,
}

func ReadVirtualHosts(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	requestDto dto.ReadVirtualHostsRequest,
) (responseDto dto.ReadVirtualHostsResponse, err error) {
	responseDto, err = vhostQueryRepo.Read(requestDto)
	if err != nil {
		slog.Error("ReadVirtualHostsError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadVirtualHostsInfraError")
	}

	return responseDto, nil
}
