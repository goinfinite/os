package apiController

import (
	"github.com/labstack/echo/v4"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	"github.com/speedianet/os/src/presentation/service"
)

type RuntimeController struct {
	runtimeService *service.RuntimeService
}

func NewRuntimeController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
) *RuntimeController {
	return &RuntimeController{
		runtimeService: service.NewRuntimeService(persistentDbService),
	}
}

// ReadPhpConfigs godoc
// @Summary      ReadPhpConfigs
// @Description  Get php version, modules and settings for a hostname.
// @Tags         runtime
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        hostname 	  path   string  true  "Hostname"
// @Success      200 {object} entity.PhpConfigs
// @Router       /v1/runtime/php/{hostname}/ [get]
func (controller *RuntimeController) ReadPhpConfigs(c echo.Context) error {
	requestBody := map[string]interface{}{
		"hostname": c.Param("hostname"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.runtimeService.ReadPhpConfigs(requestBody),
	)
}

// UpdatePhpConfigs godoc
// @Summary      UpdatePhpConfigs
// @Description  Update php version, modules and settings for a hostname.
// @Tags         runtime
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        hostname 	  path   string  true  "Hostname"
// @Param        updatePhpConfigsDto	body dto.UpdatePhpConfigs	true	"modules and settings are optional."
// @Success      200 {object} object{} "PhpConfigsUpdated"
// @Router       /v1/runtime/php/{hostname}/ [put]
func (controller *RuntimeController) UpdatePhpConfigs(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}
	requestBody["hostname"] = c.Param("hostname")

	return apiHelper.ServiceResponseWrapper(
		c, controller.runtimeService.UpdatePhpConfigs(requestBody),
	)
}
