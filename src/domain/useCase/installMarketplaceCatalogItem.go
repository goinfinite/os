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
	dtoDataFields []valueObject.DataField,
	requiredDataFields []valueObject.DataField,
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
	mktplaceQueryRepo repository.MktplaceQueryRepo,
	mktplaceCmdRepo repository.MktplaceCmdRepo,
	vhostQueryRepo vhostInfra.VirtualHostQueryRepo,
	vhostCmdRepo vhostInfra.VirtualHostCmdRepo,
	installMktplaceCatalogItem dto.InstallMarketplaceCatalogItem,
) error {
	vhost, err := vhostQueryRepo.GetByHostname(installMktplaceCatalogItem.Hostname)
	if err != nil {
		return errors.New("VhostNotFound")
	}

	mktplaceCatalogItem, err := mktplaceQueryRepo.GetItemById(
		installMktplaceCatalogItem.Id,
	)
	if err != nil {
		return errors.New("MktplaceCatalogItemNotFound")
	}

	hasRequiredDataFields := hasRequiredDataFields(
		installMktplaceCatalogItem.DataFields,
		mktplaceCatalogItem.DataFields,
	)
	if !hasRequiredDataFields {
		return errors.New("MissingRequiredDataFieldKeys")
	}

	mktplaceItemRootDir := installMktplaceCatalogItem.RootDirectory
	rawCorrectRootDir := vhost.RootDirectory.String() + mktplaceItemRootDir.String()
	rootDirAbsolutePath, err := valueObject.NewUnixFilePath(rawCorrectRootDir)
	if err != nil {
		return err
	}
	installMktplaceCatalogItem.RootDirectory = rootDirAbsolutePath

	err = mktplaceCmdRepo.InstallItem(installMktplaceCatalogItem)
	if err != nil {
		log.Printf("InstallMktplaceCatalogItemError: %s", err.Error())
		return errors.New("InstallMktplaceCatalogItemInfraError")
	}

	return nil
}
