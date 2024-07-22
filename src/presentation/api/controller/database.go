package apiController

import (
	"github.com/labstack/echo/v4"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	"github.com/speedianet/os/src/presentation/service"
)

type DatabaseController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	databaseService *service.DatabaseService
}

func NewDatabaseController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *DatabaseController {
	return &DatabaseController{
		persistentDbSvc: persistentDbSvc,
		databaseService: service.NewDatabaseService(persistentDbSvc),
	}
}

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
func (controller *DatabaseController) Read(c echo.Context) error {
	requestBody := map[string]interface{}{
		"dbType": c.Param("dbType"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.databaseService.Read(requestBody),
	)
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
func (controller *DatabaseController) Create(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}
	requestBody["dbType"] = c.Param("dbType")

	return apiHelper.ServiceResponseWrapper(
		c, controller.databaseService.Create(requestBody),
	)
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
func (controller *DatabaseController) Delete(c echo.Context) error {
	requestBody := map[string]interface{}{
		"dbType": c.Param("dbType"),
		"dbName": c.Param("dbName"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.databaseService.Delete(requestBody),
	)
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
func (controller *DatabaseController) CreateUser(c echo.Context) error {
	requestBody, err := apiHelper.ReadRequestBody(c)
	if err != nil {
		return err
	}
	requestBody["dbType"] = c.Param("dbType")
	requestBody["dbName"] = c.Param("dbName")

	privilegesSlice := []string{}
	for _, rawPrivilege := range requestBody["privileges"].([]interface{}) {
		privilegesSlice = append(privilegesSlice, rawPrivilege.(string))
	}
	requestBody["privileges"] = privilegesSlice

	return apiHelper.ServiceResponseWrapper(
		c, controller.databaseService.CreateUser(requestBody),
	)
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
// @Param        dbUser path string true "DatabaseUsername to delete."
// @Success      200 {object} object{} "DatabaseUserDeleted"
// @Router       /v1/database/{dbType}/{dbName}/user/{dbUser}/ [delete]
func (controller *DatabaseController) DeleteUser(c echo.Context) error {
	requestBody := map[string]interface{}{
		"dbType":   c.Param("dbType"),
		"dbName":   c.Param("dbName"),
		"username": c.Param("dbUser"),
	}

	return apiHelper.ServiceResponseWrapper(
		c, controller.databaseService.DeleteUser(requestBody),
	)
}
