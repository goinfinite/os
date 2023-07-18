package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	apiHelper "github.com/speedianet/sam/src/presentation/api/helper"
)

// GetDatabases	 godoc
// @Summary      GetDatabases
// @Description  List databases names, users and sizes.
// @Tags         database
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        dbType path valueObject.DatabaseType true "DatabaseType"
// @Success      200 {array} entity.Database
// @Router       /database/{dbType}/ [get]
func GetDatabasesController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))

	databasesQueryRepo := infra.NewDatabaseQueryRepo(dbType)

	databasesList, err := useCase.GetDatabases(databasesQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, databasesList)
}
