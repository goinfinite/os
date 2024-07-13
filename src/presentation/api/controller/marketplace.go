package apiController

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	"github.com/speedianet/os/src/presentation/service"
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
// @Success      200 {array} entity.MarketplaceCatalogItem
// @Router       /v1/marketplace/catalog/ [get]
func (controller *MarketplaceController) ReadCatalog(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(
		c, controller.marketplaceService.ReadCatalog(),
	)
}

func (controller *MarketplaceController) parseDataFields(
	rawDataFields []interface{},
) ([]valueObject.MarketplaceInstallableItemDataField, error) {
	dataFields := []valueObject.MarketplaceInstallableItemDataField{}

	for fieldIndex, rawDataField := range rawDataFields {
		errPrefix := "[index " + strconv.Itoa(fieldIndex) + "] "

		rawDataFieldMap, assertOk := rawDataField.(map[string]interface{})
		if !assertOk {
			return dataFields, errors.New(errPrefix + "InvalidDataFieldStructure")
		}

		rawName, exists := rawDataFieldMap["name"]
		if !exists {
			rawName, exists = rawDataFieldMap["key"]
			if !exists {
				return dataFields, errors.New(errPrefix + "DataFieldNameNotFound")
			}
			rawDataFieldMap["name"] = rawName
		}

		fieldName, err := valueObject.NewDataFieldName(rawName)
		if err != nil {
			return dataFields, errors.New(errPrefix + err.Error())
		}

		fieldValue, err := valueObject.NewDataFieldValue(rawDataFieldMap["value"])
		if err != nil {
			return dataFields, errors.New(errPrefix + err.Error())
		}

		dataField := valueObject.NewMarketplaceInstallableItemDataFieldPanic(
			fieldName, fieldValue,
		)

		dataFields = append(dataFields, dataField)
	}

	return dataFields, nil
}

// InstallCatalogItem	 godoc
// @Summary      InstallCatalogItem
// @Description  Install a marketplace catalog item.
// @Tags         marketplace
// @Accept       json
// @Produce      json
// @Param        InstallMarketplaceCatalogItem 	  body    dto.InstallMarketplaceCatalogItem  true  "urlPath is both the install directory and HTTP sub-directory."
// @Success      201 {object} object{} "MarketplaceCatalogItemInstalled"
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
		_, isMapStringInterface := requestBody["dataFields"].(map[string]interface{})
		if isMapStringInterface {
			requestBody["dataFields"] = []interface{}{requestBody["dataFields"]}
		}

		dataFieldsSlice, assertOk := requestBody["dataFields"].([]interface{})
		if !assertOk {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "DataFieldsMustBeArray",
			)
		}

		dataFields, err := controller.parseDataFields(dataFieldsSlice)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		requestBody["dataFields"] = dataFields
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
// @Success      200 {array} entity.MarketplaceInstalledItem
// @Router       /v1/marketplace/installed/ [get]
func (controller *MarketplaceController) ReadInstalledItems(c echo.Context) error {
	return apiHelper.ServiceResponseWrapper(
		c, controller.marketplaceService.ReadInstalledItems(),
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
		c, controller.marketplaceService.DeleteInstalledItem(requestBody),
	)
}
