package apiController

import (
	"log/slog"
	"strings"
	"time"

	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	marketplaceInfra "github.com/goinfinite/os/src/infra/marketplace"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/labstack/echo/v4"
)

type MarketplaceController struct {
	marketplaceService *service.MarketplaceService
	persistentDbSvc    *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceController {
	return &MarketplaceController{
		marketplaceService: service.NewMarketplaceService(persistentDbSvc),
		persistentDbSvc:    persistentDbSvc,
	}
}

// ReadCatalog godoc
// @Summary      ReadCatalog
// @Description  List marketplace catalog items.
// @Tags         marketplace
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        itemId query  uint  false  "Id"
// @Param        itemSlug query  string  false  "Slug"
// @Param        itemName query  string  false  "Name"
// @Param        itemType query  string  false  "Type"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadMarketplaceCatalogItemsResponse
// @Router       /v1/marketplace/catalog/ [get]
func (controller *MarketplaceController) ReadCatalog(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.marketplaceService.ReadCatalog(requestBody),
	)
}

func (controller *MarketplaceController) transformDataFieldsIntoMap(
	rawDataFields string,
) []map[string]interface{} {
	dataFieldsMapSlice := []map[string]interface{}{}
	if len(rawDataFields) == 0 {
		return dataFieldsMapSlice
	}

	rawDataFieldsSlice := strings.Split(rawDataFields, ";")
	for _, rawDataField := range rawDataFieldsSlice {
		rawDataFieldParts := strings.Split(rawDataField, ":")
		if len(rawDataFieldParts) != 2 {
			slog.Error(
				"InvalidDataFieldStringStructure",
				slog.String("rawDataField", rawDataField),
			)
			continue
		}

		dataFieldsMapSlice = append(
			dataFieldsMapSlice,
			map[string]interface{}{rawDataFieldParts[0]: rawDataFieldParts[1]},
		)
	}

	return dataFieldsMapSlice
}

func (controller *MarketplaceController) parseDataFieldMap(
	rawDataFields map[string]interface{},
) []valueObject.MarketplaceInstallableItemDataField {
	dataFields := []valueObject.MarketplaceInstallableItemDataField{}

	fieldIndex := 0
	for rawFieldName, rawFieldValue := range rawDataFields {
		fieldIndex++

		fieldName, err := valueObject.NewDataFieldName(rawFieldName)
		if err != nil {
			slog.Error(err.Error(), slog.Int("fieldIndex", fieldIndex))
			continue
		}

		fieldValue, err := valueObject.NewDataFieldValue(rawFieldValue)
		if err != nil {
			slog.Error(err.Error(), slog.String("fieldName", fieldName.String()))
			continue
		}

		dataField, err := valueObject.NewMarketplaceInstallableItemDataField(
			fieldName, fieldValue,
		)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.String("fieldName", fieldName.String()),
				slog.String("fieldValue", fieldValue.String()),
			)
			continue
		}

		dataFields = append(dataFields, dataField)
	}

	return dataFields
}

// DataFields has multiple possible structures which this parser can handle:
// "dataFieldName:dataFieldValue;dataFieldName:dataFieldValue" (string slice, semicolon separated items)
// { "dataFieldName": "dataFieldValue" } (map[string]interface{})
// [{ "dataFieldName": "dataFieldValue" }] (map[string]interface{} slice)
func (controller *MarketplaceController) parseDataFields(
	dataFieldsAsUnknownType any,
) []valueObject.MarketplaceInstallableItemDataField {
	dataFields := []valueObject.MarketplaceInstallableItemDataField{}

	rawDataFieldsSlice := []interface{}{}
	switch dataFieldsValues := dataFieldsAsUnknownType.(type) {
	case map[string]interface{}:
		rawDataFieldsSlice = []interface{}{dataFieldsValues}
	case string:
		dataFieldsMaps := controller.transformDataFieldsIntoMap(dataFieldsValues)
		for _, dataFieldMap := range dataFieldsMaps {
			rawDataFieldsSlice = append(rawDataFieldsSlice, dataFieldMap)
		}
	case []interface{}:
		rawDataFieldsSlice = dataFieldsValues
	}

	for _, rawDataField := range rawDataFieldsSlice {
		rawDataFieldMap, assertOk := rawDataField.(map[string]interface{})
		if !assertOk {
			slog.Error(
				"InvalidDataFieldStructure", slog.Any("rawDataField", rawDataField),
			)
			continue
		}
		dataFields = append(dataFields, controller.parseDataFieldMap(rawDataFieldMap)...)
	}

	return dataFields
}

// InstallCatalogItem	 godoc
// @Summary      InstallCatalogItem
// @Description  Install a marketplace catalog item.
// @Tags         marketplace
// @Accept       json
// @Produce      json
// @Param        InstallMarketplaceCatalogItem 	  body    dto.InstallMarketplaceCatalogItem  true  "urlPath is both the install directory and HTTP sub-directory."
// @Success      201 {object} object{} "MarketplaceCatalogItemInstallationScheduled"
// @Router       /v1/marketplace/catalog/ [post]
func (controller *MarketplaceController) InstallCatalogItem(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	possibleUrlPathKeys := []string{"urlPath", "directory", "installDirectory"}
	for _, key := range possibleUrlPathKeys {
		if requestBody[key] == nil {
			continue
		}

		requestBody["urlPath"] = requestBody[key]
		break
	}

	if requestBody["dataFields"] != nil {
		requestBody["dataFields"] = controller.parseDataFields(requestBody["dataFields"])
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.marketplaceService.InstallCatalogItem(requestBody, true),
	)
}

// ReadInstalledItems godoc
// @Summary      ReadInstalledItems
// @Description  List marketplace installed items.
// @Tags         marketplace
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        itemId query  uint  false  "Id"
// @Param        itemSlug query  string  false  "Slug"
// @Param        itemName query  string  false  "Name"
// @Param        itemType query  string  false  "Type"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadMarketplaceCatalogItemsResponse
// @Router       /v1/marketplace/installed/ [get]
func (controller *MarketplaceController) ReadInstalledItems(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.marketplaceService.ReadInstalledItems(requestBody),
	)
}

// DeleteInstalledItem godoc
// @Summary      DeleteInstalledItem
// @Description  Delete/Uninstall an installed item.
// @Tags         marketplace
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        installedId path uint true "MarketplaceInstalledItemId to delete."
// @Param        shouldUninstallServices query boolean false "Should uninstall all services not being used? Default is 'true'."
// @Success      200 {object} object{} "MarketplaceInstalledItemDeleted"
// @Router       /v1/marketplace/installed/{installedId}/ [delete]
func (controller *MarketplaceController) DeleteInstalledItem(c echo.Context) error {
	requestBody := map[string]interface{}{
		"installedId": c.Param("installedId"),
	}

	if c.QueryParam("shouldUninstallServices") != "" {
		requestBody["shouldUninstallServices"] = c.QueryParam("shouldUninstallServices")
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.marketplaceService.DeleteInstalledItem(requestBody, true),
	)
}

func (controller *MarketplaceController) AutoRefreshMarketplaceItems() {
	taskInterval := time.Duration(2) * time.Minute
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(
		controller.persistentDbSvc,
	)
	for range timer.C {
		useCase.RefreshMarketplaceItems(marketplaceCmdRepo)
	}
}
