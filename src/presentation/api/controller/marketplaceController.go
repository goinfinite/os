package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/useCase"
	mktplaceInfra "github.com/speedianet/os/src/infra/marketplace"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// GetMarketplaceCatalog godoc
// @Summary      GetMarketplaceCatalog
// @Description  List marketplace catalog services names, types, steps and more.
// @Tags         marketplace
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Success      200 {string} "AllCatalogProducts"
// @Router       /marketplace/catalog/ [get]
func GetCatalogController(c echo.Context) error {
	mktplaceQueryRepo := mktplaceInfra.MarketplaceQueryRepo{}
	mktplaceItems, err := useCase.GetMarketplaceCatalog(mktplaceQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, mktplaceItems)
}
