package mktplaceInfra

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

//go:embed assets/*
var assets embed.FS

type MktplaceCatalogQueryRepo struct{}

func (repo MktplaceCatalogQueryRepo) getMarketplaceCatalogItemFromName(
	mktplaceItemName valueObject.MktplaceItemName,
) (entity.MarketplaceCatalogItem, error) {
	var mktplaceCatalogItem entity.MarketplaceCatalogItem

	mktplaceItemFilePath := "assets/" + mktplaceItemName.String() + ".json"
	mktplaceCatalogItemFile, err := assets.Open(mktplaceItemFilePath)
	if err != nil {
		return mktplaceCatalogItem, errors.New(
			"FailedToOpenMktCatalogItemFile: " + err.Error(),
		)
	}
	defer mktplaceCatalogItemFile.Close()

	mktplaceCatalogItemJsonDecoder := json.NewDecoder(mktplaceCatalogItemFile)
	err = mktplaceCatalogItemJsonDecoder.Decode(&mktplaceCatalogItem)
	if err != nil {
		return mktplaceCatalogItem, errors.New(
			"FailedToDecodeMktCatalogItemFile: " + err.Error(),
		)
	}

	return mktplaceCatalogItem, nil
}

func (repo MktplaceCatalogQueryRepo) GetItems() ([]entity.MarketplaceCatalogItem, error) {
	mktplaceCatalogItems := []entity.MarketplaceCatalogItem{}

	assetsEntries, err := fs.ReadDir(assets, "assets")
	if err != nil {
		return mktplaceCatalogItems, errors.New("FailedToGetMktAssetsFiles: " + err.Error())
	}

	if len(assetsEntries) == 0 {
		return mktplaceCatalogItems, errors.New("MktAssetsEmpty")
	}

	return mktplaceCatalogItems, nil
}
