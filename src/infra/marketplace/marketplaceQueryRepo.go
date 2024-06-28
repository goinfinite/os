package marketplaceInfra

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"slices"
	"sort"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
	"golang.org/x/exp/maps"
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
) (catalogItemMap map[string]interface{}, err error) {
	itemFileHandler, err := assets.Open(catalogItemFilePath.String())
	if err != nil {
		return catalogItemMap, err
	}
	defer itemFileHandler.Close()

	itemFileExt, err := catalogItemFilePath.GetFileExtension()
	if err != nil {
		return catalogItemMap, err
	}

	isYamlFile := itemFileExt == "yml" || itemFileExt == "yaml"
	if isYamlFile {
		itemYamlDecoder := yaml.NewDecoder(itemFileHandler)
		err = itemYamlDecoder.Decode(&catalogItemMap)
		if err != nil {
			return catalogItemMap, err
		}

		return catalogItemMap, nil
	}

	itemJsonDecoder := json.NewDecoder(itemFileHandler)
	err = itemJsonDecoder.Decode(&catalogItemMap)
	if err != nil {
		return catalogItemMap, err
	}

	return catalogItemMap, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemServices(
	catalogItemServices interface{},
) (serviceNamesWithVersions []valueObject.ServiceNameWithVersion, err error) {
	rawServices, assertOk := catalogItemServices.([]interface{})
	if !assertOk {
		return serviceNamesWithVersions, errors.New("InvalidCatalogItemServices")
	}

	for _, rawService := range rawServices {
		rawServiceNameWithVersion, assertOk := rawService.(string)
		if !assertOk {
			log.Printf("InvalidCatalogItemService: %s", rawService)
			continue
		}

		serviceNameWithVersion, err := valueObject.NewServiceNameWithVersionFromString(
			rawServiceNameWithVersion,
		)
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawServiceNameWithVersion)
			continue
		}

		serviceNamesWithVersions = append(serviceNamesWithVersions, serviceNameWithVersion)
	}

	return serviceNamesWithVersions, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemMappings(
	catalogItemMappingsMap interface{},
) (itemMappings []valueObject.MarketplaceItemMapping, err error) {
	if catalogItemMappingsMap == nil {
		return itemMappings, nil
	}

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

		var targetValuePtr *valueObject.MappingTargetValue
		if rawItemMappingMap["targetValue"] != nil {
			rawTargetValue, assertOk := rawItemMappingMap["targetValue"].(string)
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemMappingTargetValue: %s",
					rawPath,
				)
				continue
			}

			targetValue, err := valueObject.NewMappingTargetValue(
				rawTargetValue, targetType,
			)
			if err != nil {
				log.Printf("%s (%s): %s", err.Error(), rawPath, rawMatchPattern)
				continue
			}
			targetValuePtr = &targetValue
		}

		var targetHttpResponseCodePtr *valueObject.HttpResponseCode
		if rawItemMappingMap["targetHttpResponseCode"] != nil {
			rawTargetHttpResponseCode, assertOk := rawItemMappingMap["targetHttpResponseCode"].(string)
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemMappingTargetHttpResponseCode: %s",
					rawPath,
				)
				continue
			}

			targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
				rawTargetHttpResponseCode,
			)
			if err != nil {
				log.Printf("%s (%s): %s", err.Error(), rawPath, rawMatchPattern)
				continue
			}
			targetHttpResponseCodePtr = &targetHttpResponseCode
		}

		itemMapping := valueObject.NewMarketplaceItemMapping(
			path,
			matchPattern,
			targetType,
			targetValuePtr,
			targetHttpResponseCodePtr,
		)
		itemMappings = append(itemMappings, itemMapping)
	}

	return itemMappings, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemDataFields(
	catalogItemDataFieldsMap interface{},
) (itemDataFields []valueObject.MarketplaceCatalogItemDataField, err error) {
	if catalogItemDataFieldsMap == nil {
		return itemDataFields, nil
	}

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

		rawKey, assertOk := rawItemDataFieldMap["name"].(string)
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemDataFieldKey: %s", rawKey)
			continue
		}
		key, err := valueObject.NewDataFieldName(rawKey)
		if err != nil {
			log.Printf("%s (%s): %s", err.Error(), rawKey, rawKey)
			continue
		}

		rawLabel, assertOk := rawItemDataFieldMap["label"].(string)
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemDataFieldLabel: %s", rawKey)
			continue
		}
		label, err := valueObject.NewDataFieldLabel(rawLabel)
		if err != nil {
			log.Printf("%s (%s): %s", err.Error(), rawKey, rawLabel)
			continue
		}

		rawFieldType, assertOk := rawItemDataFieldMap["type"].(string)
		if !assertOk {
			log.Printf("InvalidMarketplaceCatalogItemDataFieldType: %s", rawKey)
			continue
		}
		fieldType, err := valueObject.NewDataFieldType(rawFieldType)
		if err != nil {
			log.Printf("%s (%s): %s", err.Error(), rawKey, rawFieldType)
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

		options := []valueObject.DataFieldValue{}
		if rawItemDataFieldMap["options"] != nil {
			rawOptions, assertOk := rawItemDataFieldMap["options"].([]interface{})
			if !assertOk {
				log.Printf(
					"InvalidMarketplaceCatalogItemDataFieldOptions: %s", rawKey,
				)
				continue
			}

			for _, rawOption := range rawOptions {
				option, err := valueObject.NewDataFieldValue(rawOption)
				if err != nil {
					log.Printf("%s (%s): %s", err.Error(), rawKey, rawOption)
					continue
				}
				options = append(options, option)
			}
		}

		itemDataField, err := valueObject.NewMarketplaceCatalogItemDataField(
			key,
			label,
			fieldType,
			defaultValuePtr,
			options,
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
) (itemCmdSteps []valueObject.MarketplaceItemCmdStep, err error) {
	if catalogItemCmdStepsMap == nil {
		return itemCmdSteps, nil
	}

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

func (repo *MarketplaceQueryRepo) parseCatalogItemUninstallFilesToRemove(
	catalogItemUninstallFilesToRemove interface{},
) (itemUninstallFilesToRemove []valueObject.UnixFileName, err error) {
	if catalogItemUninstallFilesToRemove == nil {
		return itemUninstallFilesToRemove, nil
	}

	rawItemFilesToDelete, assertOk := catalogItemUninstallFilesToRemove.([]string)
	if !assertOk {
		return itemUninstallFilesToRemove, errors.New(
			"InvalidMarketplaceCatalogItemFilesToDelete",
		)
	}

	for _, rawItemFileToDelete := range rawItemFilesToDelete {
		itemUninstallFileToRemove, err := valueObject.NewUnixFileName(rawItemFileToDelete)
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawItemFileToDelete)
			continue
		}

		itemUninstallFilesToRemove = append(
			itemUninstallFilesToRemove, itemUninstallFileToRemove,
		)
	}

	return itemUninstallFilesToRemove, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemScreenshotUrls(
	catalogItemUrlsMap interface{},
) (itemUrls []valueObject.Url, err error) {
	if catalogItemUrlsMap == nil {
		return itemUrls, nil
	}

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
) (catalogItem entity.MarketplaceCatalogItem, err error) {
	itemMap, err := repo.getCatalogItemMapFromFilePath(catalogItemFilePath)
	if err != nil {
		return catalogItem, err
	}

	itemId, _ := valueObject.NewMarketplaceItemId(0)
	rawItemId, exists := itemMap["id"]
	if exists {
		itemId, _ = valueObject.NewMarketplaceItemId(rawItemId)
	}

	itemSlugs := []valueObject.MarketplaceItemSlug{}
	if itemMap["slugs"] != nil {
		rawItemSlugs, assertOk := itemMap["slugs"].([]interface{})
		if !assertOk {
			return catalogItem, errors.New("InvalidMarketplaceItemSlugs")
		}
		for _, rawItemSlug := range rawItemSlugs {
			itemSlug, err := valueObject.NewMarketplaceItemSlug(rawItemSlug)
			if err != nil {
				return catalogItem, err
			}
			itemSlugs = append(itemSlugs, itemSlug)
		}
	}

	rawItemName, assertOk := itemMap["name"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceItemName")
	}
	itemName, err := valueObject.NewMarketplaceItemName(rawItemName)
	if err != nil {
		return catalogItem, err
	}

	rawItemType, assertOk := itemMap["type"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceItemType")
	}
	itemType, err := valueObject.NewMarketplaceItemType(rawItemType)
	if err != nil {
		return catalogItem, err
	}

	rawItemDescription, assertOk := itemMap["description"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceItemDescription")
	}
	itemDescription, err := valueObject.NewMarketplaceItemDescription(rawItemDescription)
	if err != nil {
		return catalogItem, err
	}

	itemServices := []valueObject.ServiceNameWithVersion{}
	if itemMap["services"] != nil {
		itemServices, err = repo.parseCatalogItemServices(itemMap["services"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemMappings := []valueObject.MarketplaceItemMapping{}
	if itemMap["mappings"] != nil {
		itemMappings, err = repo.parseCatalogItemMappings(itemMap["mappings"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemDataFields := []valueObject.MarketplaceCatalogItemDataField{}
	if itemMap["dataFields"] != nil {
		itemDataFields, err = repo.parseCatalogItemDataFields(itemMap["dataFields"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemInstallCmdSteps := []valueObject.MarketplaceItemCmdStep{}
	if itemMap["installCmdSteps"] != nil {
		itemInstallCmdSteps, err = repo.parseCatalogItemCmdSteps(itemMap["installCmdSteps"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemUninstallCmdSteps := []valueObject.MarketplaceItemCmdStep{}
	if itemMap["uninstallCmdSteps"] != nil {
		itemUninstallCmdSteps, err = repo.parseCatalogItemCmdSteps(itemMap["uninstallCmdSteps"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemUninstallFilesToRemove := []valueObject.UnixFileName{}
	if itemMap["filesToRemove"] != nil {
		itemUninstallFilesToRemove, err = repo.parseCatalogItemUninstallFilesToRemove(
			itemMap["filesToRemove"],
		)
		if err != nil {
			return catalogItem, err
		}
	}

	estimatedSizeBytes := valueObject.Byte(1000000000)
	if itemMap["estimatedSizeBytes"] == nil {
		estimatedSizeBytes, err = valueObject.NewByte(itemMap["estimatedSizeBytes"])
		if err != nil {
			return catalogItem, err
		}
	}

	rawItemAvatarUrl, assertOk := itemMap["avatarUrl"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceItemAvatarUrl")
	}
	itemAvatarUrl, err := valueObject.NewUrl(rawItemAvatarUrl)
	if err != nil {
		return catalogItem, err
	}

	itemScreenshotUrls := []valueObject.Url{}
	if itemMap["screenshotUrls"] != nil {
		itemScreenshotUrls, err = repo.parseCatalogItemScreenshotUrls(itemMap["screenshotUrls"])
		if err != nil {
			return catalogItem, err
		}
	}

	return entity.NewMarketplaceCatalogItem(
		itemId,
		itemSlugs,
		itemName,
		itemType,
		itemDescription,
		itemServices,
		itemMappings,
		itemDataFields,
		itemInstallCmdSteps,
		itemUninstallCmdSteps,
		itemUninstallFilesToRemove,
		estimatedSizeBytes,
		itemAvatarUrl,
		itemScreenshotUrls,
	), nil
}

func (repo *MarketplaceQueryRepo) ReadCatalogItems() (
	catalogItems []entity.MarketplaceCatalogItem, err error,
) {
	itemsFiles, err := fs.ReadDir(assets, "assets")
	if err != nil {
		return catalogItems, errors.New(
			"GetMarketplaceCatalogItemsFilesError: " + err.Error(),
		)
	}

	catalogItemsIdsMap := map[uint]interface{}{}
	for _, itemFileEntry := range itemsFiles {
		itemFileName := itemFileEntry.Name()

		itemFilePathStr := "assets/" + itemFileName
		itemFilePath, err := valueObject.NewUnixFilePath(itemFilePathStr)
		if err != nil {
			log.Printf("%s (%s): %s", err.Error(), itemFileName, itemFilePathStr)
			continue
		}

		catalogItem, err := repo.catalogItemFactory(itemFilePath)
		if err != nil {
			log.Printf(
				"ReadMarketplaceCatalogItemError (%s): %s", itemFileName, err.Error(),
			)
			continue
		}

		_, idAlreadyUsed := catalogItemsIdsMap[catalogItem.Id.Get()]
		if idAlreadyUsed {
			catalogItem.Id, _ = valueObject.NewMarketplaceItemId(0)
		}

		if catalogItem.Id.Get() != 0 {
			catalogItemsIdsMap[catalogItem.Id.Get()] = nil
		}

		catalogItems = append(catalogItems, catalogItem)
	}

	catalogItemsIds := maps.Keys(catalogItemsIdsMap)
	slices.Sort(catalogItemsIds)

	for itemIndex, catalogItem := range catalogItems {
		if catalogItem.Id.Get() != 0 {
			continue
		}

		lastIdUsed := catalogItemsIds[len(catalogItemsIds)-1]
		nextAvailableId, err := valueObject.NewMarketplaceItemId(lastIdUsed + 1)
		if err != nil {
			log.Printf(
				"GenerateNewMarketplaceItemIdError (%s): %s",
				catalogItem.Name.String(), err.Error(),
			)
			continue
		}

		catalogItems[itemIndex].Id = nextAvailableId
		catalogItemsIds = append(catalogItemsIds, nextAvailableId.Get())
	}

	sort.SliceStable(catalogItems, func(i, j int) bool {
		return catalogItems[i].Id.Get() < catalogItems[j].Id.Get()
	})

	return catalogItems, nil
}

func (repo *MarketplaceQueryRepo) ReadCatalogItemById(
	catalogId valueObject.MarketplaceItemId,
) (catalogItem entity.MarketplaceCatalogItem, err error) {
	catalogItems, err := repo.ReadCatalogItems()
	if err != nil {
		return catalogItem, err
	}

	for _, catalogItem := range catalogItems {
		if catalogItem.Id.Get() != catalogId.Get() {
			continue
		}

		return catalogItem, nil
	}

	return catalogItem, errors.New("CatalogItemNotFound")
}

func (repo *MarketplaceQueryRepo) ReadCatalogItemBySlug(
	slug valueObject.MarketplaceItemSlug,
) (catalogItem entity.MarketplaceCatalogItem, err error) {
	catalogItems, err := repo.ReadCatalogItems()
	if err != nil {
		return catalogItem, err
	}

	for _, catalogItem := range catalogItems {
		for _, catalogItemSlug := range catalogItem.Slugs {
			if catalogItemSlug.String() != slug.String() {
				continue
			}

			return catalogItem, nil
		}
	}

	return catalogItem, errors.New("CatalogItemNotFound")
}

func (repo *MarketplaceQueryRepo) ReadInstalledItems() (
	[]entity.MarketplaceInstalledItem, error,
) {
	entities := []entity.MarketplaceInstalledItem{}

	models := []dbModel.MarketplaceInstalledItem{}
	err := repo.persistentDbSvc.Handler.
		Model(&dbModel.MarketplaceInstalledItem{}).
		Preload("Mappings").
		Find(&models).Error
	if err != nil {
		return entities, errors.New("ReadDatabaseEntriesError")
	}

	for _, installedItemModel := range models {
		entity, err := installedItemModel.ToEntity()
		if err != nil {
			log.Printf("MarketplaceInstalledItemModelToEntityError: %s", err.Error())
			continue
		}

		entities = append(
			entities,
			entity,
		)
	}

	return entities, nil
}

func (repo *MarketplaceQueryRepo) ReadInstalledItemById(
	installedId valueObject.MarketplaceItemId,
) (entity entity.MarketplaceInstalledItem, err error) {
	var model dbModel.MarketplaceInstalledItem
	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.MarketplaceInstalledItem{}).
		Where("id = ?", installedId.Get()).
		Preload("Mappings").
		Find(&model).Error
	if err != nil {
		return entity, errors.New("ReadDatabaseEntryError")
	}

	entity, err = model.ToEntity()
	if err != nil {
		return entity, errors.New("ModelToEntityError")
	}

	return entity, nil
}
