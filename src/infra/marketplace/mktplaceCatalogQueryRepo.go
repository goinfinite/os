package mktplaceInfra

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	"gopkg.in/yaml.v3"
)

//go:embed assets/*
var assets embed.FS

type MktplaceCatalogQueryRepo struct{}

func (repo MktplaceCatalogQueryRepo) getMktCatalogItemFromFilePath(
	mktplaceItemFilePath valueObject.UnixFilePath,
) (entity.MarketplaceCatalogItem, error) {
	var mktplaceCatalogItem entity.MarketplaceCatalogItem

	mktplaceCatalogItemFile, err := assets.Open(mktplaceItemFilePath.String())
	if err != nil {
		return mktplaceCatalogItem, err
	}
	defer mktplaceCatalogItemFile.Close()

	mktplaceItemFileExt, _ := mktplaceItemFilePath.GetFileExtension()
	if mktplaceItemFileExt == "json" {
		mktplaceCatalogItemJsonDecoder := json.NewDecoder(mktplaceCatalogItemFile)
		err = mktplaceCatalogItemJsonDecoder.Decode(&mktplaceCatalogItem)
		if err != nil {
			return mktplaceCatalogItem, err
		}

		return mktplaceCatalogItem, nil
	}

	mktplaceCatalogItemYamlDecoder := yaml.NewDecoder(mktplaceCatalogItemFile)
	err = mktplaceCatalogItemYamlDecoder.Decode(&mktplaceCatalogItem)
	if err != nil {
		return mktplaceCatalogItem, err
	}

	return mktplaceCatalogItem, nil
}

func (repo MktplaceCatalogQueryRepo) GetItems() (
	[]entity.MarketplaceCatalogItem, error,
) {
	mktplaceCatalogItems := []entity.MarketplaceCatalogItem{}

	mktplaceItemFiles, err := fs.ReadDir(assets, "assets")
	if err != nil {
		return mktplaceCatalogItems, errors.New(
			"GetMktItemsFilesError: " + err.Error(),
		)
	}

	if len(mktplaceItemFiles) == 0 {
		return mktplaceCatalogItems, errors.New("MktItemsEmpty")
	}

	for mktplaceItemFileIndex, mktplaceItemFile := range mktplaceItemFiles {
		mktplaceItemFileName := mktplaceItemFile.Name()

		mktplaceItemFilePathStr := "assets/" + mktplaceItemFileName
		mktplaceItemFilePath, err := valueObject.NewUnixFilePath(
			mktplaceItemFilePathStr,
		)
		if err != nil {
			log.Printf(
				"%s (%s): %s", err.Error(),
				mktplaceItemFileName,
				mktplaceItemFilePathStr,
			)
			continue
		}

		mktplaceCatalogItem, err := repo.getMktCatalogItemFromFilePath(
			mktplaceItemFilePath,
		)
		if err != nil {
			log.Printf(
				"GetMktCatalogItemError (%s): %s",
				mktplaceItemFileName,
				err.Error(),
			)
			continue
		}

		mktplaceItemIdInt := mktplaceItemFileIndex + 1
		mktplaceItemId, _ := valueObject.NewMktplaceItemId(mktplaceItemIdInt)
		mktplaceCatalogItem.Id = mktplaceItemId

		mktplaceCatalogItems = append(mktplaceCatalogItems, mktplaceCatalogItem)
	}

	return mktplaceCatalogItems, nil
}
