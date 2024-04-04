package marketplaceInfra

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	"gopkg.in/yaml.v3"
)

//go:embed assets/*
var assets embed.FS

type MarketplaceQueryRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceQueryRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceQueryRepo {
	return &MarketplaceQueryRepo{
		persistentDbSvc: persistentDbSvc,
	}
}

func (repo *MarketplaceQueryRepo) getCatalogItemFromFilePath(
	catalogItemFilePath valueObject.UnixFilePath,
) (entity.MarketplaceCatalogItem, error) {
	var catalogItem entity.MarketplaceCatalogItem

	catalogItemFile, err := assets.Open(catalogItemFilePath.String())
	if err != nil {
		return catalogItem, err
	}
	defer catalogItemFile.Close()

	catalogItemFileExt, _ := catalogItemFilePath.GetFileExtension()
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

func (repo *MarketplaceQueryRepo) GetCatalogItems() (
	[]entity.MarketplaceCatalogItem, error,
) {
	catalogItems := []entity.MarketplaceCatalogItem{}

	catalogItemFiles, err := fs.ReadDir(assets, "assets")
	if err != nil {
		return catalogItems, errors.New(
			"GetMarketplaceCatalogItemsFilesError: " + err.Error(),
		)
	}

	if len(catalogItemFiles) == 0 {
		return catalogItems, errors.New("MarketplaceItemsEmpty")
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

		catalogItem, err := repo.getCatalogItemFromFilePath(
			catalogItemFilePath,
		)
		if err != nil {
			log.Printf(
				"GetMarketplaceCatalogItemError (%s): %s",
				catalogItemFileName,
				err.Error(),
			)
			continue
		}

		catalogItemIdInt := catalogItemFileIndex + 1
		catalogItemId, _ := valueObject.NewMarketplaceItemId(catalogItemIdInt)
		catalogItem.Id = catalogItemId

		catalogItems = append(catalogItems, catalogItem)
	}

	return catalogItems, nil
}

func (repo *MarketplaceQueryRepo) GetCatalogItemById(
	id valueObject.MarketplaceItemId,
) (entity.MarketplaceCatalogItem, error) {
	var catalogItem entity.MarketplaceCatalogItem

	catalogItems, err := repo.GetCatalogItems()
	if err != nil {
		return catalogItem, err
	}

	for _, catalogItem := range catalogItems {
		if catalogItem.Id.Get() != id.Get() {
			continue
		}

		return catalogItem, nil
	}

	return catalogItem, nil
}

func (repo *MarketplaceQueryRepo) GetInstalledItems() (
	[]entity.MarketplaceInstalledItem, error,
) {
	installedItemEntities := []entity.MarketplaceInstalledItem{}

	installedItemModels := []dbModel.MarketplaceInstalledItem{}
	err := repo.persistentDbSvc.Handler.Model(&dbModel.MarketplaceInstalledItem{}).
		Find(&installedItemModels).Error
	if err != nil {
		return installedItemEntities, errors.New(
			"DatabaseQueryMarketplaceInstalledItemsError",
		)
	}

	for _, installedItemModel := range installedItemModels {
		installedItemEntity, err := installedItemModel.ToEntity()
		if err != nil {
			log.Printf("MarketplaceInstalledItemModelToEntityError: %s", err.Error())
			continue
		}

		installedItemEntities = append(
			installedItemEntities,
			installedItemEntity,
		)
	}

	return installedItemEntities, nil
}
