package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	databaseInfra "github.com/speedianet/os/src/infra/database"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
)

// GetDatabases godoc
// @Summary      GetDatabases
// @Description  List databases names, users and sizes.
// @Tags         database
// @Security     Bearer
// @Accept       json
// @Produce      json
// @Param        dbType path string true "DatabaseType (like mysql, postgres)"
// @Success      200 {array} entity.Database
// @Router       /v1/database/{dbType}/ [get]
func GetDatabasesController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))
	sharedHelper.StopIfServiceUnavailable(dbType.String())

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)

	databasesList, err := useCase.GetDatabases(databaseQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, databasesList)
}

// CreateDatabase godoc
// @Summary      CreateDatabase
// @Description  Create a new database.
// @Tags         database
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        dbType path string true "DatabaseType (like mysql, postgres)"
// @Param        createDatabaseDto body dto.CreateDatabase true "All props are required."
// @Success      201 {object} object{} "DatabaseCreated"
// @Router       /v1/database/{dbType}/ [post]
func CreateDatabaseController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))
	sharedHelper.StopIfServiceUnavailable(dbType.String())

	requiredParams := []string{"dbName"}
	requestBody, _ := apiHelper.GetRequestBody(c)
	apiHelper.CheckMissingParams(requestBody, requiredParams)

	dbName := valueObject.NewDatabaseNamePanic(requestBody["dbName"].(string))
	createDatabaseDto := dto.NewCreateDatabase(dbName)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err := useCase.CreateDatabase(
		databaseQueryRepo,
		databaseCmdRepo,
		createDatabaseDto,
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
// @Param        dbType path string true "DatabaseType (like mysql, postgres)"
// @Param        dbName path string true "DatabaseName"
// @Success      200 {object} object{} "DatabaseDeleted"
// @Router       /v1/database/{dbType}/{dbName}/ [delete]
func DeleteDatabaseController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))
	sharedHelper.StopIfServiceUnavailable(dbType.String())

	dbName := valueObject.NewDatabaseNamePanic(c.Param("dbName"))

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

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

// CreateDatabaseUser godoc
// @Summary      CreateDatabaseUser
// @Description  Create a new database user.
// @Tags         database
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        dbType path string true "DatabaseType (like mysql, postgres)"
// @Param        dbName path string true "DatabaseName"
// @Param        createDatabaseUserDto body dto.CreateDatabaseUser true "privileges is optional. When not provided, privileges will be 'ALL'."
// @Success      201 {object} object{} "DatabaseUserCreated"
// @Router       /v1/database/{dbType}/{dbName}/user/ [post]
func CreateDatabaseUserController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))
	sharedHelper.StopIfServiceUnavailable(dbType.String())

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

	createDatabaseUserDto := dto.NewCreateDatabaseUser(
		dbName,
		username,
		password,
		privileges,
	)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err := useCase.CreateDatabaseUser(
		databaseQueryRepo,
		databaseCmdRepo,
		createDatabaseUserDto,
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
// @Param        dbType path string true "DatabaseType (like mysql, postgres)"
// @Param        dbName path string true "DatabaseName"
// @Param        dbUser path string true "Database username that will be deleted."
// @Success      200 {object} object{} "DatabaseUserDeleted"
// @Router       /v1/database/{dbType}/{dbName}/user/{dbUser}/ [delete]
func DeleteDatabaseUserController(c echo.Context) error {
	dbType := valueObject.NewDatabaseTypePanic(c.Param("dbType"))
	sharedHelper.StopIfServiceUnavailable(dbType.String())

	dbName := valueObject.NewDatabaseNamePanic(c.Param("dbName"))
	dbUser := valueObject.NewDatabaseUsernamePanic(c.Param("dbUser"))

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

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
