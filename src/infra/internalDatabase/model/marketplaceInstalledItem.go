package dbModel

import (
	"log"
	"strings"
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	ID               uint   `gorm:"primarykey"`
	Name             string `gorm:"not null"`
	Hostname         string `gorm:"not null"`
	Type             string `gorm:"not null"`
	UrlPath          string `gorm:"not null"`
	InstallDirectory string `gorm:"not null"`
	InstallUuid      string `gorm:"not null"`
	Services         string
	Mappings         []Mapping
	AvatarUrl        string    `gorm:"not null"`
	Slug             string    `gorm:"not null"`
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time `gorm:"not null"`
}

func (MarketplaceInstalledItem) TableName() string {
	return "marketplace_installed_items"
}

func (model MarketplaceInstalledItem) ToEntity() (
	entity.MarketplaceInstalledItem, error,
) {
	var marketplaceInstalledItem entity.MarketplaceInstalledItem

	id, err := valueObject.NewMarketplaceItemId(model.ID)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	itemName, err := valueObject.NewMarketplaceItemName(model.Name)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	hostname, err := valueObject.NewFqdn(model.Hostname)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	itemType, err := valueObject.NewMarketplaceItemType(model.Type)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	urlPath, err := valueObject.NewUrlPath(model.UrlPath)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	installDirectory, err := valueObject.NewUnixFilePath(model.InstallDirectory)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	installUuid, err := valueObject.NewMarketplaceInstalledItemUuid(
		model.InstallUuid,
	)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	serviceNamesWithVersion := []valueObject.ServiceNameWithVersion{}
	if len(model.Services) > 0 {
		rawServicesList := strings.Split(model.Services, ",")
		for _, rawService := range rawServicesList {
			serviceNameWithVersion, err := valueObject.NewServiceNameWithVersionFromString(
				rawService,
			)
			if err != nil {
				log.Printf("%s: %s", err.Error(), rawService)
			}
			serviceNamesWithVersion = append(serviceNamesWithVersion, serviceNameWithVersion)
		}
	}

	mappings := []entity.Mapping{}
	if len(model.Mappings) > 0 {
		for _, mappingModel := range model.Mappings {
			mapping, err := mappingModel.ToEntity()
			if err != nil {
				return marketplaceInstalledItem, err
			}
			mappings = append(mappings, mapping)
		}
	}

	avatarUrl, err := valueObject.NewUrl(model.AvatarUrl)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	slug, err := valueObject.NewMarketplaceItemSlug(model.Slug)
	if err != nil {
		return marketplaceInstalledItem, err
	}

	return entity.NewMarketplaceInstalledItem(
		id,
		itemName,
		hostname,
		itemType,
		urlPath,
		installDirectory,
		installUuid,
		serviceNamesWithVersion,
		mappings,
		avatarUrl,
		slug,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
