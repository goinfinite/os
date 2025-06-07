package liaison

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	databaseInfra "github.com/goinfinite/os/src/infra/database"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
)

type DatabaseLiaison struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
	availabilityInspector *sharedHelper.ServiceAvailabilityInspector
}

func NewDatabaseLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *DatabaseLiaison {
	return &DatabaseLiaison{
		persistentDbSvc:       persistentDbSvc,
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
		availabilityInspector: sharedHelper.NewServiceAvailabilityInspector(
			persistentDbSvc,
		),
	}
}

func (liaison *DatabaseLiaison) Read(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"dbType"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return NewLiaisonOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	requestPagination, err := liaisonHelper.PaginationParser(
		untrustedInput, useCase.DatabasesDefaultPagination,
	)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	var databaseNamePtr *valueObject.DatabaseName
	if untrustedInput["name"] != nil {
		databaseName, err := valueObject.NewDatabaseName(untrustedInput["name"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		databaseNamePtr = &databaseName
	}

	var usernamePtr *valueObject.DatabaseUsername
	if untrustedInput["username"] != nil {
		username, err := valueObject.NewDatabaseUsername(untrustedInput["username"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
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
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, responseDto)
}

func (liaison *DatabaseLiaison) Create(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"dbType", "dbName"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbName, err := valueObject.NewDatabaseName(untrustedInput["dbName"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return NewLiaisonOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateDatabase(dbName, operatorAccountId, operatorIpAddress)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.CreateDatabase(
		databaseQueryRepo, databaseCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "DatabaseCreated")
}

func (liaison *DatabaseLiaison) Delete(untrustedInput map[string]any) LiaisonOutput {
	requiredParams := []string{"dbType", "dbName"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return NewLiaisonOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	dbName, err := valueObject.NewDatabaseName(untrustedInput["dbName"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteDatabase(dbName, operatorAccountId, operatorIpAddress)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.DeleteDatabase(
		databaseQueryRepo, databaseCmdRepo, liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "DatabaseDeleted")
}

func (liaison *DatabaseLiaison) CreateUser(
	untrustedInput map[string]any,
) LiaisonOutput {
	requiredParams := []string{"dbType", "dbName", "username", "password"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return NewLiaisonOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	dbName, err := valueObject.NewDatabaseName(untrustedInput["dbName"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbUsername, err := valueObject.NewDatabaseUsername(untrustedInput["username"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbPassword, err := valueObject.NewPassword(untrustedInput["password"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbPrivileges := []valueObject.DatabasePrivilege{
		valueObject.DatabasePrivilege("ALL"),
	}
	if untrustedInput["privileges"] != nil {
		var assertOk bool
		dbPrivileges, assertOk = untrustedInput["privileges"].([]valueObject.DatabasePrivilege)
		if !assertOk {
			return NewLiaisonOutput(UserError, "InvalidDatabasePrivileges")
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	createDto := dto.NewCreateDatabaseUser(
		dbName, dbUsername, dbPassword, dbPrivileges, operatorAccountId,
		operatorIpAddress,
	)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.CreateDatabaseUser(
		databaseQueryRepo, databaseCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "DatabaseUserCreated")
}

func (liaison *DatabaseLiaison) DeleteUser(
	untrustedInput map[string]any,
) LiaisonOutput {
	requiredParams := []string{"dbType", "dbName", "dbUser"}
	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return NewLiaisonOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	dbName, err := valueObject.NewDatabaseName(untrustedInput["dbName"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	dbUsername, err := valueObject.NewDatabaseUsername(untrustedInput["dbUser"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteDatabaseUser(
		dbName, dbUsername, operatorAccountId, operatorIpAddress,
	)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.DeleteDatabaseUser(
		databaseQueryRepo, databaseCmdRepo, liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "DatabaseUserDeleted")
}
