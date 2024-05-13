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
// @Router       /marketplace/catalog/ [get]
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
			panic("InvalidDataField")
		}

		dataField := valueObject.NewMarketplaceInstallableItemDataFieldPanic(
			valueObject.NewDataFieldNamePanic(dataFieldMap["name"].(string)),
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
// @Param        InstallMarketplaceCatalogItem 	  body    dto.InstallMarketplaceCatalogItem  true  "InstallMarketplaceCatalogItem (directory is optional)"
// @Success      201 {object} object{} "MarketplaceCatalogItemInstalled"
// @Router       /marketplace/catalog/ [post]
func (controller *MarketplaceController) InstallCatalogItem(c echo.Context) error {
	requiredParams := []string{"id", "hostname"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	catalogId := valueObject.NewMarketplaceCatalogItemIdPanic(requestBody["id"])
	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))

	var urlPathPtr *valueObject.UrlPath
	if requestBody["directory"] != nil {
		urlPath := valueObject.NewUrlPathPanic(
			requestBody["directory"].(string),
		)
		urlPathPtr = &urlPath
	}

	dataFields := parseDataFieldsFromBody(requestBody["dataFields"])

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)
	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(controller.persistentDbSvc)
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	dto := dto.NewInstallMarketplaceCatalogItem(catalogId, hostname, urlPathPtr, dataFields)
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
// @Router       /marketplace/installed/ [get]
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
// @Param        installedId path uint true "MarketplaceInstalledItemId"
// @Param        shouldUninstallServices query boolean false "ShouldUninstallServices"
// @Param        shouldRemoveFiles query boolean false "ShouldRemoveFiles"
// @Success      200 {object} object{} "MarketplaceInstalledItemDeleted"
// @Router       /marketplace/installed/{installedId}/ [delete]
func (controller *MarketplaceController) DeleteInstalledItem(c echo.Context) error {
	installedId := valueObject.NewMarketplaceInstalledItemIdPanic(
		c.Param("installedId"),
	)

	var err error

	shouldUninstallServices := true
	if c.QueryParam("shouldUninstallServices") != "" {
		shouldUninstallServices, err = apiHelper.ParseBoolParam(
			c.QueryParam("shouldUninstallServices"),
		)
		if err != nil {
			shouldUninstallServices = false
		}
	}

	shouldRemoveFiles := true
	if c.QueryParam("shouldRemoveFiles") != "" {
		shouldRemoveFiles, err = apiHelper.ParseBoolParam(
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
