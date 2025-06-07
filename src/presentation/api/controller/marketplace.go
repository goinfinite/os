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
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/labstack/echo/v4"
)

type MarketplaceController struct {
	marketplaceLiaison *liaison.MarketplaceLiaison
	persistentDbSvc    *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MarketplaceController {
	return &MarketplaceController{
		marketplaceLiaison: liaison.NewMarketplaceLiaison(persistentDbSvc, trailDbSvc),
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
// @Param        id query  uint  false  "Id"
// @Param        slug query  string  false  "Slug"
// @Param        name query  string  false  "Name"
// @Param        type query  string  false  "Type"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadMarketplaceCatalogItemsResponse
// @Router       /v1/marketplace/catalog/ [get]
func (controller *MarketplaceController) ReadCatalog(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.marketplaceLiaison.ReadCatalog(requestInputData),
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
			slog.Debug(
				"InvalidDataFieldStringStructure",
				slog.String("rawDataField", rawDataField),
			)
			continue
		}

		dataFieldMap := map[string]interface{}{
			"name":  rawDataFieldParts[0],
			"value": rawDataFieldParts[1],
		}
		dataFieldsMapSlice = append(dataFieldsMapSlice, dataFieldMap)
	}

	return dataFieldsMapSlice
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

	for index, rawDataField := range rawDataFieldsSlice {
		rawDataFieldMap, assertOk := rawDataField.(map[string]interface{})
		if !assertOk {
			slog.Debug("InvalidDataFieldStructure", slog.Any("fieldIndex", index))
			continue
		}

		fieldName, err := valueObject.NewDataFieldName(rawDataFieldMap["name"])
		if err != nil {
			slog.Debug(err.Error(), slog.Any("fieldIndex", index))
			continue
		}

		fieldValue, err := valueObject.NewDataFieldValue(rawDataFieldMap["value"])
		if err != nil {
			slog.Debug(err.Error(), slog.Any("fieldName", fieldName.String()))
			continue
		}

		dataField := valueObject.NewMarketplaceInstallableItemDataField(
			fieldName, fieldValue,
		)
		dataFields = append(dataFields, dataField)
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	possibleUrlPathKeys := []string{"urlPath", "directory", "installDirectory"}
	for _, key := range possibleUrlPathKeys {
		if requestInputData[key] == nil {
			continue
		}

		requestInputData["urlPath"] = requestInputData[key]
		break
	}

	if requestInputData["dataFields"] != nil {
		requestInputData["dataFields"] = controller.parseDataFields(
			requestInputData["dataFields"],
		)
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.marketplaceLiaison.InstallCatalogItem(requestInputData, true),
	)
}

// ReadInstalledItems godoc
// @Summary      ReadInstalledItems
// @Description  List marketplace installed items.
// @Tags         marketplace
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        id query  uint  false  "Id"
// @Param        hostname query  string  false  "Hostname"
// @Param        type query  string  false  "Type"
// @Param        installationUuid query  string  false  "InstallUuid"
// @Param        pageNumber query  uint  false  "PageNumber (Pagination)"
// @Param        itemsPerPage query  uint  false  "ItemsPerPage (Pagination)"
// @Param        sortBy query  string  false  "SortBy (Pagination)"
// @Param        sortDirection query  string  false  "SortDirection (Pagination)"
// @Param        lastSeenId query  string  false  "LastSeenId (Pagination)"
// @Success      200 {object} dto.ReadMarketplaceInstalledItemsResponse
// @Router       /v1/marketplace/installed/ [get]
func (controller *MarketplaceController) ReadInstalledItems(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.marketplaceLiaison.ReadInstalledItems(requestInputData),
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	return apiHelper.LiaisonResponseWrapper(
		c, controller.marketplaceLiaison.DeleteInstalledItem(requestInputData, true),
	)
}

func (controller *MarketplaceController) AutoRefreshMarketplaceCatalogItems() {
	refreshIntervalHours := 24 / useCase.RefreshMarketplaceCatalogItemsAmountPerDay

	taskInterval := time.Duration(refreshIntervalHours) * time.Hour
	timer := time.NewTicker(taskInterval)
	defer timer.Stop()

	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(
		controller.persistentDbSvc,
	)
	for range timer.C {
		useCase.RefreshMarketplaceCatalogItems(marketplaceCmdRepo)
	}
}
