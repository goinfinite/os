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

func mktplaceItemMappingFactory(
	vhost valueObject.Fqdn,
	rootDirectory valueObject.UnixFilePath,
) (dto.CreateMapping, error) {
	var mktplaceItemMapping dto.CreateMapping

	rootDirectoryWithoutTrailingSlash := strings.TrimSuffix(rootDirectory.String(), "/")
	mktplaceItemMappingPath, err := valueObject.NewMappingPath(
		rootDirectoryWithoutTrailingSlash,
	)
	if err != nil {
		return mktplaceItemMapping, err
	}

	mktplaceItemMappingMatchPattern, err := valueObject.NewMappingMatchPattern("begins-with")
	if err != nil {
		return mktplaceItemMapping, err
	}

	mktplaceItemMappingTargetType, err := valueObject.NewMappingTargetType("static-files")
	if err != nil {
		return mktplaceItemMapping, err
	}

	mktplaceItemMapping = dto.NewCreateMapping(
		vhost,
		mktplaceItemMappingPath,
		mktplaceItemMappingMatchPattern,
		mktplaceItemMappingTargetType,
		nil,
		nil,
		nil,
		nil,
	)

	return mktplaceItemMapping, nil
}

func InstallMarketplaceCatalogItem(
	MktplaceQueryRepo repository.MktplaceQueryRepo,
	MktplaceCmdRepo repository.MktplaceCmdRepo,
	vhostQueryRepo vhostInfra.VirtualHostQueryRepo,
	vhostCmdRepo vhostInfra.VirtualHostCmdRepo,
	installMktplaceCatalogItem dto.InstallMarketplaceCatalogItem,
) error {
	vhost, err := vhostQueryRepo.GetByHostname(installMktplaceCatalogItem.Hostname)
	if err != nil {
		return errors.New("VhostNotFound")
	}

	mktplaceCatalogItem, err := MktplaceQueryRepo.GetItemById(
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

	err = MktplaceCmdRepo.InstallItem(installMktplaceCatalogItem)
	if err != nil {
		log.Printf("InstallMktplaceCatalogItemError: %s", err.Error())
		return errors.New("InstallMktplaceCatalogItemInfraError")
	}

	mktplaceItemMapping, err := mktplaceItemMappingFactory(
		installMktplaceCatalogItem.Hostname,
		installMktplaceCatalogItem.RootDirectory,
	)
	if err != nil {
		log.Printf("CreateMktplaceItemMappingError: %s", err.Error())
		return errors.New("CreateMktplaceItemMappingError")
	}

	err = vhostCmdRepo.CreateMapping(mktplaceItemMapping)
	if err != nil {
		log.Printf("CreateMktplaceItemMappingError: %s", err.Error())
		return errors.New("CreateMktplaceItemMappingInfraError")
	}

	return nil
}
