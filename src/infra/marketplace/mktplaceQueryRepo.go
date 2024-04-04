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

func (repo *MarketplaceQueryRepo) getMarketplaceCatalogItemFromFilePath(
	marketplaceCatalogItemFilePath valueObject.UnixFilePath,
) (entity.MarketplaceCatalogItem, error) {
	var catalogItem entity.MarketplaceCatalogItem

	catalogItemFile, err := assets.Open(marketplaceCatalogItemFilePath.String())
	if err != nil {
		return catalogItem, err
	}
	defer catalogItemFile.Close()

	catalogItemFileExt, _ := marketplaceCatalogItemFilePath.GetFileExtension()
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

func (repo *MarketplaceQueryRepo) GetItems() (
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

		catalogItem, err := repo.getMarketplaceCatalogItemFromFilePath(
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

func (repo *MarketplaceQueryRepo) GetItemById(
	id valueObject.MarketplaceItemId,
) (entity.MarketplaceCatalogItem, error) {
	var marketplaceCatalogItem entity.MarketplaceCatalogItem

	marketplaceCatalogItems, err := repo.GetItems()
	if err != nil {
		return marketplaceCatalogItem, err
	}

	for _, catalogItem := range marketplaceCatalogItems {
		if catalogItem.Id.Get() != id.Get() {
			continue
		}

		marketplaceCatalogItem = catalogItem
	}

	return marketplaceCatalogItem, nil
}

func (repo *MarketplaceQueryRepo) GetInstalledItems() (
	[]entity.MarketplaceInstalledItem, error,
) {
	marketplaceInstalledItemEntities := []entity.MarketplaceInstalledItem{}

	marketplaceInstalledItemModels := []dbModel.MarketplaceInstalledItem{}
	err := repo.persistentDbSvc.Handler.Model(&dbModel.MarketplaceInstalledItem{}).
		Find(&marketplaceInstalledItemModels).Error
	if err != nil {
		return marketplaceInstalledItemEntities, errors.New(
			"DatabaseQueryMarketplaceInstalledItemsError",
		)
	}

	for _, marketplaceInstalledItemModel := range marketplaceInstalledItemModels {
		marketplaceInstalledItem, err := marketplaceInstalledItemModel.ToEntity()
		if err != nil {
			log.Printf("MarketplaceInstalledItemModelToEntityError: %s", err.Error())
			continue
		}

		marketplaceInstalledItemEntities = append(
			marketplaceInstalledItemEntities,
			marketplaceInstalledItem,
		)
	}

	return marketplaceInstalledItemEntities, nil
}
