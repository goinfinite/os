package useCase

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

func requiredDataFieldsInspector(
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
) error {
	requiredDataFieldNames := []string{}
	for _, dataField := range catalogDataFields {
		if !dataField.IsRequired {
			continue
		}

		requiredDataFieldNames = append(requiredDataFieldNames, dataField.Name.String())
	}

	if len(requiredDataFieldNames) == 0 {
		return nil
	}

	receivedDataFieldNames := map[string]interface{}{}
	for _, dataField := range receivedDataFields {
		receivedDataFieldNames[dataField.Name.String()] = nil
	}

	missingDataFieldNames := []string{}
	for _, requiredDataFieldName := range requiredDataFieldNames {
		if _, isPresent := receivedDataFieldNames[requiredDataFieldName]; isPresent {
			continue
		}
		missingDataFieldNames = append(missingDataFieldNames, requiredDataFieldName)
	}

	if len(missingDataFieldNames) > 0 {
		return errors.New(
			"MissingRequiredDataFields: " + strings.Join(missingDataFieldNames, ","),
		)
	}

	return nil
}

func InstallMarketplaceCatalogItem(
	vhostQueryRepo repository.VirtualHostQueryRepo,
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	marketplaceCmdRepo repository.MarketplaceCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
	installDto dto.InstallMarketplaceCatalogItem,
) error {
	if installDto.Id == nil && installDto.Slug == nil {
		return errors.New("ItemIdOrSlugRequired")
	}

	_, err := vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &installDto.Hostname,
	})
	if err != nil {
		return errors.New("VirtualHostNotFound")
	}

	catalogItem, err := marketplaceQueryRepo.ReadFirstCatalogItem(
		dto.ReadMarketplaceCatalogItemsRequest{
			MarketplaceCatalogItemId:   installDto.Id,
			MarketplaceCatalogItemSlug: installDto.Slug,
		},
	)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}
	installDto.Id = &catalogItem.Id

	err = requiredDataFieldsInspector(catalogItem.DataFields, installDto.DataFields)
	if err != nil {
		return err
	}

	err = marketplaceCmdRepo.InstallItem(installDto)
	if err != nil {
		slog.Error("InstallMarketplaceCatalogItem", slog.String("err", err.Error()))
		return errors.New("InstallMarketplaceCatalogItemInfraError")
	}

	NewCreateSecurityActivityRecord(activityRecordCmdRepo).
		InstallMarketplaceCatalogItem(installDto)

	return nil
}
