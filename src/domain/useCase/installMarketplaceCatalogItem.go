package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

func hasRequiredDataFields(
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	requiredDataFields []valueObject.MarketplaceCatalogItemDataField,
) bool {
	receivedDataFieldsKeysStr := []string{}
	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldsKeysStr = append(
			receivedDataFieldsKeysStr,
			receivedDataField.Key.String(),
		)
	}

	hasRequiredDataFields := true
	for _, requiredDataField := range requiredDataFields {
		if !requiredDataField.IsRequired {
			continue
		}

		requiredDataFieldStr := requiredDataField.Key.String()
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
	vhostQueryRepo vhostInfra.VirtualHostQueryRepo,
	vhostCmdRepo vhostInfra.VirtualHostCmdRepo,
	installDto dto.InstallMarketplaceCatalogItem,
) error {
	_, err := vhostQueryRepo.GetByHostname(installDto.Hostname)
	if err != nil {
		return errors.New("VhostNotFound")
	}

	catalogItem, err := marketplaceQueryRepo.GetCatalogItemById(
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
