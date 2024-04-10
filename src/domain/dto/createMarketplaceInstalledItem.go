package dto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type CreateMarketplaceInstalledItem struct {
	Name             valueObject.MarketplaceItemName `json:"name"`
	Type             valueObject.MarketplaceItemType `json:"type"`
	InstallDirectory valueObject.UnixFilePath        `json:"installDirectory"`
	ServiceNames     []valueObject.ServiceName       `json:"serviceNames"`
	Mappings         []entity.Mapping                `json:"mappings"`
	AvatarUrl        valueObject.Url                 `json:"avatarUrl"`
}

func CreateNewMarketplaceInstalledItem(
	itemName valueObject.MarketplaceItemName,
	itemType valueObject.MarketplaceItemType,
	installDirectory valueObject.UnixFilePath,
	serviceNames []valueObject.ServiceName,
	mappings []entity.Mapping,
	avatarUrl valueObject.Url,
) CreateMarketplaceInstalledItem {
	return CreateMarketplaceInstalledItem{
		Name:             itemName,
		Type:             itemType,
		InstallDirectory: installDirectory,
		ServiceNames:     serviceNames,
		Mappings:         mappings,
		AvatarUrl:        avatarUrl,
	}
}
