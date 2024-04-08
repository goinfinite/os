package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CreateMarketplaceInstalledItem struct {
	Name          valueObject.MarketplaceItemName `json:"name"`
	Type          valueObject.MarketplaceItemType `json:"type"`
	RootDirectory valueObject.UnixFilePath        `json:"rootDirectory"`
	ServiceNames  []valueObject.ServiceName       `json:"serviceNames"`
	Mappings      []entity.Mapping                `json:"mappings"`
	AvatarUrl     valueObject.Url                 `json:"avatarUrl"`
}

func CreateNewMarketplaceInstalledItem(
	itemName valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	rootDirectory valueObject.UnixFilePath,
	serviceNames []valueObject.ServiceName,
	mappings []entity.Mapping,
	avatarUrl valueObject.Url,
) CreateMarketplaceInstalledItem {
	return CreateMarketplaceInstalledItem{
		Name:          itemName,
		Type:          itemType,
		RootDirectory: rootDirectory,
		ServiceNames:  serviceNames,
		Mappings:      mappings,
		AvatarUrl:     avatarUrl,
	}
}
