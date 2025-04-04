package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func UninstallMarketplaceInstalledItemServices(
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	installedItemEntity entity.MarketplaceInstalledItem,
) {
	for _, serviceWithVersion := range installedItemEntity.Services {
		isServiceUninstallable := false
		serviceNameStr := serviceWithVersion.Name.String()

		serviceEntity, err := servicesQueryRepo.ReadFirstInstalledItem(
			dto.ReadFirstInstalledServiceItemsRequest{
				ServiceName: &serviceWithVersion.Name,
			},
		)
		if err != nil {
			slog.Error(
				"ReadServiceEntityError",
				slog.String("serviceName", serviceNameStr),
				slog.String("err", err.Error()),
			)
			continue
		}

		if serviceEntity.Type == valueObject.ServiceTypeDatabase {
			slog.Debug(
				"SkippingDatabaseServiceUninstall",
				slog.String("serviceName", serviceNameStr),
				slog.String("marketplaceItemId", installedItemEntity.Id.String()),
			)
			continue
		}

		targetValue, err := valueObject.NewMappingTargetValue(
			serviceNameStr, valueObject.MappingTargetTypeService,
		)
		if err != nil {
			slog.Error(
				"ServiceNameNotValidTargetValue",
				slog.String("serviceName", serviceNameStr),
				slog.String("err", err.Error()),
			)
			continue
		}

		mappingReadResponse, err := mappingQueryRepo.Read(dto.ReadMappingsRequest{
			Pagination:  dto.PaginationUnpaginated,
			TargetValue: &targetValue,
		})
		if err != nil {
			slog.Error(
				"ReadServiceMappingError",
				slog.String("serviceName", serviceNameStr),
				slog.String("mappingPath", serviceNameStr),
				slog.String("err", err.Error()),
			)
			continue
		}
		isServiceUninstallable = len(mappingReadResponse.Mappings) == 0

		if !isServiceUninstallable {
			slog.Debug(
				"SkippingServiceWithMappingsUninstall",
				slog.Uint64("mappingCount", uint64(len(mappingReadResponse.Mappings))),
				slog.String("serviceName", serviceNameStr),
				slog.String("mappingPath", serviceNameStr),
			)
			continue
		}

		err = servicesCmdRepo.Delete(serviceWithVersion.Name)
		if err != nil {
			slog.Error(
				"DeleteServiceInfraError",
				slog.String("serviceName", serviceNameStr),
				slog.String("err", err.Error()),
			)
			continue
		}
	}
}

func DeleteMarketplaceInstalledItem(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	marketplaceCmdRepo repository.MarketplaceCmdRepo,
	mappingQueryRepo repository.MappingQueryRepo,
	mappingCmdRepo repository.MappingCmdRepo,
	servicesQueryRepo repository.ServicesQueryRepo,
	servicesCmdRepo repository.ServicesCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	deleteDto dto.DeleteMarketplaceInstalledItem,
) error {
	installedItemEntity, err := marketplaceQueryRepo.ReadFirstInstalledItem(
		dto.ReadMarketplaceInstalledItemsRequest{
			MarketplaceInstalledItemId: &deleteDto.InstalledId,
		},
	)
	if err != nil {
		slog.Error("ReadMarketplaceInstalledItemError", slog.String("err", err.Error()))
		return errors.New("MarketplaceInstalledItemNotFound")
	}

	for _, mappingEntity := range installedItemEntity.Mappings {
		err = mappingCmdRepo.Delete(mappingEntity.Id)
		if err != nil {
			slog.Error(
				"DeleteInstalledItemMappingError",
				slog.String("mappingId", mappingEntity.Id.String()),
				slog.String("mappingPath", mappingEntity.Path.String()),
				slog.String("err", err.Error()),
			)
			continue
		}
	}

	if deleteDto.ShouldUninstallServices {
		UninstallMarketplaceInstalledItemServices(
			servicesQueryRepo, servicesCmdRepo, mappingQueryRepo, installedItemEntity,
		)
	}

	err = marketplaceCmdRepo.UninstallItem(deleteDto)
	if err != nil {
		slog.Error("UninstallMarketplaceItemError", slog.String("err", err.Error()))
		return errors.New("UninstallMarketplaceItemInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		DeleteMarketplaceInstalledItem(deleteDto)

	return nil
}
