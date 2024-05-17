package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func hasRequiredDataFields(
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
) bool {
	receivedDataFieldsKeysStr := []string{}
	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldsKeysStr = append(
			receivedDataFieldsKeysStr,
			receivedDataField.Name.String(),
		)
	}

	hasRequiredDataFields := true
	for _, catalogDataField := range catalogDataFields {
		if !catalogDataField.IsRequired {
			continue
		}

		requiredDataFieldStr := catalogDataField.Name.String()
		if !slices.Contains(receivedDataFieldsKeysStr, requiredDataFieldStr) {
			hasRequiredDataFields = false
			break
		}

		continue
	}

	return hasRequiredDataFields
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

	catalogItem, err := marketplaceQueryRepo.ReadCatalogItemById(
		installDto.Id,
	)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	hasRequiredDataFields := hasRequiredDataFields(
		installDto.DataFields,
		catalogItem.DataFields,
	)
	if !hasRequiredDataFields {
		return errors.New("MissingRequiredDataFieldKeys")
	}

	err = marketplaceCmdRepo.InstallItem(installDto)
	if err != nil {
		log.Printf("InstallMarketplaceCatalogItemError: %s", err.Error())
		return errors.New("InstallMarketplaceCatalogItemInfraError")
	}

	return nil
}
