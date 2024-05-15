package useCase

import (
	"errors"
	"log"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
)

func requiredDataFieldsInspector(
	requiredDataFields []valueObject.MarketplaceCatalogItemDataField,
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
) error {
	requiredDataFieldNames := []string{}
	for _, dataField := range requiredDataFields {
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

	catalogItem, err := marketplaceQueryRepo.ReadCatalogItemById(
		installDto.Id,
	)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	err = requiredDataFieldsInspector(
		catalogItem.DataFields,
		installDto.DataFields,
	)
	if err != nil {
		return err
	}

	err = marketplaceCmdRepo.InstallItem(installDto)
	if err != nil {
		log.Printf("InstallMarketplaceCatalogItemError: %s", err.Error())
		return errors.New("InstallMarketplaceCatalogItemInfraError")
	}

	return nil
}
