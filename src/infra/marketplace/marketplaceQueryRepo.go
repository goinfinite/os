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
	if catalogItemFileExt == "yaml" {
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

func (repo *MarketplaceQueryRepo) catalogItemMappingFactory(
	catalogItemMappingMap map[string]interface{},
) (valueObject.MarketplaceItemMapping, error) {
	var mapping valueObject.MarketplaceItemMapping

	rawPath, assertOk := catalogItemMappingMap["path"].(string)
	if !assertOk {
		return mapping, errors.New("InvalidMappingPath")
	}
	path, err := valueObject.NewMappingPath(rawPath)
	if err != nil {
		return mapping, err
	}

	rawMatchPattern, assertOk := catalogItemMappingMap["matchPattern"].(string)
	if !assertOk {
		return mapping, errors.New("InvalidMappingMatchPattern")
	}
	matchPattern, err := valueObject.NewMappingMatchPattern(rawMatchPattern)
	if err != nil {
		return mapping, err
	}

	rawTargetType, assertOk := catalogItemMappingMap["targetType"].(string)
	if !assertOk {
		return mapping, errors.New("InvalidMappingTargetType")
	}
	targetType, err := valueObject.NewMappingTargetType(rawTargetType)
	if err != nil {
		return mapping, err
	}

	var targetSvcNamePtr *valueObject.ServiceName
	_, rawTargetSvcNameExists := catalogItemMappingMap["targetServiceName"]
	if rawTargetSvcNameExists {
		rawTargetSvcName, assertOk := catalogItemMappingMap["targetServiceName"].(string)
		if !assertOk {
			return mapping, errors.New("InvalidMappingTargetSvcName")
		}
		targetSvcName, err := valueObject.NewServiceName(rawTargetSvcName)
		if err != nil {
			return mapping, err
		}
		targetSvcNamePtr = &targetSvcName
	}

	var targetUrlPtr *valueObject.Url
	_, rawTargetUrlExists := catalogItemMappingMap["targetUrl"]
	if rawTargetUrlExists {
		rawTargetUrl, assertOk := catalogItemMappingMap["targetUrl"].(string)
		if !assertOk {
			return mapping, errors.New("InvalidMappingTargetUrl")
		}
		targetUrl, err := valueObject.NewUrl(rawTargetUrl)
		if err != nil {
			return mapping, err
		}
		targetUrlPtr = &targetUrl
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	_, rawTargetHttpResponseCodeExists := catalogItemMappingMap["targetHttpResponseCode"]
	if rawTargetHttpResponseCodeExists {
		rawTargetHttpResponseCode, assertOk := catalogItemMappingMap["targetHttpResponseCode"].(string)
		if !assertOk {
			return mapping, errors.New("InvalidMappingTargetHttpResponseCode")
		}
		targetHttpResponseCode, err := valueObject.NewHttpResponseCode(rawTargetHttpResponseCode)
		if err != nil {
			return mapping, err
		}
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	var targetInlineHtmlContentPtr *valueObject.InlineHtmlContent
	_, rawTargetInlineHtmlContentExists := catalogItemMappingMap["targetInlineHtmlContent"]
	if rawTargetInlineHtmlContentExists {
		rawTargetInlineHtmlContent, assertOk := catalogItemMappingMap["targetInlineHtmlContent"].(string)
		if !assertOk {
			return mapping, errors.New("InvalidMappingTargetInlineHtmlContent")
		}
		targetInlineHtmlContent, err := valueObject.NewInlineHtmlContent(rawTargetInlineHtmlContent)
		if err != nil {
			return mapping, err
		}
		targetInlineHtmlContentPtr = &targetInlineHtmlContent
	}

	return valueObject.NewMarketplaceItemMapping(
		path,
		matchPattern,
		targetType,
		targetSvcNamePtr,
		targetUrlPtr,
		targetHttpResponseCodePtr,
		targetInlineHtmlContentPtr,
	), nil
}

func (repo *MarketplaceQueryRepo) catalogItemDataFieldFactory(
	catalogItemDataFieldMap map[string]interface{},
) (valueObject.MarketplaceItemDataField, error) {
	var dataField valueObject.MarketplaceItemDataField

	rawKey, assertOk := catalogItemDataFieldMap["key"].(string)
	if !assertOk {
		return dataField, errors.New("InvalidDataFieldKey")
	}
	key, err := valueObject.NewDataFieldKey(rawKey)
	if err != nil {
		return dataField, err
	}

	rawValue, assertOk := catalogItemDataFieldMap["value"].(string)
	if !assertOk {
		return dataField, errors.New("InvalidDataFieldValue")
	}
	value, err := valueObject.NewDataFieldValue(rawValue)
	if err != nil {
		return dataField, err
	}

	rawIsRequired, assertOk := catalogItemDataFieldMap["isRequired"].(bool)
	if !assertOk {
		return dataField, errors.New("InvalidDataFieldIsRequired")
	}

	_, defaultValueExists := catalogItemDataFieldMap["defaultValue"]
	var defaultValuePtr *valueObject.DataFieldValue
	if defaultValueExists {
		rawDefaultValue, assertOk := catalogItemDataFieldMap["defaultValue"].(string)
		if !assertOk {
			return dataField, errors.New("InvalidDataFieldDefaultValue")
		}
		defaultValue, err := valueObject.NewDataFieldValue(rawDefaultValue)
		if err != nil {
			return dataField, err
		}
		defaultValuePtr = &defaultValue
	}

	return valueObject.NewMarketplaceItemDataField(
		key,
		value,
		rawIsRequired,
		defaultValuePtr,
	)
}

func (repo *MarketplaceQueryRepo) catalogItemFactory(
	catalogItemFilePath valueObject.UnixFilePath,
) (entity.MarketplaceCatalogItem, error) {
	var catalogItem entity.MarketplaceCatalogItem

	catalogItemMap, err := repo.getCatalogItemMapFromFilePath(catalogItemFilePath)
	if err != nil {
		return catalogItem, err
	}

	catalogItemId, _ := valueObject.NewMarketplaceItemId(1)

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

	rawCatalogItemSvcNames, assertOk := catalogItemMap["serviceNames"].([]interface{})
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemServiceNames")
	}
	catalogItemSvcNames := []valueObject.ServiceName{}
	for _, rawCatalogItemSvcName := range rawCatalogItemSvcNames {
		catalogItemSvcName, err := valueObject.NewServiceName(rawCatalogItemSvcName.(string))
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawCatalogItemSvcName)
			continue
		}

		catalogItemSvcNames = append(catalogItemSvcNames, catalogItemSvcName)
	}

	rawCatalogItemMappings, assertOk := catalogItemMap["mappings"].([]interface{})
	if !assertOk {
		log.Printf("InvalidMarketplaceCatalogItemMappings")
	}
	catalogItemMappings := []valueObject.MarketplaceItemMapping{}
	for _, rawCatalogItemMapping := range rawCatalogItemMappings {
		catalogItemMapping, err := repo.catalogItemMappingFactory(rawCatalogItemMapping.(map[string]interface{}))
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawCatalogItemMapping)
			continue
		}

		catalogItemMappings = append(catalogItemMappings, catalogItemMapping)
	}

	rawCatalogItemDataFields, assertOk := catalogItemMap["dataFields"].([]interface{})
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemDataFields")
	}
	catalogItemDataFields := []valueObject.MarketplaceItemDataField{}
	for _, rawCatalogItemDataField := range rawCatalogItemDataFields {
		catalogItemDataField, err := repo.catalogItemDataFieldFactory(rawCatalogItemDataField.(map[string]interface{}))
		if err != nil {
			log.Printf("%s: %v", err.Error(), rawCatalogItemDataField)
			return catalogItem, err
		}

		catalogItemDataFields = append(catalogItemDataFields, catalogItemDataField)
	}

	rawCatalogItemCmdSteps, assertOk := catalogItemMap["cmdSteps"].([]interface{})
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemCmdSteps")
	}
	catalogItemCmdSteps := []valueObject.MarketplaceItemInstallStep{}
	for _, rawCatalogItemCmdStep := range rawCatalogItemCmdSteps {
		catalogItemCmdStep, err := valueObject.NewMarketplaceItemInstallStep(rawCatalogItemCmdStep.(string))
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawCatalogItemCmdStep)
			return catalogItem, err
		}

		catalogItemCmdSteps = append(catalogItemCmdSteps, catalogItemCmdStep)
	}

	rawCatalogEstimatedSizeBytes, assertOk := catalogItemMap["estimatedSizeBytes"].(float64)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogEstimatedSizeBytes")
	}
	catalogEstimatedSizeBytes := valueObject.Byte(rawCatalogEstimatedSizeBytes)

	rawCatalogItemAvatarUrl, assertOk := catalogItemMap["avatarUrl"].(string)
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemAvatarUrl")
	}
	catalogItemAvatarUrl, err := valueObject.NewUrl(rawCatalogItemAvatarUrl)
	if err != nil {
		return catalogItem, err
	}

	rawCatalogItemScreenshotUrls, assertOk := catalogItemMap["screenshotUrls"].([]interface{})
	if !assertOk {
		return catalogItem, errors.New("InvalidMarketplaceCatalogItemScreenshotUrls")
	}
	catalogItemScreenshotUrls := []valueObject.Url{}
	for _, rawCatalogItemScreenshotUrl := range rawCatalogItemScreenshotUrls {
		catalogItemScreenshotUrl, err := valueObject.NewUrl(rawCatalogItemScreenshotUrl.(string))
		if err != nil {
			log.Printf("%s: %s", err.Error(), rawCatalogItemScreenshotUrl)
			continue
		}

		catalogItemScreenshotUrls = append(catalogItemScreenshotUrls, catalogItemScreenshotUrl)
	}

	return entity.NewMarketplaceCatalogItem(
		catalogItemId,
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
		catalogItemId, _ := valueObject.NewMarketplaceItemId(catalogItemIdInt)
		catalogItem.Id = catalogItemId

		catalogItems = append(catalogItems, catalogItem)
	}

	return catalogItems, nil
}

func (repo *MarketplaceQueryRepo) GetCatalogItemById(
	id valueObject.MarketplaceItemId,
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

	return catalogItem, nil
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
