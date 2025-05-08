package marketplaceInfra

import (
	"errors"
	"log/slog"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
	tkInfra "github.com/goinfinite/tk/src/infra"
	"github.com/iancoleman/strcase"
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

func (repo *MarketplaceQueryRepo) catalogItemServicesFactory(
	catalogItemServices interface{},
) (serviceNamesWithVersions []valueObject.ServiceNameWithVersion, err error) {
	rawServiceNamesWithVersion, assertOk := catalogItemServices.([]interface{})
	if !assertOk {
		return serviceNamesWithVersions, errors.New(
			"InvalidCatalogItemServicesStructure",
		)
	}

	for _, rawServiceNameWithVersion := range rawServiceNamesWithVersion {
		serviceNameWithVersion, err := valueObject.NewServiceNameWithVersionFromString(
			rawServiceNameWithVersion,
		)
		if err != nil {
			slog.Debug(
				err.Error(),
				slog.Any("serviceNameWithVersion", rawServiceNameWithVersion),
			)
			continue
		}

		serviceNamesWithVersions = append(
			serviceNamesWithVersions, serviceNameWithVersion,
		)
	}

	return serviceNamesWithVersions, nil
}

func (repo *MarketplaceQueryRepo) catalogItemMappingsFactory(
	catalogItemMappingsMap interface{},
) (itemMappings []valueObject.MarketplaceItemMapping, err error) {
	rawItemMappings, assertOk := catalogItemMappingsMap.([]interface{})
	if !assertOk {
		return itemMappings, errors.New(
			"InvalidMarketplaceCatalogItemMappingsStructure",
		)
	}

	for mappingIndex, rawItemMapping := range rawItemMappings {
		rawItemMappingMap, assertOk := rawItemMapping.(map[string]interface{})
		if !assertOk {
			slog.Debug(
				"InvalidMarketplaceCatalogItemMapping", slog.Int("index", mappingIndex),
			)
			continue
		}

		path, err := valueObject.NewMappingPath(rawItemMappingMap["path"])
		if err != nil {
			slog.Debug(err.Error(), slog.Int("index", mappingIndex))
			continue
		}

		matchPattern, err := valueObject.NewMappingMatchPattern(
			rawItemMappingMap["matchPattern"],
		)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("index", mappingIndex))
			continue
		}

		targetType, err := valueObject.NewMappingTargetType(
			rawItemMappingMap["targetType"],
		)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("index", mappingIndex))
			continue
		}

		var targetValuePtr *valueObject.MappingTargetValue
		if rawItemMappingMap["targetValue"] != nil {
			targetValue, err := valueObject.NewMappingTargetValue(
				rawItemMappingMap["targetValue"], targetType,
			)
			if err != nil {
				slog.Debug(err.Error(), slog.Int("index", mappingIndex))
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
				slog.Debug(err.Error(), slog.Int("index", mappingIndex))
				continue
			}
			targetHttpResponseCodePtr = &targetHttpResponseCode
		}

		var shouldUpgradeInsecureRequestsPtr *bool
		if rawItemMappingMap["shouldUpgradeInsecureRequests"] != nil {
			shouldUpgradeInsecureRequests, err := tkVoUtil.InterfaceToBool(
				rawItemMappingMap["shouldUpgradeInsecureRequests"],
			)
			if err != nil {
				slog.Debug("ShouldUpgradeInsecureInvalidFormat", slog.Int("index", mappingIndex))
				shouldUpgradeInsecureRequests = false
			}
			shouldUpgradeInsecureRequestsPtr = &shouldUpgradeInsecureRequests
		}

		itemMapping := valueObject.NewMarketplaceItemMapping(
			path, matchPattern, targetType, targetValuePtr, targetHttpResponseCodePtr,
			shouldUpgradeInsecureRequestsPtr,
		)
		itemMappings = append(itemMappings, itemMapping)
	}

	return itemMappings, nil
}

func (repo *MarketplaceQueryRepo) specificTypeDataFieldValueGenerator(
	dataFieldSpecificType valueObject.DataFieldSpecificType,
) (valueObject.DataFieldValue, error) {
	synthesizer := tkInfra.Synthesizer{}

	var dummyValue string
	switch dataFieldSpecificType.String() {
	case "password":
		dummyValue = synthesizer.PasswordFactory(16, false)
	case "username":
		dummyValue = synthesizer.UsernameFactory()
	case "email":
		dummyValue = synthesizer.MailAddressFactory(nil)
	}

	return valueObject.NewDataFieldValue(dummyValue)
}

func (repo *MarketplaceQueryRepo) catalogItemDataFieldsFactory(
	catalogItemDataFieldsMap interface{},
) (itemDataFields []valueObject.MarketplaceCatalogItemDataField, err error) {
	rawItemDataFields, assertOk := catalogItemDataFieldsMap.([]interface{})
	if !assertOk {
		return itemDataFields, errors.New(
			"InvalidMarketplaceCatalogItemDataFieldsStructure",
		)
	}

	for _, rawItemDataField := range rawItemDataFields {
		rawItemDataFieldMap, assertOk := rawItemDataField.(map[string]interface{})
		if !assertOk {
			slog.Debug(
				"InvalidMarketplaceCatalogItemDataField",
				slog.Any("itemDataField", rawItemDataField),
			)
			continue
		}

		rawKey := rawItemDataFieldMap["name"]
		key, err := valueObject.NewDataFieldName(rawKey)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("key", rawKey))
			continue
		}

		rawLabel := rawItemDataFieldMap["label"]
		label, err := valueObject.NewDataFieldLabel(rawLabel)
		if err != nil {
			slog.Debug(
				err.Error(), slog.Any("key", rawKey),
				slog.Any("label", rawLabel),
			)
			continue
		}

		rawFieldType := rawItemDataFieldMap["type"]
		fieldType, err := valueObject.NewDataFieldType(rawFieldType)
		if err != nil {
			slog.Debug(
				err.Error(), slog.Any("key", rawKey),
				slog.Any("type", rawFieldType),
			)
			continue
		}

		var fieldSpecificTypePtr *valueObject.DataFieldSpecificType
		if rawItemDataFieldMap["specificType"] != nil {
			rawFieldSpecificType := rawItemDataFieldMap["specificType"]
			fieldSpecificType, err := valueObject.NewDataFieldSpecificType(
				rawFieldSpecificType,
			)
			if err != nil {
				slog.Debug(
					err.Error(), slog.Any("key", rawKey),
					slog.Any("specificType", rawFieldSpecificType),
				)
				continue
			}
			fieldSpecificTypePtr = &fieldSpecificType
		}

		var defaultValuePtr *valueObject.DataFieldValue
		if rawItemDataFieldMap["defaultValue"] != nil {
			rawDefaultValue := rawItemDataFieldMap["defaultValue"]
			defaultValue, err := valueObject.NewDataFieldValue(rawDefaultValue)
			if err != nil {
				slog.Debug(
					err.Error(), slog.Any("key", rawKey),
					slog.Any("defaultValue", rawDefaultValue),
				)
				continue
			}
			defaultValuePtr = &defaultValue
		}

		if fieldSpecificTypePtr != nil && defaultValuePtr == nil {
			defaultValue, err := repo.specificTypeDataFieldValueGenerator(
				*fieldSpecificTypePtr,
			)
			if err != nil {
				slog.Debug(
					err.Error(), slog.Any("key", rawKey),
					slog.Any("specificType", fieldSpecificTypePtr.String()),
				)
				continue
			}
			defaultValuePtr = &defaultValue
		}

		isRequired := false
		if rawItemDataFieldMap["isRequired"] != nil {
			rawIsRequired, err := tkVoUtil.InterfaceToBool(
				rawItemDataFieldMap["isRequired"],
			)
			if err != nil {
				slog.Debug(
					"InvalidMarketplaceCatalogItemDataFieldIsRequired",
					slog.String("err", err.Error()), slog.Any("key", rawKey),
					slog.Bool("isRequired", rawIsRequired),
				)
				continue
			}
			isRequired = rawIsRequired
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
				slog.Debug(
					"InvalidMarketplaceCatalogItemDataFieldOptions",
					slog.Any("key", rawKey), slog.Any("options", rawOptions),
				)
				continue
			}

			for _, rawOption := range rawOptions {
				option, err := valueObject.NewDataFieldValue(rawOption)
				if err != nil {
					slog.Debug(
						err.Error(), slog.Any("key", rawKey),
						slog.Any("options", rawOption),
					)
					continue
				}
				options = append(options, option)
			}
		}

		itemDataField, err := valueObject.NewMarketplaceCatalogItemDataField(
			key, label, fieldType, fieldSpecificTypePtr, defaultValuePtr, options,
			isRequired,
		)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("key", rawKey))
			continue
		}
		itemDataFields = append(itemDataFields, itemDataField)
	}

	return itemDataFields, nil
}

func (repo *MarketplaceQueryRepo) catalogItemCmdStepsFactory(
	catalogItemCmdStepsMap interface{},
) (itemCmdSteps []valueObject.UnixCommand, err error) {
	rawItemCmdSteps, assertOk := catalogItemCmdStepsMap.([]interface{})
	if !assertOk {
		return itemCmdSteps, errors.New(
			"InvalidMarketplaceCatalogItemCmdStepsStructure",
		)
	}

	for _, rawItemCmdStep := range rawItemCmdSteps {
		itemCmdStep, err := valueObject.NewUnixCommand(rawItemCmdStep)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("cmdStep", rawItemCmdStep))
			continue
		}

		itemCmdSteps = append(itemCmdSteps, itemCmdStep)
	}

	return itemCmdSteps, nil
}

func (repo *MarketplaceQueryRepo) catalogItemUninstallFileNamesFactory(
	catalogItemUninstallFileNames interface{},
) (itemUninstallFileNames []valueObject.UnixFileName, err error) {
	rawItemUninstallFileNames, assertOk := catalogItemUninstallFileNames.([]interface{})
	if !assertOk {
		return itemUninstallFileNames, errors.New(
			"InvalidMarketplaceCatalogItemUninstallFileNamesStructure",
		)
	}

	for _, rawItemUninstallFileName := range rawItemUninstallFileNames {
		itemUninstallUninstallFileNames, err := valueObject.NewUnixFileName(
			rawItemUninstallFileName,
		)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("fileName", rawItemUninstallFileName))
			continue
		}

		itemUninstallFileNames = append(
			itemUninstallFileNames, itemUninstallUninstallFileNames,
		)
	}

	return itemUninstallFileNames, nil
}

func (repo *MarketplaceQueryRepo) catalogItemScreenshotUrlsFactory(
	catalogItemUrlsMap interface{},
) (itemUrls []valueObject.Url, err error) {
	if catalogItemUrlsMap == nil {
		return itemUrls, nil
	}

	rawItemUrls, assertOk := catalogItemUrlsMap.([]interface{})
	if !assertOk {
		return itemUrls, errors.New(
			"InvalidMarketplaceCatalogItemUrlsStructure",
		)
	}

	for _, rawItemUrl := range rawItemUrls {
		itemUrl, err := valueObject.NewUrl(rawItemUrl)
		if err != nil {
			slog.Debug(err.Error(), slog.Any("url", rawItemUrl))
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

	requiredFields := []string{
		"name", "type", "description", "avatarUrl",
	}
	missingFields := []string{}
	for _, requiredField := range requiredFields {
		if _, exists := itemMap[requiredField]; !exists {
			missingFields = append(missingFields, requiredField)
		}
	}
	if len(missingFields) > 0 {
		return catalogItem, errors.New(
			"MissingItemFields: " + strings.Join(missingFields, ", "),
		)
	}

	itemManifestVersion, _ := valueObject.NewMarketplaceItemManifestVersion("v1")
	if itemMap["manifestVersion"] != nil {
		itemManifestVersion, err = valueObject.NewMarketplaceItemManifestVersion(
			itemMap["manifestVersion"],
		)
		if err != nil {
			return catalogItem, err
		}
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
			return catalogItem, errors.New("InvalidMarketplaceItemSlugsStructure")
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
		itemServices, err = repo.catalogItemServicesFactory(itemMap["services"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemMappings := []valueObject.MarketplaceItemMapping{}
	if itemMap["mappings"] != nil {
		itemMappings, err = repo.catalogItemMappingsFactory(itemMap["mappings"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemDataFields := []valueObject.MarketplaceCatalogItemDataField{}
	if itemMap["dataFields"] != nil {
		itemDataFields, err = repo.catalogItemDataFieldsFactory(itemMap["dataFields"])
		if err != nil {
			return catalogItem, err
		}
	}

	itemInstallTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if itemMap["installTimeoutSecs"] != nil {
		itemInstallTimeoutSecs, err = valueObject.NewUnixTime(
			itemMap["installTimeoutSecs"],
		)
		if err != nil {
			return catalogItem, err
		}
	}

	itemInstallCmdSteps := []valueObject.UnixCommand{}
	if itemMap["installCmdSteps"] != nil {
		itemInstallCmdSteps, err = repo.catalogItemCmdStepsFactory(
			itemMap["installCmdSteps"],
		)
		if err != nil {
			return catalogItem, err
		}
	}

	itemUninstallTimeoutSecs, _ := valueObject.NewUnixTime(600)
	if itemMap["uninstallTimeoutSecs"] != nil {
		itemUninstallTimeoutSecs, err = valueObject.NewUnixTime(
			itemMap["uninstallTimeoutSecs"],
		)
		if err != nil {
			return catalogItem, err
		}
	}

	itemUninstallCmdSteps := []valueObject.UnixCommand{}
	if itemMap["uninstallCmdSteps"] != nil {
		itemUninstallCmdSteps, err = repo.catalogItemCmdStepsFactory(
			itemMap["uninstallCmdSteps"],
		)
		if err != nil {
			return catalogItem, err
		}
	}

	itemUninstallFileNames := []valueObject.UnixFileName{}
	if itemMap["uninstallFileNames"] != nil {
		itemUninstallFileNames, err = repo.catalogItemUninstallFileNamesFactory(
			itemMap["uninstallFileNames"],
		)
		if err != nil {
			return catalogItem, err
		}
	}

	estimatedSizeBytes := valueObject.Byte(1000000000)
	if itemMap["estimatedSizeBytes"] != nil {
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
		itemScreenshotUrls, err = repo.catalogItemScreenshotUrlsFactory(
			itemMap["screenshotUrls"],
		)
		if err != nil {
			return catalogItem, err
		}
	}

	return entity.NewMarketplaceCatalogItem(
		itemManifestVersion, itemId, itemSlugs, itemName, itemType, itemDescription,
		itemServices, itemMappings, itemDataFields, itemInstallTimeoutSecs,
		itemInstallCmdSteps, itemUninstallTimeoutSecs, itemUninstallCmdSteps,
		itemUninstallFileNames, estimatedSizeBytes, itemAvatarUrl, itemScreenshotUrls,
	), nil
}

func (repo *MarketplaceQueryRepo) ReadCatalogItems(
	requestDto dto.ReadMarketplaceCatalogItemsRequest,
) (responseDto dto.ReadMarketplaceCatalogItemsResponse, err error) {
	_, err = os.Stat(infraEnvs.MarketplaceCatalogItemsDir)
	if err != nil {
		marketplaceCmdRepo := NewMarketplaceCmdRepo(repo.persistentDbSvc)
		err = marketplaceCmdRepo.RefreshCatalogItems()
		if err != nil {
			return responseDto, errors.New(
				"RefreshMarketplaceCatalogItemsError: " + err.Error(),
			)
		}
	}

	rawCatalogFilesList, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "find " + infraEnvs.MarketplaceCatalogItemsDir + " -type f " +
			"\\( -name '*.json' -o -name '*.yaml' -o -name '*.yml' \\) " +
			"-not -path '*/.*' -not -name '.*'",
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return responseDto, errors.New("ReadMarketplaceFilesError: " + err.Error())
	}

	if len(rawCatalogFilesList) == 0 {
		return responseDto, errors.New("NoMarketplaceFilesFound")
	}

	rawCatalogFilesListParts := strings.Split(rawCatalogFilesList, "\n")
	if len(rawCatalogFilesListParts) == 0 {
		return responseDto, errors.New("NoMarketplaceFilesFound")
	}

	catalogItems := []entity.MarketplaceCatalogItem{}
	catalogItemsIdsMap := map[uint16]struct{}{}
	for _, rawFilePath := range rawCatalogFilesListParts {
		itemFilePath, err := valueObject.NewUnixFilePath(rawFilePath)
		if err != nil {
			slog.Debug(err.Error(), slog.String("filePath", rawFilePath))
			continue
		}

		catalogItem, err := repo.catalogItemFactory(itemFilePath)
		if err != nil {
			slog.Debug(
				"CatalogMarketplaceItemFactoryError",
				slog.String("filePath", itemFilePath.String()),
				slog.String("err", err.Error()),
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
			slog.Debug(
				"CreateNewCatalogMarketplaceItemIdError",
				slog.String("itemName", catalogItem.Name.String()),
				slog.String("err", err.Error()),
			)
			continue
		}

		catalogItems[itemIndex].Id = nextAvailableId
		itemsIdsSlice = append(itemsIdsSlice, nextAvailableId.Uint16())
	}

	filteredCatalogItems := []entity.MarketplaceCatalogItem{}
	for _, catalogItem := range catalogItems {
		itemId := requestDto.MarketplaceCatalogItemId
		if itemId != nil && catalogItem.Id != *itemId {
			continue
		}

		if requestDto.MarketplaceCatalogItemSlug != nil {
			if !slices.Contains(catalogItem.Slugs, *requestDto.MarketplaceCatalogItemSlug) {
				continue
			}
		}

		if requestDto.MarketplaceCatalogItemName != nil {
			if catalogItem.Name != *requestDto.MarketplaceCatalogItemName {
				continue
			}
		}

		if requestDto.MarketplaceCatalogItemType != nil {
			if catalogItem.Type != *requestDto.MarketplaceCatalogItemType {
				continue
			}
		}

		filteredCatalogItems = append(filteredCatalogItems, catalogItem)
	}

	if len(filteredCatalogItems) > int(requestDto.Pagination.ItemsPerPage) {
		filteredCatalogItems = filteredCatalogItems[:requestDto.Pagination.ItemsPerPage]
	}

	sortDirectionStr := "asc"
	if requestDto.Pagination.SortDirection != nil {
		sortDirectionStr = requestDto.Pagination.SortDirection.String()
	}

	if requestDto.Pagination.SortBy != nil {
		slices.SortStableFunc(filteredCatalogItems, func(a, b entity.MarketplaceCatalogItem) int {
			firstElement := a
			secondElement := b
			if sortDirectionStr != "asc" {
				firstElement = b
				secondElement = a
			}

			switch requestDto.Pagination.SortBy.String() {
			case "id":
				if firstElement.Id.Uint16() < secondElement.Id.Uint16() {
					return -1
				}
				if firstElement.Id.Uint16() > secondElement.Id.Uint16() {
					return 1
				}
				return 0
			case "name":
				return strings.Compare(
					firstElement.Name.String(), secondElement.Name.String(),
				)
			case "type":
				return strings.Compare(
					firstElement.Type.String(), secondElement.Type.String(),
				)
			default:
				return 0
			}
		})
	}

	itemsTotal := uint64(len(filteredCatalogItems))
	pagesTotal := uint32(itemsTotal / uint64(requestDto.Pagination.ItemsPerPage))

	paginationDto := requestDto.Pagination
	paginationDto.ItemsTotal = &itemsTotal
	paginationDto.PagesTotal = &pagesTotal

	return dto.ReadMarketplaceCatalogItemsResponse{
		Pagination:              paginationDto,
		MarketplaceCatalogItems: filteredCatalogItems,
	}, nil
}

func (repo *MarketplaceQueryRepo) ReadFirstCatalogItem(
	requestDto dto.ReadMarketplaceCatalogItemsRequest,
) (catalogItem entity.MarketplaceCatalogItem, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.ReadCatalogItems(requestDto)
	if err != nil {
		return catalogItem, err
	}

	if len(responseDto.MarketplaceCatalogItems) == 0 {
		return catalogItem, errors.New("MarketplaceCatalogItemNotFound")
	}

	return responseDto.MarketplaceCatalogItems[0], nil
}

func (repo *MarketplaceQueryRepo) ReadInstalledItems(
	requestDto dto.ReadMarketplaceInstalledItemsRequest,
) (responseDto dto.ReadMarketplaceInstalledItemsResponse, err error) {
	model := dbModel.MarketplaceInstalledItem{}
	if requestDto.MarketplaceInstalledItemId != nil {
		model.ID = requestDto.MarketplaceInstalledItemId.Uint16()
	}
	if requestDto.MarketplaceInstalledItemHostname != nil {
		model.Hostname = requestDto.MarketplaceInstalledItemHostname.String()
	}
	if requestDto.MarketplaceInstalledItemType != nil {
		model.Type = requestDto.MarketplaceInstalledItemType.String()
	}
	if requestDto.MarketplaceInstalledItemUuid != nil {
		model.InstallUuid = requestDto.MarketplaceInstalledItemUuid.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Model(&model).
		Where(&model).
		Preload("Mappings")

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return responseDto, errors.New(
			"CountMarketplaceInstalledItemsTotalError: " + err.Error(),
		)
	}

	dbQuery.Limit(int(requestDto.Pagination.ItemsPerPage))
	if requestDto.Pagination.LastSeenId == nil {
		offset := int(requestDto.Pagination.PageNumber) * int(requestDto.Pagination.ItemsPerPage)
		dbQuery = dbQuery.Offset(offset)
	} else {
		dbQuery = dbQuery.Where("id > ?", requestDto.Pagination.LastSeenId.String())
	}
	if requestDto.Pagination.SortBy != nil {
		orderStatement := requestDto.Pagination.SortBy.String()
		orderStatement = strcase.ToSnake(orderStatement)
		if orderStatement == "id" {
			orderStatement = "ID"
		}

		if requestDto.Pagination.SortDirection != nil {
			orderStatement += " " + requestDto.Pagination.SortDirection.String()
		}

		dbQuery = dbQuery.Order(orderStatement)
	}

	models := []dbModel.MarketplaceInstalledItem{}
	err = dbQuery.Find(&models).Error
	if err != nil {
		return responseDto, errors.New("ReadMarketplaceInstalledItemsError")
	}

	entities := []entity.MarketplaceInstalledItem{}
	for _, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			slog.Debug(
				"MarketplaceInstalledItemModelToEntityError",
				slog.Uint64("id", uint64(model.ID)), slog.String("err", err.Error()),
			)
			continue
		}

		entities = append(entities, entity)
	}

	itemsTotalUint := uint64(itemsTotal)
	pagesTotal := uint32(
		math.Ceil(float64(itemsTotal) / float64(requestDto.Pagination.ItemsPerPage)),
	)
	responsePagination := dto.Pagination{
		PageNumber:    requestDto.Pagination.PageNumber,
		ItemsPerPage:  requestDto.Pagination.ItemsPerPage,
		SortBy:        requestDto.Pagination.SortBy,
		SortDirection: requestDto.Pagination.SortDirection,
		PagesTotal:    &pagesTotal,
		ItemsTotal:    &itemsTotalUint,
	}

	return dto.ReadMarketplaceInstalledItemsResponse{
		Pagination:                responsePagination,
		MarketplaceInstalledItems: entities,
	}, nil
}

func (repo *MarketplaceQueryRepo) ReadFirstInstalledItem(
	requestDto dto.ReadMarketplaceInstalledItemsRequest,
) (installedItem entity.MarketplaceInstalledItem, err error) {
	requestDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.ReadInstalledItems(requestDto)
	if err != nil {
		return installedItem, err
	}

	if len(responseDto.MarketplaceInstalledItems) == 0 {
		return installedItem, errors.New("MarketplaceInstalledItemNotFound")
	}

	return responseDto.MarketplaceInstalledItems[0], nil
}
