package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/sam/src/domain/dto"
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

	databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)

	databasesList, err := useCase.GetDatabases(databaseQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, databasesList)
}

// AddDatabase godoc
// @Summary      AddDatabase
// @Description  Add a new database.
// @Tags         database
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        dbType path valueObject.DatabaseType true "DatabaseType"
// @Param        addDatabaseDto body dto.AddDatabase true "AddDatabase"
// @Success      201 {object} object{} "DatabaseAdded"
// @Router       /database/{dbType}/ [post]
func AddDatabaseController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))

	requiredParams := []string{"dbName"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)
	dbName := valueObject.NewDatabaseNamePanic(requestBody["dbName"].(string))
	addDatabaseDto := dto.NewAddDatabase(dbName)

	databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := infra.NewDatabaseCmdRepo(dbType)

	err := useCase.AddDatabase(
		databaseQueryRepo,
		databaseCmdRepo,
		addDatabaseDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, nil)
}
