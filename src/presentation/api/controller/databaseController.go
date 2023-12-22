package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// GetDatabases godoc
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
// @Success      201 {object} object{} "DatabaseCreated"
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

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "DatabaseCreated")
}

// DeleteDatabase godoc
// @Summary      DeleteDatabase
// @Description  Delete a database.
// @Tags         database
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        dbType path valueObject.DatabaseType true "DatabaseType"
// @Param        dbName path string true "DatabaseName"
// @Success      200 {object} object{} "DatabaseDeleted"
// @Router       /database/{dbType}/{dbName}/ [delete]
func DeleteDatabaseController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))
	dbName := valueObject.NewDatabaseNamePanic(c.Param("dbName"))

	databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := infra.NewDatabaseCmdRepo(dbType)

	err := useCase.DeleteDatabase(
		databaseQueryRepo,
		databaseCmdRepo,
		dbName,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "DatabaseDeleted")
}

// AddDatabaseUser godoc
// @Summary      AddDatabaseUser
// @Description  Add a new database user.
// @Tags         database
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        dbType path valueObject.DatabaseType true "DatabaseType"
// @Param        dbName path string true "DatabaseName"
// @Param        addDatabaseUserDto body dto.AddDatabaseUser true "AddDatabaseUser"
// @Success      201 {object} object{} "DatabaseUserCreated"
// @Router       /database/{dbType}/{dbName}/user/ [post]
func AddDatabaseUserController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))
	dbName := valueObject.NewDatabaseNamePanic(c.Param("dbName"))

	requiredParams := []string{"username", "password"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)
	username := valueObject.NewDatabaseUsernamePanic(requestBody["username"].(string))
	password := valueObject.NewPasswordPanic(requestBody["password"].(string))

	privileges := []valueObject.DatabasePrivilege{}
	if requestBody["privileges"] != nil {
		for _, privilege := range requestBody["privileges"].([]interface{}) {
			privilege := valueObject.NewDatabasePrivilegePanic(privilege.(string))
			privileges = append(privileges, privilege)
		}
	}

	addDatabaseUserDto := dto.NewAddDatabaseUser(
		dbName,
		username,
		password,
		privileges,
	)

	databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := infra.NewDatabaseCmdRepo(dbType)

	err := useCase.AddDatabaseUser(
		databaseQueryRepo,
		databaseCmdRepo,
		addDatabaseUserDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "DatabaseUserCreated")
}

// DeleteDatabaseUser godoc
// @Summary      DeleteDatabaseUser
// @Description  Delete a database user.
// @Tags         database
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        dbType path valueObject.DatabaseType true "DatabaseType"
// @Param        dbName path string true "DatabaseName"
// @Param        dbUser path string true "DatabaseUsername"
// @Success      200 {object} object{} "DatabaseUserDeleted"
// @Router       /database/{dbType}/{dbName}/user/{dbUser}/ [delete]
func DeleteDatabaseUserController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))
	dbName := valueObject.NewDatabaseNamePanic(c.Param("dbName"))
	dbUser := valueObject.NewDatabaseUsernamePanic(c.Param("dbUser"))

	databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := infra.NewDatabaseCmdRepo(dbType)

	err := useCase.DeleteDatabaseUser(
		databaseQueryRepo,
		databaseCmdRepo,
		dbName,
		dbUser,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "DatabaseUserDeleted")
}
