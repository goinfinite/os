package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	marketplaceInfra "github.com/speedianet/os/src/infra/marketplace"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
)

type MarketplaceController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceController {
	return &MarketplaceController{
		persistentDbSvc: persistentDbSvc,
	}
}

// ReadCatalog godoc
// @Summary      ReadCatalog
// @Description  List marketplace catalog services names, types, steps and more.
// @Tags         marketplace
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {array} entity.MarketplaceCatalogItem
// @Router       /v1/marketplace/catalog/ [get]
func (controller *MarketplaceController) ReadCatalog(c echo.Context) error {
	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)

	catalogItems, err := useCase.ReadMarketplaceCatalog(marketplaceQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, catalogItems)
}

func parseDataFieldsFromBody(
	dataFieldsBodyInput interface{},
) []valueObject.MarketplaceInstallableItemDataField {
	dataFields := []valueObject.MarketplaceInstallableItemDataField{}

	if dataFieldsBodyInput == nil {
		return dataFields
	}

	dataFieldsInterfaceSlice, assertOk := dataFieldsBodyInput.([]interface{})
	if !assertOk {
		panic("InvalidDataField")
	}

	for _, dataFieldsInterface := range dataFieldsInterfaceSlice {
		dataFieldMap, assertOk := dataFieldsInterface.(map[string]interface{})
		if !assertOk {
			panic("InvalidDataFieldStructure")
		}

		nameStr, assertOk := dataFieldMap["name"].(string)
		if !assertOk {
			nameStr, assertOk = dataFieldMap["key"].(string)
			if !assertOk {
				panic("InvalidDataField")
			}
		}

		dataField := valueObject.NewMarketplaceInstallableItemDataFieldPanic(
			valueObject.NewDataFieldNamePanic(nameStr),
			valueObject.NewDataFieldValuePanic(dataFieldMap["value"].(string)),
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
// @Param        InstallMarketplaceCatalogItem 	  body    dto.InstallMarketplaceCatalogItem  true  "directory will be the virtual host root directory."
// @Success      201 {object} object{} "MarketplaceCatalogItemInstalled"
// @Router       /v1/marketplace/catalog/ [post]
func (controller *MarketplaceController) InstallCatalogItem(c echo.Context) error {
	requiredParams := []string{"hostname"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))

	var idPtr *valueObject.MarketplaceItemId
	if requestBody["id"] != nil {
		id := valueObject.NewMarketplaceItemIdPanic(requestBody["id"])
		idPtr = &id
	}

	var slugPtr *valueObject.MarketplaceItemSlug
	if requestBody["slug"] != nil {
		slug := valueObject.NewMarketplaceItemSlugPanic(requestBody["slug"])
		slugPtr = &slug
	}

	var urlPathPtr *valueObject.UrlPath
	if requestBody["directory"] != nil {
		urlPath := valueObject.NewUrlPathPanic(requestBody["directory"].(string))
		urlPathPtr = &urlPath
	}
	if requestBody["installDirectory"] != nil {
		urlPath := valueObject.NewUrlPathPanic(requestBody["installDirectory"].(string))
		urlPathPtr = &urlPath
	}

	dataFields := parseDataFieldsFromBody(requestBody["dataFields"])

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)
	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(controller.persistentDbSvc)
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	dto := dto.NewInstallMarketplaceCatalogItem(idPtr, slugPtr, hostname, urlPathPtr, dataFields)
	err := useCase.InstallMarketplaceCatalogItem(
		marketplaceQueryRepo,
		marketplaceCmdRepo,
		vhostQueryRepo,
		vhostCmdRepo,
		dto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "MarketplaceCatalogItemInstalled")
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
	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)

	installedItems, err := useCase.ReadMarketplaceInstalledItems(marketplaceQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, installedItems)
}

// DeleteInstalledItem godoc
// @Summary      DeleteInstalledItem
// @Description  Delete/Uninstall a marketplace installed item.
// @Tags         marketplace
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        installedId path uint true "Marketplace installed item ID that will be deleted."
// @Param        shouldUninstallServices query boolean false "Should uninstall all services that were installed with the marketplace item installation? Default is 'true'."
// @Param        shouldRemoveFiles query boolean false "Should remove all files that were created with the marketplace item installation? Default is 'true'."
// @Success      200 {object} object{} "MarketplaceInstalledItemDeleted"
// @Router       /v1/marketplace/installed/{installedId}/ [delete]
func (controller *MarketplaceController) DeleteInstalledItem(c echo.Context) error {
	installedId := valueObject.NewMarketplaceItemIdPanic(
		c.Param("installedId"),
	)

	var err error

	shouldUninstallServices := true
	if c.QueryParam("shouldUninstallServices") != "" {
		shouldUninstallServices, err = sharedHelper.ParseBoolParam(
			c.QueryParam("shouldUninstallServices"),
		)
		if err != nil {
			shouldUninstallServices = false
		}
	}

	shouldRemoveFiles := true
	if c.QueryParam("shouldRemoveFiles") != "" {
		shouldRemoveFiles, err = sharedHelper.ParseBoolParam(
			c.QueryParam("shouldRemoveFiles"),
		)
		if err != nil {
			shouldRemoveFiles = false
		}
	}

	deleteMarketplaceInstalledItem := dto.NewDeleteMarketplaceInstalledItem(
		installedId, shouldUninstallServices, shouldRemoveFiles,
	)

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)
	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(controller.persistentDbSvc)

	err = useCase.DeleteMarketplaceInstalledItem(
		marketplaceQueryRepo,
		marketplaceCmdRepo,
		deleteMarketplaceInstalledItem,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "MarketplaceInstalledItemDeleted")
}
