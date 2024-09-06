package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	databaseInfra "github.com/speedianet/os/src/infra/database"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
)

type DatabaseService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewDatabaseService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *DatabaseService {
	return &DatabaseService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *DatabaseService) Read(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"dbType"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(input["dbType"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}
	err = sharedHelper.EnsureServiceAvailability(service.persistentDbSvc, serviceName)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)

	databasesList, err := useCase.ReadDatabases(databaseQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, databasesList)
}

func (service *DatabaseService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"dbType", "dbName"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(input["dbType"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbName, err := valueObject.NewDatabaseName(input["dbName"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}
	err = sharedHelper.EnsureServiceAvailability(service.persistentDbSvc, serviceName)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	dto := dto.NewCreateDatabase(dbName)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.CreateDatabase(databaseQueryRepo, databaseCmdRepo, dto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "DatabaseCreated")
}

func (service *DatabaseService) Delete(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"dbType", "dbName"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(input["dbType"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}
	err = sharedHelper.EnsureServiceAvailability(service.persistentDbSvc, serviceName)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	dbName, err := valueObject.NewDatabaseName(input["dbName"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.DeleteDatabase(databaseQueryRepo, databaseCmdRepo, dbName)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "DatabaseDeleted")
}

func (service *DatabaseService) CreateUser(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"dbType", "dbName", "username", "password"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(input["dbType"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}
	err = sharedHelper.EnsureServiceAvailability(service.persistentDbSvc, serviceName)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	dbName, err := valueObject.NewDatabaseName(input["dbName"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbUsername, err := valueObject.NewDatabaseUsername(input["username"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbPassword, err := valueObject.NewPassword(input["password"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbPrivileges := []valueObject.DatabasePrivilege{}
	if input["privileges"] != nil {
		for _, rawPrivilege := range input["privileges"].([]string) {
			dbPrivilege, err := valueObject.NewDatabasePrivilege(rawPrivilege)
			if err != nil {
				return NewServiceOutput(UserError, err.Error())
			}
			dbPrivileges = append(dbPrivileges, dbPrivilege)
		}
	}

	dto := dto.NewCreateDatabaseUser(dbName, dbUsername, dbPassword, dbPrivileges)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.CreateDatabaseUser(databaseQueryRepo, databaseCmdRepo, dto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "DatabaseUserCreated")
}

func (service *DatabaseService) DeleteUser(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"dbType", "dbName", "username"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(input["dbType"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}
	err = sharedHelper.EnsureServiceAvailability(service.persistentDbSvc, serviceName)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	dbName, err := valueObject.NewDatabaseName(input["dbName"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbUsername, err := valueObject.NewDatabaseUsername(input["username"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.DeleteDatabaseUser(
		databaseQueryRepo, databaseCmdRepo, dbName, dbUsername,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "DatabaseUserDeleted")
}
