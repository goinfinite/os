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

func inspectReceivedDataFields(
	receivedDataFields []valueObject.MarketplaceInstallableItemDataField,
	requiredDataFields []valueObject.MarketplaceCatalogItemDataField,
) error {
	receivedDataFieldsStrList := []string{}
	for _, receivedDataField := range receivedDataFields {
		receivedDataFieldsStrList = append(
			receivedDataFieldsStrList,
			receivedDataField.Name.String(),
		)
	}

	missingRequiredDataFields := []string{}
	for _, requiredDataField := range requiredDataFields {
		requiredDataFieldNameStr := requiredDataField.Name.String()
		if !slices.Contains(receivedDataFieldsStrList, requiredDataFieldNameStr) {
			missingRequiredDataFields = append(
				missingRequiredDataFields,
				requiredDataFieldNameStr,
			)
		}
	}

	if len(missingRequiredDataFields) > 0 {
		missingRequiredDataFieldsStr := strings.Join(missingRequiredDataFields, ",")
		return errors.New("MissingRequiredDataFields: " + missingRequiredDataFieldsStr)
	}

	return nil
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

	err = inspectReceivedDataFields(
		installDto.DataFields,
		catalogItem.DataFields,
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
