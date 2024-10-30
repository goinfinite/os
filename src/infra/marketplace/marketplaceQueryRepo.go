package marketplaceInfra

import (
	"errors"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
)

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

func (repo *MarketplaceQueryRepo) parseCatalogItemServices(
	catalogItemServices interface{},
) (serviceNamesWithVersions []valueObject.ServiceNameWithVersion, err error) {
	rawServices, assertOk := catalogItemServices.([]interface{})
	if !assertOk {
		return serviceNamesWithVersions, errors.New("InvalidCatalogItemServices")
	}

	for _, rawService := range rawServices {
		serviceNameWithVersion, err := valueObject.NewServiceNameWithVersionFromString(
			rawService,
		)
		if err != nil {
			slog.Error(err.Error(), slog.Any("serviceNameWithVersion", rawService))
			continue
		}

		serviceNamesWithVersions = append(
			serviceNamesWithVersions, serviceNameWithVersion,
		)
	}

	return serviceNamesWithVersions, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemMappings(
	catalogItemMappingsMap interface{},
) (itemMappings []valueObject.MarketplaceItemMapping, err error) {
	rawItemMappings, assertOk := catalogItemMappingsMap.([]interface{})
	if !assertOk {
		return itemMappings, errors.New("InvalidMarketplaceCatalogItemMappings")
	}

	for mappingIndex, rawItemMapping := range rawItemMappings {
		rawItemMappingMap, assertOk := rawItemMapping.(map[string]interface{})
		if !assertOk {
			slog.Error(
				"InvalidMarketplaceCatalogItemMapping", slog.Int("index", mappingIndex),
			)
			continue
		}

		path, err := valueObject.NewMappingPath(rawItemMappingMap["path"])
		if err != nil {
			slog.Error(err.Error(), slog.Int("index", mappingIndex))
			continue
		}

		matchPattern, err := valueObject.NewMappingMatchPattern(
			rawItemMappingMap["matchPattern"],
		)
		if err != nil {
			slog.Error(err.Error(), slog.Int("index", mappingIndex))
			continue
		}

		targetType, err := valueObject.NewMappingTargetType(
			rawItemMappingMap["targetType"],
		)
		if err != nil {
			slog.Error(err.Error(), slog.Int("index", mappingIndex))
			continue
		}

		var targetValuePtr *valueObject.MappingTargetValue
		if rawItemMappingMap["targetValue"] != nil {
			targetValue, err := valueObject.NewMappingTargetValue(
				rawItemMappingMap["targetValue"], targetType,
			)
			if err != nil {
				slog.Error(err.Error(), slog.Int("index", mappingIndex))
				continue
			}
			targetValuePtr = &targetValue
		}

		var targetHttpResponseCodePtr *valueObject.HttpResponseCode
		if rawItemMappingMap["targetHttpResponseCode"] != nil {
			targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
				rawItemMappingMap["targetHttpResponseCode"],
			)
			if err != nil {
				slog.Error(err.Error(), slog.Int("index", mappingIndex))
				continue
			}
			targetHttpResponseCodePtr = &targetHttpResponseCode
		}

		itemMapping := valueObject.NewMarketplaceItemMapping(
			path, matchPattern, targetType, targetValuePtr, targetHttpResponseCodePtr,
		)
		itemMappings = append(itemMappings, itemMapping)
	}

	return itemMappings, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemDataFields(
	catalogItemDataFieldsMap interface{},
) (itemDataFields []valueObject.MarketplaceCatalogItemDataField, err error) {
	rawItemDataFields, assertOk := catalogItemDataFieldsMap.([]interface{})
	if !assertOk {
		return itemDataFields, errors.New("InvalidMarketplaceCatalogItemDataFields")
	}

	for _, rawItemDataField := range rawItemDataFields {
		rawItemDataFieldMap, assertOk := rawItemDataField.(map[string]interface{})
		if !assertOk {
			slog.Error(
				"InvalidMarketplaceCatalogItemDataField",
				slog.Any("itemDataField", rawItemDataField),
			)
			continue
		}

		rawKey := rawItemDataFieldMap["name"]
		key, err := valueObject.NewDataFieldName(rawKey)
		if err != nil {
			slog.Error(err.Error(), slog.Any("key", rawKey))
			continue
		}

		rawLabel := rawItemDataFieldMap["label"]
		label, err := valueObject.NewDataFieldLabel(rawLabel)
		if err != nil {
			slog.Error(
				err.Error(), slog.Any("key", rawKey),
				slog.Any("label", rawLabel),
			)
			continue
		}

		rawFieldType := rawItemDataFieldMap["type"]
		fieldType, err := valueObject.NewDataFieldType(rawFieldType)
		if err != nil {
			slog.Error(
				err.Error(), slog.Any("key", rawKey),
				slog.Any("type", rawFieldType),
			)
			continue
		}

		isRequired := false
		if rawItemDataFieldMap["isRequired"] != nil {
			rawIsRequired, err := voHelper.InterfaceToBool(
				rawItemDataFieldMap["isRequired"],
			)
			if err != nil {
				slog.Error(
					"InvalidMarketplaceCatalogItemDataFieldIsRequired",
					slog.Any("err", err), slog.Any("key", rawKey),
					slog.Bool("isRequired", rawIsRequired),
				)
				continue
			}
			isRequired = rawIsRequired
		}

		var defaultValuePtr *valueObject.DataFieldValue
		if rawItemDataFieldMap["defaultValue"] != nil {
			rawDefaultValue := rawItemDataFieldMap["defaultValue"]
			defaultValue, err := valueObject.NewDataFieldValue(rawDefaultValue)
			if err != nil {
				slog.Error(
					err.Error(), slog.Any("key", rawKey),
					slog.Any("defaultValue", rawDefaultValue),
				)
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
				slog.Error(
					"InvalidMarketplaceCatalogItemDataFieldOptions",
					slog.Any("key", rawKey), slog.Any("options", rawOptions),
				)
				continue
			}

			for _, rawOption := range rawOptions {
				option, err := valueObject.NewDataFieldValue(rawOption)
				if err != nil {
					slog.Error(
						err.Error(), slog.Any("key", rawKey),
						slog.Any("options", rawOption),
					)
					continue
				}
				options = append(options, option)
			}
		}

		itemDataField, err := valueObject.NewMarketplaceCatalogItemDataField(
			key, label, fieldType, defaultValuePtr, options, isRequired,
		)
		if err != nil {
			slog.Error(err.Error(), slog.Any("key", rawKey))
			continue
		}
		itemDataFields = append(itemDataFields, itemDataField)
	}

	return itemDataFields, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemCmdSteps(
	catalogItemCmdStepsMap interface{},
) (itemCmdSteps []valueObject.UnixCommand, err error) {
	rawItemCmdSteps, assertOk := catalogItemCmdStepsMap.([]interface{})
	if !assertOk {
		return itemCmdSteps, errors.New("InvalidMarketplaceCatalogItemCmdSteps")
	}

	for _, rawItemCmdStep := range rawItemCmdSteps {
		itemCmdStep, err := valueObject.NewUnixCommand(rawItemCmdStep)
		if err != nil {
			slog.Error(err.Error(), slog.Any("cmdStep", rawItemCmdStep))
			continue
		}

		itemCmdSteps = append(itemCmdSteps, itemCmdStep)
	}

	return itemCmdSteps, nil
}

func (repo *MarketplaceQueryRepo) parseCatalogItemUninstallFileNames(
	catalogItemUninstallFileNames interface{},
) (itemUninstallFileNames []valueObject.UnixFileName, err error) {
	rawItemUninstallFileNames, assertOk := catalogItemUninstallFileNames.([]interface{})
	if !assertOk {
		return itemUninstallFileNames, errors.New(
			"InvalidMarketplaceCatalogItemUninstallFileNames",
		)
	}

	for _, rawItemUninstallFileName := range rawItemUninstallFileNames {
		itemUninstallUninstallFileNames, err := valueObject.NewUnixFileName(
			rawItemUninstallFileName,
		)
		if err != nil {
			slog.Error(err.Error(), slog.Any("fileName", rawItemUninstallFileName))
			continue
		}

		itemUninstallFileNames = append(
			itemUninstallFileNames, itemUninstallUninstallFileNames,
		)
	}

	return itemUninstallFileNames, nil
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
		itemUrl, err := valueObject.NewUrl(rawItemUrl)
		if err != nil {
			slog.Error(err.Error(), slog.Any("url", rawItemUrl))
			continue
		}

		itemUrls = append(itemUrls, itemUrl)
	}

	return itemUrls, nil
}

func (repo *MarketplaceQueryRepo) catalogItemFactory(
	catalogItemFilePath valueObject.UnixFilePath,
) (catalogItem entity.MarketplaceCatalogItem, err error) {
	itemMap, err := infraHelper.FileSerializedDataToMap(catalogItemFilePath)
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

	itemName, err := valueObject.NewMarketplaceItemName(itemMap["name"])
	if err != nil {
		return catalogItem, err
	}

	itemType, err := valueObject.NewMarketplaceItemType(itemMap["type"])
	if err != nil {
		return catalogItem, err
	}

	itemDescription, err := valueObject.NewMarketplaceItemDescription(
		itemMap["description"],
	)
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

	itemInstallCmdSteps := []valueObject.UnixCommand{}
	if itemMap["installCmdSteps"] != nil {
		itemInstallCmdSteps, err = repo.parseCatalogItemCmdSteps(itemMap["installCmdSteps"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemUninstallCmdSteps := []valueObject.UnixCommand{}
	if itemMap["uninstallCmdSteps"] != nil {
		itemUninstallCmdSteps, err = repo.parseCatalogItemCmdSteps(itemMap["uninstallCmdSteps"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemUninstallFileNames := []valueObject.UnixFileName{}
	if itemMap["uninstallFileNames"] != nil {
		itemUninstallFileNames, err = repo.parseCatalogItemUninstallFileNames(
			itemMap["uninstallFileNames"],
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

	itemAvatarUrl, err := valueObject.NewUrl(itemMap["avatarUrl"])
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
		itemUninstallFileNames,
		estimatedSizeBytes,
		itemAvatarUrl,
		itemScreenshotUrls,
	), nil
}

func (repo *MarketplaceQueryRepo) ReadCatalogItems() (
	[]entity.MarketplaceCatalogItem, error,
) {
	catalogItems := []entity.MarketplaceCatalogItem{}

	_, err := os.Stat(infraEnvs.MarketplaceItemsDir)
	if err != nil {
		marketplaceCmdRepo := NewMarketplaceCmdRepo(repo.persistentDbSvc)
		err = marketplaceCmdRepo.RefreshItems()
		if err != nil {
			return catalogItems, errors.New(
				"RefreshMarketplaceItemsError: " + err.Error(),
			)
		}
	}

	rawCatalogFilesList, err := infraHelper.RunCmdWithSubShell(
		"find " + infraEnvs.MarketplaceItemsDir + " -type f " +
			"\\( -name '*.json' -o -name '*.yaml' -o -name '*.yml' \\) " +
			"-not -path '*/.*' -not -name '.*'",
	)
	if err != nil {
		return catalogItems, errors.New("ReadMarketplaceFilesError: " + err.Error())
	}

	if len(rawCatalogFilesList) == 0 {
		return catalogItems, errors.New("NoMarketplaceFilesFound")
	}

	rawCatalogFilesListParts := strings.Split(rawCatalogFilesList, "\n")
	if len(rawCatalogFilesListParts) == 0 {
		return catalogItems, errors.New("NoMarketplaceFilesFound")
	}

	catalogItemsIdsMap := map[uint16]struct{}{}
	for _, rawFilePath := range rawCatalogFilesListParts {
		itemFilePath, err := valueObject.NewUnixFilePath(rawFilePath)
		if err != nil {
			slog.Error(err.Error(), slog.String("filePath", rawFilePath))
			continue
		}

		catalogItem, err := repo.catalogItemFactory(itemFilePath)
		if err != nil {
			slog.Error(
				"CatalogMarketplaceItemFactoryError",
				slog.String("filePath", itemFilePath.String()), slog.Any("err", err),
			)
			continue
		}

		itemIdUint16 := catalogItem.Id.Uint16()
		_, idAlreadyUsed := catalogItemsIdsMap[itemIdUint16]
		if idAlreadyUsed {
			catalogItem.Id, _ = valueObject.NewMarketplaceItemId(0)
		}

		catalogItems = append(catalogItems, catalogItem)

		if catalogItem.Id.Uint16() != 0 {
			catalogItemsIdsMap[itemIdUint16] = struct{}{}
		}
	}

	itemsIdsSlice := []uint16{}
	for itemId := range catalogItemsIdsMap {
		itemsIdsSlice = append(itemsIdsSlice, itemId)
	}
	slices.Sort(itemsIdsSlice)

	if len(itemsIdsSlice) == 0 {
		itemsIdsSlice = append(itemsIdsSlice, 0)
	}

	for itemIndex, catalogItem := range catalogItems {
		if catalogItem.Id.Uint16() != 0 {
			continue
		}

		lastIdUsed := itemsIdsSlice[len(itemsIdsSlice)-1]
		nextAvailableId, err := valueObject.NewMarketplaceItemId(lastIdUsed + 1)
		if err != nil {
			slog.Error(
				"CreateNewCatalogMarketplaceItemIdError",
				slog.String("itemName", catalogItem.Name.String()),
				slog.Any("err", err),
			)
			continue
		}

		catalogItems[itemIndex].Id = nextAvailableId
		itemsIdsSlice = append(itemsIdsSlice, nextAvailableId.Uint16())
	}

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
		if catalogItem.Id.Uint16() != catalogId.Uint16() {
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
			slog.Error(
				"MarketplaceInstalledItemModelToEntityError", slog.Any("error", err),
			)
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
		Where("id = ?", installedId.Uint16()).
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
