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
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
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
			slog.Error(
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
			slog.Error(err.Error(), slog.Any("cmdStep", rawItemCmdStep))
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
			slog.Error(err.Error(), slog.Any("fileName", rawItemUninstallFileName))
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

	itemInstallCmdSteps := []valueObject.UnixCommand{}
	if itemMap["installCmdSteps"] != nil {
		itemInstallCmdSteps, err = repo.catalogItemCmdStepsFactory(
			itemMap["installCmdSteps"],
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
		itemScreenshotUrls, err = repo.catalogItemScreenshotUrlsFactory(
			itemMap["screenshotUrls"],
		)
		if err != nil {
			return catalogItem, err
		}
	}

	return entity.NewMarketplaceCatalogItem(
		itemId, itemSlugs, itemName, itemType, itemDescription, itemServices,
		itemMappings, itemDataFields, itemInstallCmdSteps, itemUninstallCmdSteps,
		itemUninstallFileNames, estimatedSizeBytes, itemAvatarUrl, itemScreenshotUrls,
	), nil
}

func (repo *MarketplaceQueryRepo) ReadCatalogItems(
	readDto dto.ReadMarketplaceCatalogItemsRequest,
) (catalogItemsDto dto.ReadMarketplaceCatalogItemsResponse, err error) {
	_, err = os.Stat(infraEnvs.MarketplaceItemsDir)
	if err != nil {
		marketplaceCmdRepo := NewMarketplaceCmdRepo(repo.persistentDbSvc)
		err = marketplaceCmdRepo.RefreshItems()
		if err != nil {
			return catalogItemsDto, errors.New(
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
		return catalogItemsDto, errors.New("ReadMarketplaceFilesError: " + err.Error())
	}

	if len(rawCatalogFilesList) == 0 {
		return catalogItemsDto, errors.New("NoMarketplaceFilesFound")
	}

	rawCatalogFilesListParts := strings.Split(rawCatalogFilesList, "\n")
	if len(rawCatalogFilesListParts) == 0 {
		return catalogItemsDto, errors.New("NoMarketplaceFilesFound")
	}

	catalogItems := []entity.MarketplaceCatalogItem{}
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

	filteredCatalogItems := []entity.MarketplaceCatalogItem{}
	for _, catalogItem := range catalogItems {
		if len(catalogItems) >= int(readDto.Pagination.ItemsPerPage) {
			break
		}

		if readDto.Id != nil && catalogItem.Id != *readDto.Id {
			continue
		}

		if readDto.Slug != nil {
			if !slices.Contains(catalogItem.Slugs, *readDto.Slug) {
				continue
			}
		}

		if readDto.Name != nil {
			if !strings.EqualFold(catalogItem.Name.String(), readDto.Name.String()) {
				continue
			}
		}

		if readDto.Type != nil && catalogItem.Type != *readDto.Type {
			if !strings.EqualFold(catalogItem.Type.String(), readDto.Type.String()) {
				continue
			}
		}

		filteredCatalogItems = append(filteredCatalogItems, catalogItem)
	}

	sortDirectionStr := "asc"
	if readDto.Pagination.SortDirection != nil {
		sortDirectionStr = readDto.Pagination.SortDirection.String()
	}

	if readDto.Pagination.SortBy != nil {
		slices.SortStableFunc(filteredCatalogItems, func(a, b entity.MarketplaceCatalogItem) int {
			firstElement := a
			secondElement := b
			if sortDirectionStr != "asc" {
				firstElement = b
				secondElement = a
			}

			switch readDto.Pagination.SortBy.String() {
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
	pagesTotal := uint32(itemsTotal / uint64(readDto.Pagination.ItemsPerPage))

	paginationDto := readDto.Pagination
	paginationDto.ItemsTotal = &itemsTotal
	paginationDto.PagesTotal = &pagesTotal

	return dto.ReadMarketplaceCatalogItemsResponse{
		Pagination: paginationDto,
		Items:      filteredCatalogItems,
	}, nil
}

func (repo *MarketplaceQueryRepo) ReadUniqueCatalogItem(
	readDto dto.ReadMarketplaceCatalogItemsRequest,
) (catalogItem entity.MarketplaceCatalogItem, err error) {
	readDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.ReadCatalogItems(readDto)
	if err != nil {
		return catalogItem, err
	}

	if len(responseDto.Items) == 0 {
		return catalogItem, errors.New("MarketplaceCatalogItemNotFound")
	}

	foundCatalogItem := responseDto.Items[0]
	return foundCatalogItem, nil
}

func (repo *MarketplaceQueryRepo) ReadInstalledItems(
	readDto dto.ReadMarketplaceInstalledItemsRequest,
) (installedItemsDto dto.ReadMarketplaceInstalledItemsResponse, err error) {
	model := dbModel.MarketplaceInstalledItem{}
	if readDto.Id != nil {
		model.ID = uint(readDto.Id.Uint16())
	}
	if readDto.Hostname != nil {
		model.Hostname = readDto.Hostname.String()
	}
	if readDto.Type != nil {
		model.Type = readDto.Type.String()
	}
	if readDto.InstallationUuid != nil {
		model.InstallUuid = readDto.InstallationUuid.String()
	}

	dbQuery := repo.persistentDbSvc.Handler.
		Where(&model).
		Limit(int(readDto.Pagination.ItemsPerPage))
	if readDto.Pagination.LastSeenId == nil {
		offset := int(readDto.Pagination.PageNumber) * int(readDto.Pagination.ItemsPerPage)
		dbQuery = dbQuery.Offset(offset)
	} else {
		dbQuery = dbQuery.Where("id > ?", readDto.Pagination.LastSeenId.String())
	}
	if readDto.Pagination.SortBy != nil {
		orderStatement := readDto.Pagination.SortBy.String()
		orderStatement = strcase.ToSnake(orderStatement)
		if orderStatement == "id" {
			orderStatement = "ID"
		}

		if readDto.Pagination.SortDirection != nil {
			orderStatement += " " + readDto.Pagination.SortDirection.String()
		}

		dbQuery = dbQuery.Order(orderStatement)
	}

	models := []dbModel.MarketplaceInstalledItem{}
	err = dbQuery.Preload("Mappings").Find(&models).Error
	if err != nil {
		return installedItemsDto, errors.New("ReadMarketplaceInstalledItemsError")
	}

	var itemsTotal int64
	err = dbQuery.Count(&itemsTotal).Error
	if err != nil {
		return installedItemsDto, errors.New(
			"CountMarketplaceInstalledItemsTotalError: " + err.Error(),
		)
	}

	entities := []entity.MarketplaceInstalledItem{}
	for _, model := range models {
		entity, err := model.ToEntity()
		if err != nil {
			slog.Error(
				"MarketplaceInstalledItemModelToEntityError",
				slog.Uint64("id", uint64(model.ID)), slog.Any("error", err),
			)
			continue
		}

		entities = append(entities, entity)
	}

	itemsTotalUint := uint64(itemsTotal)
	pagesTotal := uint32(
		math.Ceil(float64(itemsTotal) / float64(readDto.Pagination.ItemsPerPage)),
	)
	responsePagination := dto.Pagination{
		PageNumber:    readDto.Pagination.PageNumber,
		ItemsPerPage:  readDto.Pagination.ItemsPerPage,
		SortBy:        readDto.Pagination.SortBy,
		SortDirection: readDto.Pagination.SortDirection,
		PagesTotal:    &pagesTotal,
		ItemsTotal:    &itemsTotalUint,
	}

	return dto.ReadMarketplaceInstalledItemsResponse{
		Pagination: responsePagination,
		Items:      entities,
	}, nil
}

func (repo *MarketplaceQueryRepo) ReadUniqueInstalledItem(
	readDto dto.ReadMarketplaceInstalledItemsRequest,
) (installedItem entity.MarketplaceInstalledItem, err error) {
	readDto.Pagination = dto.Pagination{
		PageNumber:   0,
		ItemsPerPage: 1,
	}
	responseDto, err := repo.ReadInstalledItems(readDto)
	if err != nil {
		return installedItem, err
	}

	if len(responseDto.Items) == 0 {
		return installedItem, errors.New("MarketplaceInstalledItemNotFound")
	}

	foundInstalledItem := responseDto.Items[0]
	return foundInstalledItem, nil
}
