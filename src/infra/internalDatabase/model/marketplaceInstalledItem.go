package dbModel

import (
	"log"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	ID            uint `gorm:"primarykey"`
	Name          string
	Type          string
	RootDirectory string
	Services      string
	MappingsIds   string
	AvatarUrl     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (MarketplaceInstalledItem) TableName() string {
	return "mktplace_installed_items"
}

func (MarketplaceInstalledItem) ToModel(
	entity entity.MarketplaceInstalledItem,
	withoutId bool,
) (MarketplaceInstalledItem, error) {
	idUint := uint(entity.Id.Get())
	if withoutId {
		idUint = 0
	}

	var svcNamesListStr []string
	for _, svcName := range entity.Services {
		svcNamesListStr = append(svcNamesListStr, svcName.String())
	}
	svcNamesStr := strings.Join(svcNamesListStr, ",")

	var mappingIdsListStr []string
	for _, mapping := range entity.Mappings {
		mappingIdsListStr = append(mappingIdsListStr, mapping.Id.String())
	}
	mappingIdsStr := strings.Join(mappingIdsListStr, ",")

	return MarketplaceInstalledItem{
		ID:            idUint,
		Name:          entity.Name.String(),
		Type:          entity.Type.String(),
		RootDirectory: entity.RootDirectory.String(),
		Services:      svcNamesStr,
		MappingsIds:   mappingIdsStr,
		AvatarUrl:     entity.AvatarUrl.String(),
		CreatedAt:     entity.CreatedAt.GetUnixTime(),
		UpdatedAt:     entity.UpdatedAt.GetUnixTime(),
	}, nil
}

func (model MarketplaceInstalledItem) ToEntity() (
	entity.MarketplaceInstalledItem, error,
) {
	var mktplaceInstalledItem entity.MarketplaceInstalledItem

	id, err := valueObject.NewMktplaceItemId(model.ID)
	if err != nil {
		return mktplaceInstalledItem, err
	}

	itemName, err := valueObject.NewMktplaceItemName(model.Name)
	if err != nil {
		return mktplaceInstalledItem, err
	}

	itemType, err := valueObject.NewMktplaceItemType(model.Type)
	if err != nil {
		return mktplaceInstalledItem, err
	}

	rootDirectory, err := valueObject.NewUnixFilePath(model.RootDirectory)
	if err != nil {
		return mktplaceInstalledItem, err
	}

	svcsNameList := []valueObject.ServiceName{}
	if len(model.Services) > 0 {
		rawSvcsNameList := strings.Split(model.Services, ",")
		for _, rawSvcName := range rawSvcsNameList {
			svcName, err := valueObject.NewServiceName(rawSvcName)
			if err != nil {
				log.Printf("%s: %s", err.Error(), rawSvcName)
			}

			svcsNameList = append(svcsNameList, svcName)
		}
	}

	// Arrumar isso aqui depois
	mappings := []entity.Mapping{}

	avatarUrl, err := valueObject.NewUrl(model.AvatarUrl)
	if err != nil {
		return mktplaceInstalledItem, err
	}

	rawCreatedAtUnix := model.CreatedAt.UTC().Unix()
	createdAt := valueObject.UnixTime(rawCreatedAtUnix)

	rawUpdatedAtUnix := model.UpdatedAt.UTC().Unix()
	updatedAt := valueObject.UnixTime(rawUpdatedAtUnix)

	return entity.NewMarketplaceInstalledItem(
		id,
		itemName,
		itemType,
		rootDirectory,
		svcsNameList,
		mappings,
		avatarUrl,
		createdAt,
		updatedAt,
	), nil
}
