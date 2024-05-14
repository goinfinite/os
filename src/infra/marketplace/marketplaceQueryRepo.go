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
) ([]valueObject.MarketplaceCatalogItemDataField, error) {
	itemDataFields := []valueObject.MarketplaceCatalogItemDataField{}

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
) ([]valueObject.MarketplaceItemCmdStep, error) {
	itemCmdSteps := []valueObject.MarketplaceItemCmdStep{}

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

func (repo *MarketplaceQueryRepo) parseCatalogItemScreenshotUrls(
	catalogItemUrlsMap interface{},
) ([]valueObject.Url, error) {
	itemUrls := []valueObject.Url{}

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

func (repo *MarketplaceQueryRepo) ReadCatalogItems() (
	[]entity.MarketplaceCatalogItem, error,
) {
	catalogItems := []entity.MarketplaceCatalogItem{}

	catalogItemFiles, err := fs.ReadDir(assets, "assets")
	if err != nil {
		return catalogItems, errors.New(
			"GetMarketplaceCatalogItemsFilesError: " + err.Error(),
		)
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

func (repo *MarketplaceQueryRepo) ReadCatalogItemById(
	catalogId valueObject.MarketplaceCatalogItemId,
) (entity.MarketplaceCatalogItem, error) {
	var catalogItem entity.MarketplaceCatalogItem

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

func (repo *MarketplaceQueryRepo) ReadInstalledItems() (
	[]entity.MarketplaceInstalledItem, error,
) {
	entities := []entity.MarketplaceInstalledItem{}

	models := []dbModel.MarketplaceInstalledItem{}
	err := repo.persistentDbSvc.Handler.
		Model(models).
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
	installedId valueObject.MarketplaceInstalledItemId,
) (entity entity.MarketplaceInstalledItem, err error) {
	query := dbModel.Mapping{
		ID: uint(installedId.Get()),
	}

	var model dbModel.MarketplaceInstalledItem
	err = repo.persistentDbSvc.Handler.
		Model(query).
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
