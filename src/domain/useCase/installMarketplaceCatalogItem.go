package useCase

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
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

		dataFieldNameStr := dataField.Name.String()
		requiredDataFieldNames = append(requiredDataFieldNames, dataFieldNameStr)
	}

	if len(requiredDataFieldNames) == 0 {
		return nil
	}

	receivedDataFieldNames := map[string]interface{}{}
	for _, dataField := range receivedDataFields {
		dataFieldNameStr := dataField.Name.String()
		receivedDataFieldNames[dataFieldNameStr] = nil
	}

	missingDataFieldNames := []string{}
	for _, requiredDataFieldName := range requiredDataFieldNames {
		if _, isPresent := receivedDataFieldNames[requiredDataFieldName]; isPresent {
			continue
		}
		missingDataFieldNames = append(missingDataFieldNames, requiredDataFieldName)
	}

	if len(missingDataFieldNames) == 0 {
		return nil
	}

	return errors.New(
		"MissingRequiredDataFields: " + strings.Join(missingDataFieldNames, ","),
	)
}

func MarketplaceCatalogItemLookup(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	itemId *valueObject.MarketplaceItemId,
	itemSlug *valueObject.MarketplaceItemSlug,
) (itemEntity entity.MarketplaceCatalogItem, err error) {
	if itemId == nil && itemSlug == nil {
		return itemEntity, errors.New("ItemIdOrSlugRequired")
	}

	readDto := dto.ReadMarketplaceCatalogItemsRequest{}
	if itemId != nil {
		readDto.ItemId = itemId
		return marketplaceQueryRepo.ReadUniqueCatalogItem(readDto)
	}

	readDto.ItemSlug = itemSlug
	return marketplaceQueryRepo.ReadUniqueCatalogItem(readDto)
}

func InstallMarketplaceCatalogItem(
	marketplaceQueryRepo repository.MarketplaceQueryRepo,
	marketplaceCmdRepo repository.MarketplaceCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	vhostCmdRepo repository.VirtualHostCmdRepo,
	installDto dto.InstallMarketplaceCatalogItem,
) error {
	_, err := vhostQueryRepo.ReadByHostname(installDto.Hostname)
	if err != nil {
		return errors.New("VhostNotFound")
	}

	catalogItem, err := MarketplaceCatalogItemLookup(
		marketplaceQueryRepo, installDto.Id, installDto.Slug,
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
		slog.Error("InstallMarketplaceCatalogItem", slog.Any("error", err))
		return errors.New("InstallMarketplaceCatalogItemInfraError")
	}

	return nil
}
