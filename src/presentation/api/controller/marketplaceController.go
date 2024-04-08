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

// GetMarketplaceCatalog godoc
// @Summary      GetMarketplaceCatalog
// @Description  List marketplace catalog services names, types, steps and more.
// @Tags         marketplace
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {string} entity.MarketplaceCatalogItem
// @Router       /marketplace/catalog/ [get]
func (controller *MarketplaceController) GetCatalogController(c echo.Context) error {
	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)
	marketplaceItems, err := useCase.GetMarketplaceCatalog(marketplaceQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, marketplaceItems)
}

func getDataFieldsFromBody(
	dataFieldsBodyInput interface{},
) []valueObject.MarketplaceItemDataField {
	dataFields := []valueObject.MarketplaceItemDataField{}

	dataFieldsInterfaceSlice, assertOk := dataFieldsBodyInput.([]interface{})
	if !assertOk {
		panic("InvalidDataField")
	}

	for _, dataFieldsInterface := range dataFieldsInterfaceSlice {
		dataFieldMap, assertOk := dataFieldsInterface.(map[string]interface{})
		if !assertOk {
			panic("InvalidDataField")
		}

		dataField := valueObject.NewMarketplaceItemDataField(
			valueObject.NewDataFieldKeyPanic(dataFieldMap["key"].(string)),
			valueObject.NewDataFieldValuePanic(dataFieldMap["value"].(string)),
			false,
			nil,
		)

		dataFields = append(dataFields, dataField)
	}

	return dataFields
}

// InstallMarketplaceCatalogItem	 godoc
// @Summary      InstallMarketplaceCatalogItem
// @Description  Install a marketplace catalog item.
// @Tags         marketplace
// @Accept       json
// @Produce      json
// @Param        InstallMarketplaceCatalogItem 	  body    dto.InstallMarketplaceCatalogItem  true  "InstallMarketplaceCatalogItem"
// @Success      201 {object} object{} "MarketplaceCatalogItemInstalled"
// @Router       /marketplace/catalog/ [post]
func (controller *MarketplaceController) InstallCatalogItemController(c echo.Context) error {
	requiredParams := []string{"id", "hostname", "dataFields"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	id := valueObject.NewMarketplaceItemIdPanic(requestBody["id"])
	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))

	rawRootDir := requestBody["rootDirectory"]
	if rawRootDir == nil {
		rawRootDir = "/"
	}
	rootDir := valueObject.NewUnixFilePathPanic(rawRootDir.(string))

	dataFields := getDataFieldsFromBody(requestBody["dataFields"])

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)
	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(controller.persistentDbSvc)
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	dto := dto.NewInstallMarketplaceCatalogItem(id, hostname, rootDir, dataFields)
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

// GetMarketplaceInstalledItems godoc
// @Summary      GetMarketplaceInstalledItems
// @Description  List marketplace installed items.
// @Tags         marketplace
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {string} entity.MarketplaceInstalledItem
// @Router       /marketplace/installed/ [get]
func (controller *MarketplaceController) GetInstalledItemsController(c echo.Context) error {
	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)
	marketplaceInstalledItems, err := useCase.GetMarketplaceInstalledItems(marketplaceQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, marketplaceInstalledItems)
}
