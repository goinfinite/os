package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
	return apiHelper.ResponseWrapper(c, http.StatusOK, "AllCatalogProducts")
}
