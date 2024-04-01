package dbModel

import (
	"log"
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type MarketplaceInstalledItem struct {
	Id            int64     `gorm:"column:Id;primaryKey;unique;autoIncrement"`
	Name          string    `gorm:"column:Name;not null"`
	Type          string    `gorm:"column:Type;not null"`
	RootDirectory string    `gorm:"column:RootDirectory;not null"`
	Services      string    `gorm:"column:Services;not null"`
	MappingsIds   string    `gorm:"column:MappingsIds;not null"`
	AvatarUrl     string    `gorm:"column:AvatarUrl;not null"`
	CreatedAt     time.Time `gorm:"column:CreatedAt;autoCreateTime"`
	UpdatedAt     time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

func (MarketplaceInstalledItem) TableName() string {
	return "marketplaceInstalledItem"
}

func (MarketplaceInstalledItem) ToModel(
	entity entity.MarketplaceInstalledItem,
) (MarketplaceInstalledItem, error) {
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
		Id:            entity.Id.Get(),
		Name:          entity.Name.String(),
		Type:          entity.Type.String(),
		RootDirectory: entity.RootDirectory.String(),
		Services:      svcNamesStr,
		MappingsIds:   mappingIdsStr,
		AvatarUrl:     entity.AvatarUrl.String(),
	}, nil
}

func (model MarketplaceInstalledItem) ToEntity() (
	entity.MarketplaceInstalledItem, error,
) {
	var mktplaceInstalledItem entity.MarketplaceInstalledItem

	id, err := valueObject.NewMktplaceItemId(model.Id)
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
	rawSvcsNameList := strings.Split(model.Services, ",")
	for _, rawSvcName := range rawSvcsNameList {
		svcName, err := valueObject.NewServiceName(rawSvcName)
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawSvcName)
		}

		svcsNameList = append(svcsNameList, svcName)
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
