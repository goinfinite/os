package mktplaceInfra

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	"gopkg.in/yaml.v3"
)

//go:embed assets/*
var assets embed.FS

type MktplaceCatalogQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMktplaceCatalogQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MktplaceCatalogQueryRepo {
	return &MktplaceCatalogQueryRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *MktplaceCatalogQueryRepo) getMktCatalogItemFromFilePath(
	mktCatalogItemFilePath valueObject.UnixFilePath,
) (entity.MarketplaceCatalogItem, error) {
	var catalogItem entity.MarketplaceCatalogItem

	catalogItemFile, err := assets.Open(mktCatalogItemFilePath.String())
	if err != nil {
		return catalogItem, err
	}
	defer catalogItemFile.Close()

	catalogItemFileExt, _ := mktCatalogItemFilePath.GetFileExtension()
	if catalogItemFileExt == "json" {
		catalogItemJsonDecoder := json.NewDecoder(catalogItemFile)
		err = catalogItemJsonDecoder.Decode(&catalogItem)
		if err != nil {
			return catalogItem, err
		}

		return catalogItem, nil
	}

	catalogItemYamlDecoder := yaml.NewDecoder(catalogItemFile)
	err = catalogItemYamlDecoder.Decode(&catalogItem)
	if err != nil {
		return catalogItem, err
	}

	return catalogItem, nil
}

func (repo *MktplaceCatalogQueryRepo) GetItems() (
	[]entity.MarketplaceCatalogItem, error,
) {
	catalogItems := []entity.MarketplaceCatalogItem{}

	catalogItemFiles, err := fs.ReadDir(assets, "assets")
	if err != nil {
		return catalogItems, errors.New(
			"GetMktCatalogItemsFilesError: " + err.Error(),
		)
	}

	if len(catalogItemFiles) == 0 {
		return catalogItems, errors.New("MktItemsEmpty")
	}

	for catalogItemFileIndex, catalogItemFile := range catalogItemFiles {
		catalogItemFileName := catalogItemFile.Name()

		catalogItemFilePathStr := "assets/" + catalogItemFileName
		catalogItemFilePath, err := valueObject.NewUnixFilePath(
			catalogItemFilePathStr,
		)
		if err != nil {
			log.Printf(
				"%s (%s): %s", err.Error(),
				catalogItemFileName,
				catalogItemFilePathStr,
			)
			continue
		}

		catalogItem, err := repo.getMktCatalogItemFromFilePath(
			catalogItemFilePath,
		)
		if err != nil {
			log.Printf(
				"GetMktCatalogItemError (%s): %s",
				catalogItemFileName,
				err.Error(),
			)
			continue
		}

		catalogItemIdInt := catalogItemFileIndex + 1
		catalogItemId, _ := valueObject.NewMktplaceItemId(catalogItemIdInt)
		catalogItem.Id = catalogItemId

		catalogItems = append(catalogItems, catalogItem)
	}

	return catalogItems, nil
}
