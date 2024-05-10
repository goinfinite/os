package useCase

import (
	"errors"
	"log"
	"slices"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

func checkRequiredDataFields(
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	catalogDataFields []valueObject.MarketplaceCatalogItemDataField,
) []valueObject.MarketplaceCatalogItemDataField {
	missingRequiredDataFields := []valueObject.MarketplaceCatalogItemDataField{}

	receivedDataFieldsKeysStr := []string{}
	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldsKeysStr = append(
			receivedDataFieldsKeysStr,
			receivedDataField.Name.String(),
		)
	}

	for _, catalogDataField := range catalogDataFields {
		if !catalogDataField.IsRequired {
			continue
		}

		requiredDataFieldStr := catalogDataField.Name.String()
		if !slices.Contains(receivedDataFieldsKeysStr, requiredDataFieldStr) {
			missingRequiredDataFields = append(
				missingRequiredDataFields,
				catalogDataField,
			)
		}
	}

	return missingRequiredDataFields
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

	missingRequiredDataFields := checkRequiredDataFields(
		installDto.DataFields,
		catalogItem.DataFields,
	)

	hasRequiredDataFields := len(missingRequiredDataFields) == 0
	if !hasRequiredDataFields {
		missingDataFieldsStrList := []string{}
		for _, missingDataField := range missingRequiredDataFields {
			missingDataFieldsStrList = append(
				missingDataFieldsStrList,
				missingDataField.Name.String(),
			)
		}
		missingDataFieldsStr := strings.Join(missingDataFieldsStrList, ", ")

		return errors.New("MissingRequiredDataField: " + missingDataFieldsStr)
	}

	err = marketplaceCmdRepo.InstallItem(installDto)
	if err != nil {
		log.Printf("InstallMarketplaceCatalogItemError: %s", err.Error())
		return errors.New("InstallMarketplaceCatalogItemInfraError")
	}

	return nil
}
