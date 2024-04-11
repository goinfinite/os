package dbModel

import (
	"log"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	ID               uint `gorm:"primarykey"`
	Name             string
	Type             string
	InstallDirectory string
	ServiceNames     string
	AvatarUrl        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (MarketplaceInstalledItem) TableName() string {
	return "marketplace_installed_items"
}

func (model MarketplaceInstalledItem) ToEntity() (
	entity.MarketplaceInstalledItem, error,
) {
	var marketplaceInstalledItem entity.MarketplaceInstalledItem

	id, err := valueObject.NewMarketplaceInstalledItemId(model.ID)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	itemName, err := valueObject.NewMarketplaceItemName(model.Name)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	itemType, err := valueObject.NewMarketplaceItemType(model.Type)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	installDirectory, err := valueObject.NewUnixFilePath(model.InstallDirectory)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	svcsNameList := []valueObject.ServiceName{}
	if len(model.ServiceNames) > 0 {
		rawSvcsNameList := strings.Split(model.ServiceNames, ",")
		for _, rawSvcName := range rawSvcsNameList {
			svcName, err := valueObject.NewServiceName(rawSvcName)
			if err != nil {
				log.Printf("%s: %s", err.Error(), rawSvcName)
			}

			svcsNameList = append(svcsNameList, svcName)
		}
	}

	mappings := []entity.Mapping{}

	avatarUrl, err := valueObject.NewUrl(model.AvatarUrl)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	rawCreatedAtUnix := model.CreatedAt.UTC().Unix()
	createdAt := valueObject.UnixTime(rawCreatedAtUnix)

	rawUpdatedAtUnix := model.UpdatedAt.UTC().Unix()
	updatedAt := valueObject.UnixTime(rawUpdatedAtUnix)

	return entity.NewMarketplaceInstalledItem(
		id,
		itemName,
		itemType,
		installDirectory,
		svcsNameList,
		mappings,
		avatarUrl,
		createdAt,
		updatedAt,
	), nil
}
