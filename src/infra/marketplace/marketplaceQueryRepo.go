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

func (repo *MarketplaceQueryRepo) getCatalogItemMapFromFilePath(
	catalogItemFilePath valueObject.UnixFilePath,
) (map[string]interface{}, error) {
	var catalogItemMap map[string]interface{}

	catalogItemFile, err := assets.Open(catalogItemFilePath.String())
	if err != nil {
		return catalogItemMap, err
	}
	defer catalogItemFile.Close()

	catalogItemFileExt, _ := catalogItemFilePath.GetFileExtension()
	isYamlFile := catalogItemFileExt == "yml" || catalogItemFileExt == "yaml"
	if isYamlFile {
		catalogItemYamlDecoder := yaml.NewDecoder(catalogItemFile)
		err = catalogItemYamlDecoder.Decode(&catalogItemMap)
		if err != nil {
			return catalogItemMap, err
		}

		return catalogItemMap, nil
	}

	catalogItemJsonDecoder := json.NewDecoder(catalogItemFile)
	err = catalogItemJsonDecoder.Decode(&catalogItemMap)
	if err != nil {
		return catalogItemMap, err
	}

	return catalogItemMap, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemServiceNames(
	catalogItemSvcNamesMap interface{},
) ([]valueObject.ServiceName, error) {
	itemSvcNames := []valueObject.ServiceName{}

	rawItemSvcNames, assertOk := catalogItemSvcNamesMap.([]interface{})
	if !assertOk {
		return itemSvcNames, errors.New("InvalidMarketplaceCatalogItemServiceNames")
	}

	for _, rawItemSvcName := range rawItemSvcNames {
		itemSvcName, err := valueObject.NewServiceName(rawItemSvcName.(string))
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawItemSvcName)
			continue
		}

		itemSvcNames = append(itemSvcNames, itemSvcName)
	}

	return itemSvcNames, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemMappings(
	catalogItemMappingsMap interface{},
) ([]valueObject.MarketplaceItemMapping, error) {
	itemMappings := []valueObject.MarketplaceItemMapping{}

	rawItemMappings, assertOk := catalogItemMappingsMap.([]interface{})
	if !assertOk {
		return itemMappings, errors.New("InvalidMarketplaceCatalogItemMappings")
	}

	for _, rawItemMapping := range rawItemMappings {
		rawItemMappingMap, assertOk := rawItemMapping.(map[string]interface{})
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemMapping: %+v", rawItemMapping)
			continue
		}

		rawPath, assertOk := rawItemMappingMap["path"].(string)
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemMappingPath: %s", rawPath)
			continue
		}
		path, err := valueObject.NewMappingPath(rawPath)
		if err != nil {
			log.Printf("%s (%s): %s", err.Error(), rawPath, rawPath)
			continue
		}

		rawMatchPattern, assertOk := rawItemMappingMap["matchPattern"].(string)
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemMappingMatchPattern: %s", rawPath)
			continue
		}
		matchPattern, err := valueObject.NewMappingMatchPattern(rawMatchPattern)
		if err != nil {
			log.Printf("%s (%s): %s", err.Error(), rawPath, rawMatchPattern)
			continue
		}

		rawTargetType, assertOk := rawItemMappingMap["targetType"].(string)
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemMappingTargetType: %s", rawPath)
			continue
		}
		targetType, err := valueObject.NewMappingTargetType(rawTargetType)
		if err != nil {
			log.Printf("%s (%s): %s", err.Error(), rawPath, rawMatchPattern)
			continue
		}

		var targetSvcNamePtr *valueObject.ServiceName
		if rawItemMappingMap["targetServiceName"] != nil {
			rawTargetSvcName, assertOk := rawItemMappingMap["targetServiceName"].(string)
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemMappingTargetServiceName: %s",
					rawPath,
				)
				continue
			}
			targetSvcName, err := valueObject.NewServiceName(rawTargetSvcName)
			if err != nil {
				log.Printf("%s (%s): %s", err.Error(), rawPath, rawTargetSvcName)
				continue
			}
			targetSvcNamePtr = &targetSvcName
		}

		var targetUrlPtr *valueObject.Url
		if rawItemMappingMap["targetUrl"] != nil {
			rawTargetUrl, assertOk := rawItemMappingMap["targetUrl"].(string)
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemMappingTargetUrl: %s", rawPath,
				)
				continue
			}
			targetUrl, err := valueObject.NewUrl(rawTargetUrl)
			if err != nil {
				log.Printf("%s (%s): %s", err.Error(), rawPath, rawTargetUrl)
				continue
			}
			targetUrlPtr = &targetUrl
		}

		var targetHttpResCodePtr *valueObject.HttpResponseCode
		if rawItemMappingMap["targetHttpResponseCode"] != nil {
			rawTargetHttpResCode, assertOk := rawItemMappingMap["targetHttpResponseCode"].(string)
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemMappingTargetHttpResponseCode: %s",
					rawPath,
				)
				continue
			}
			targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
				rawTargetHttpResCode,
			)
			if err != nil {
				log.Printf("%s (%s): %s", err.Error(), rawPath, rawTargetHttpResCode)
				continue
			}
			targetHttpResCodePtr = &targetHttpResponseCode
		}

		var targetInlineHtmlContentPtr *valueObject.InlineHtmlContent
		if rawItemMappingMap["targetInlineHtmlContent"] != nil {
			rawTargetInlineHtmlContent, assertOk := rawItemMappingMap["targetInlineHtmlContent"].(string)
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemMappingTargetInlinteHtmlContent: %s",
					rawPath,
				)
				continue
			}
			targetInlineHtmlContent, err := valueObject.NewInlineHtmlContent(
				rawTargetInlineHtmlContent,
			)
			if err != nil {
				log.Printf(
					"%s (%s): %s", err.Error(), rawPath, rawTargetInlineHtmlContent,
				)
				continue
			}
			targetInlineHtmlContentPtr = &targetInlineHtmlContent
		}

		itemMapping := valueObject.NewMarketplaceItemMapping(
			path,
			matchPattern,
			targetType,
			targetSvcNamePtr,
			targetUrlPtr,
			targetHttpResCodePtr,
			targetInlineHtmlContentPtr,
		)
		itemMappings = append(itemMappings, itemMapping)
	}

	return itemMappings, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemDataFields(
	catalogItemDataFieldsMap interface{},
) ([]valueObject.MarketplaceCatalogItemDataField, error) {
	itemDataFields := []valueObject.MarketplaceCatalogItemDataField{}

	rawItemDataFields, assertOk := catalogItemDataFieldsMap.([]interface{})
	if !assertOk {
		return itemDataFields, errors.New("InvalidMarketplaceCatalogItemDataFields")
	}

	for _, rawItemDataField := range rawItemDataFields {
		rawItemDataFieldMap, assertOk := rawItemDataField.(map[string]interface{})
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemDataField: %+v", rawItemDataField)
			continue
		}

		rawKey, assertOk := rawItemDataFieldMap["key"].(string)
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemDataFieldKey: %s", rawKey)
			continue
		}
		key, err := valueObject.NewDataFieldKey(rawKey)
		if err != nil {
			log.Printf("%s (%s): %s", err.Error(), rawKey, rawKey)
			continue
		}

		isRequired := false
		if rawItemDataFieldMap["isRequired"] != nil {
			rawIsRequired, assertOk := rawItemDataFieldMap["isRequired"].(bool)
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemDataFieldIsRequired: %s", rawKey,
				)
				continue
			}
			isRequired = rawIsRequired
		}

		var defaultValuePtr *valueObject.DataFieldValue
		if rawItemDataFieldMap["defaultValue"] != nil {
			rawDefaultValue, assertOk := rawItemDataFieldMap["defaultValue"].(string)
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemDataFieldDefaultValue: %s", rawKey,
				)
				continue
			}
			defaultValue, err := valueObject.NewDataFieldValue(rawDefaultValue)
			if err != nil {
				log.Printf("%s (%s): %s", err.Error(), rawKey, rawDefaultValue)
				continue
			}
			defaultValuePtr = &defaultValue
		}

		if !isRequired && defaultValuePtr == nil {
			return itemDataFields, errors.New(
				"MarketplaceCatalogDataFieldWithoutDefaultValue: " +
					key.String(),
			)
		}

		itemDataField, err := valueObject.NewMarketplaceCatalogItemDataField(
			key,
			defaultValuePtr,
			isRequired,
		)
		if err != nil {
			log.Printf("%s (%s)", err.Error(), rawKey)
			continue
		}
		itemDataFields = append(itemDataFields, itemDataField)
	}

	return itemDataFields, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemCmdSteps(
	catalogItemCmdStepsMap interface{},
) ([]valueObject.MarketplaceItemCmdStep, error) {
	itemCmdSteps := []valueObject.MarketplaceItemCmdStep{}

	rawItemCmdSteps, assertOk := catalogItemCmdStepsMap.([]interface{})
	if !assertOk {
		return itemCmdSteps, errors.New("InvalidMarketplaceCatalogItemCmdSteps")
	}

	for _, rawItemCmdStep := range rawItemCmdSteps {
		itemCmdStep, err := valueObject.NewMarketplaceItemCmdStep(
			rawItemCmdStep.(string),
		)
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawItemCmdStep)
			continue
		}

		itemCmdSteps = append(itemCmdSteps, itemCmdStep)
	}

	return itemCmdSteps, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemScreenshotUrls(
	catalogItemUrlsMap interface{},
) ([]valueObject.Url, error) {
	itemUrls := []valueObject.Url{}

	rawItemUrls, assertOk := catalogItemUrlsMap.([]interface{})
	if !assertOk {
		return itemUrls, errors.New("InvalidMarketplaceCatalogItemUrls")
	}

	for _, rawItemUrl := range rawItemUrls {
		itemUrl, err := valueObject.NewUrl(rawItemUrl.(string))
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawItemUrl)
			continue
		}

		itemUrls = append(itemUrls, itemUrl)
	}

	return itemUrls, nil
}

func (repo *MarketplaceQueryRepo) catalogItemFactory(
	catalogItemFilePath valueObject.UnixFilePath,
) (entity.MarketplaceCatalogItem, error) {
	var catalogItem entity.MarketplaceCatalogItem

	catalogItemMap, err := repo.getCatalogItemMapFromFilePath(catalogItemFilePath)
	if err != nil {
		return catalogItem, err
	}

	rawCatalogItemName, assertOk := catalogItemMap["name"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemName")
	}
	catalogItemName, err := valueObject.NewMarketplaceItemName(rawCatalogItemName)
	if err != nil {
		return catalogItem, err
	}

	rawCatalogItemType, assertOk := catalogItemMap["type"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemType")
	}
	catalogItemType, err := valueObject.NewMarketplaceItemType(rawCatalogItemType)
	if err != nil {
		return catalogItem, err
	}

	rawCatalogItemDescription, assertOk := catalogItemMap["description"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemDescription")
	}
	catalogItemDescription, err := valueObject.NewMarketplaceItemDescription(rawCatalogItemDescription)
	if err != nil {
		return catalogItem, err
	}

	catalogItemSvcNames, err := repo.parseCatalogItemServiceNames(
		catalogItemMap["serviceNames"],
	)
	if err != nil {
		return catalogItem, err
	}

	catalogItemMappings, err := repo.parseCatalogItemMappings(
		catalogItemMap["mappings"],
	)
	if err != nil {
		return catalogItem, err
	}

	catalogItemDataFields, err := repo.parseCatalogItemDataFields(
		catalogItemMap["dataFields"],
	)
	if err != nil {
		return catalogItem, err
	}

	catalogItemCmdSteps, err := repo.parseCatalogItemCmdSteps(
		catalogItemMap["cmdSteps"],
	)
	if err != nil {
		return catalogItem, err
	}

	catalogEstimatedSizeBytes, err := valueObject.NewByte(
		catalogItemMap["estimatedSizeBytes"],
	)
	if err != nil {
		return catalogItem, err
	}

	rawCatalogItemAvatarUrl, assertOk := catalogItemMap["avatarUrl"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemAvatarUrl")
	}
	catalogItemAvatarUrl, err := valueObject.NewUrl(rawCatalogItemAvatarUrl)
	if err != nil {
		return catalogItem, err
	}

	catalogItemScreenshotUrls, err := repo.parseCatalogItemScreenshotUrls(
		catalogItemMap["screenshotUrls"],
	)
	if err != nil {
		return catalogItem, err
	}

	var placeholderCatalogItemId valueObject.MarketplaceCatalogItemId
	return entity.NewMarketplaceCatalogItem(
		placeholderCatalogItemId,
		catalogItemName,
		catalogItemType,
		catalogItemDescription,
		catalogItemSvcNames,
		catalogItemMappings,
		catalogItemDataFields,
		catalogItemCmdSteps,
		catalogEstimatedSizeBytes,
		catalogItemAvatarUrl,
		catalogItemScreenshotUrls,
	), nil
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

		catalogItem, err := repo.catalogItemFactory(
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
		catalogItemId, _ := valueObject.NewMarketplaceCatalogItemId(catalogItemIdInt)
		catalogItem.Id = catalogItemId

		catalogItems = append(catalogItems, catalogItem)
	}

	return catalogItems, nil
}

func (repo *MarketplaceQueryRepo) GetCatalogItemById(
	id valueObject.MarketplaceCatalogItemId,
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

	return catalogItem, errors.New("CatalogItemNotFound")
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

func (repo *MarketplaceQueryRepo) GetInstalledItemById(
	id valueObject.MarketplaceInstalledItemId,
) (entity.MarketplaceInstalledItem, error) {
	var installedItem entity.MarketplaceInstalledItem

	installedItems, err := repo.GetInstalledItems()
	if err != nil {
		return installedItem, err
	}

	for _, installedItem := range installedItems {
		if installedItem.Id.Get() != id.Get() {
			continue
		}

		return installedItem, nil
	}

	return installedItem, errors.New("InstalledItemNotFound")
}
