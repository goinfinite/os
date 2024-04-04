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
	dtoDataFields []valueObject.MarketplaceItemDataField,
	requiredDataFields []valueObject.MarketplaceItemDataField,
) bool {
	dtoDataFieldsKeysStr := []string{}
	for _, dtoDataField := range dtoDataFields {
		dtoDataFieldsKeysStr = append(
			dtoDataFieldsKeysStr,
			dtoDataField.Key.String(),
		)
	}

	hasRequiredDataFields := true
	for _, requiredDataField := range requiredDataFields {
		requiredDataFieldStr := requiredDataField.Key.String()
		if !slices.Contains(dtoDataFieldsKeysStr, requiredDataFieldStr) {
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
	dto dto.InstallMarketplaceCatalogItem,
) error {
	vhost, err := vhostQueryRepo.GetByHostname(dto.Hostname)
	if err != nil {
		return errors.New("VhostNotFound")
	}

	marketplaceCatalogItem, err := marketplaceQueryRepo.GetCatalogItemById(
		dto.Id,
	)
	if err != nil {
		return errors.New("MarketplaceCatalogItemNotFound")
	}

	hasRequiredDataFields := hasRequiredDataFields(
		dto.DataFields,
		marketplaceCatalogItem.DataFields,
	)
	if !hasRequiredDataFields {
		return errors.New("MissingRequiredDataFieldKeys")
	}

	marketplaceItemRootDir := dto.RootDirectory
	rawCorrectRootDir := vhost.RootDirectory.String() + marketplaceItemRootDir.String()
	rootDirAbsolutePath, err := valueObject.NewUnixFilePath(rawCorrectRootDir)
	if err != nil {
		return err
	}
	dto.RootDirectory = rootDirAbsolutePath

	err = marketplaceCmdRepo.InstallItem(dto)
	if err != nil {
		log.Printf("InstallMarketplaceCatalogItemError: %s", err.Error())
		return errors.New("InstallMarketplaceCatalogItemInfraError")
	}

	return nil
}
