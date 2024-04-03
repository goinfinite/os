package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	mktplaceInfra "github.com/speedianet/os/src/infra/marketplace"
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
// @Success      200 {string} "AllCatalogProducts"
// @Router       /marketplace/catalog/ [get]
func (controller *MarketplaceController) GetCatalogController(c echo.Context) error {
	mktplaceCatalogQueryRepo := mktplaceInfra.NewMktplaceCatalogQueryRepo(controller.persistentDbSvc)
	mktplaceItems, err := useCase.GetMarketplaceCatalog(mktplaceCatalogQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, mktplaceItems)
}

func getDataFieldsFromBody(
	dataFieldsBodyInput interface{},
) []valueObject.DataField {
	dataFields := []valueObject.DataField{}

	dataFieldsInterfaceSlice, assertOk := dataFieldsBodyInput.([]interface{})
	if !assertOk {
		panic("InvalidDataField")
	}

	for _, dataFieldsInterface := range dataFieldsInterfaceSlice {
		dataFieldMap, assertOk := dataFieldsInterface.(map[string]interface{})
		if !assertOk {
			panic("InvalidDataField")
		}

		dataField := valueObject.NewDataField(
			valueObject.NewDataFieldKeyPanic(dataFieldMap["key"].(string)),
			valueObject.NewDataFieldValuePanic(dataFieldMap["value"].(string)),
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

	mktplaceItemId := valueObject.NewMktplaceItemIdPanic(requestBody["id"])
	hostname := valueObject.NewFqdnPanic(requestBody["hostname"].(string))

	rawRootDir := requestBody["rootDirectory"]
	if rawRootDir == nil {
		rawRootDir = "/"
	}
	rootDir := valueObject.NewUnixFilePathPanic(rawRootDir.(string))

	dataFields := getDataFieldsFromBody(requestBody["dataFields"])

	mktplaceCatalogQueryRepo := mktplaceInfra.NewMktplaceCatalogQueryRepo(controller.persistentDbSvc)
	mktplaceCatalogCmdRepo := mktplaceInfra.NewMktplaceCatalogCmdRepo(controller.persistentDbSvc)
	vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}
	vhostCmdRepo := vhostInfra.VirtualHostCmdRepo{}

	dto := dto.NewInstallMarketplaceCatalogItem(mktplaceItemId, hostname, rootDir, dataFields)
	err := useCase.InstallMarketplaceCatalogItem(
		mktplaceCatalogQueryRepo,
		mktplaceCatalogCmdRepo,
		vhostQueryRepo,
		vhostCmdRepo,
		dto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "MarketplaceCatalogItemInstalled")
}
