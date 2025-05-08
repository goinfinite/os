package service

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	databaseInfra "github.com/goinfinite/os/src/infra/database"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
)

type DatabaseService struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
	availabilityInspector *sharedHelper.ServiceAvailabilityInspector
}

func NewDatabaseService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *DatabaseService {
	return &DatabaseService{
		persistentDbSvc:       persistentDbSvc,
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
		availabilityInspector: sharedHelper.NewServiceAvailabilityInspector(
			persistentDbSvc,
		),
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
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	requestPagination, err := serviceHelper.PaginationParser(
		input, useCase.DatabasesDefaultPagination,
	)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var databaseNamePtr *valueObject.DatabaseName
	if input["name"] != nil {
		databaseName, err := valueObject.NewDatabaseName(input["name"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		databaseNamePtr = &databaseName
	}

	var usernamePtr *valueObject.DatabaseUsername
	if input["username"] != nil {
		username, err := valueObject.NewDatabaseUsername(input["username"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		usernamePtr = &username
	}

	requestDto := dto.ReadDatabasesRequest{
		Pagination:   requestPagination,
		DatabaseName: databaseNamePtr,
		DatabaseType: &dbType,
		Username:     usernamePtr,
	}

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)

	responseDto, err := useCase.ReadDatabases(databaseQueryRepo, requestDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, responseDto)
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
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateDatabase(dbName, operatorAccountId, operatorIpAddress)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.CreateDatabase(
		databaseQueryRepo, databaseCmdRepo, service.activityRecordCmdRepo, createDto,
	)
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
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	dbName, err := valueObject.NewDatabaseName(input["dbName"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteDatabase(dbName, operatorAccountId, operatorIpAddress)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.DeleteDatabase(
		databaseQueryRepo, databaseCmdRepo, service.activityRecordCmdRepo, deleteDto,
	)
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
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
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

	dbPrivileges := []valueObject.DatabasePrivilege{
		valueObject.DatabasePrivilege("ALL"),
	}
	if input["privileges"] != nil {
		var assertOk bool
		dbPrivileges, assertOk = input["privileges"].([]valueObject.DatabasePrivilege)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidDatabasePrivileges")
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateDatabaseUser(
		dbName, dbUsername, dbPassword, dbPrivileges, operatorAccountId,
		operatorIpAddress,
	)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.CreateDatabaseUser(
		databaseQueryRepo, databaseCmdRepo, service.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "DatabaseUserCreated")
}

func (service *DatabaseService) DeleteUser(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"dbType", "dbName", "dbUser"}
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
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	dbName, err := valueObject.NewDatabaseName(input["dbName"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	dbUsername, err := valueObject.NewDatabaseUsername(input["dbUser"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteDatabaseUser(
		dbName, dbUsername, operatorAccountId, operatorIpAddress,
	)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.DeleteDatabaseUser(
		databaseQueryRepo, databaseCmdRepo, service.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "DatabaseUserDeleted")
}
