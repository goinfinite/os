package liaison

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	databaseInfra "github.com/goinfinite/os/src/infra/database"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
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

func (liaison *DatabaseLiaison) Read(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"dbType"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, sharedHelper.ServiceUnavailableError)
	}

	requestPagination, err := tkPresentation.PaginationParser(
		useCase.DatabasesDefaultPagination, untrustedInput,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	var databaseNamePtr *valueObject.DatabaseName
	if untrustedInput["name"] != nil {
		databaseName, err := valueObject.NewDatabaseName(untrustedInput["name"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
		databaseNamePtr = &databaseName
	}

	var usernamePtr *valueObject.DatabaseUsername
	if untrustedInput["username"] != nil {
		username, err := valueObject.NewDatabaseUsername(untrustedInput["username"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, responseDto)
}

func (liaison *DatabaseLiaison) Create(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"dbType", "dbName"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbName, err := valueObject.NewDatabaseName(untrustedInput["dbName"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, sharedHelper.ServiceUnavailableError)
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	createDto := dto.NewCreateDatabase(dbName, operatorAccountId, operatorIpAddress)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.CreateDatabase(
		databaseQueryRepo, databaseCmdRepo, liaison.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusCreated, "DatabaseCreated")
}

func (liaison *DatabaseLiaison) Delete(untrustedInput map[string]any) tkPresentation.LiaisonResponse {
	requiredParams := []string{"dbType", "dbName"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, sharedHelper.ServiceUnavailableError)
	}

	dbName, err := valueObject.NewDatabaseName(untrustedInput["dbName"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	deleteDto := dto.NewDeleteDatabase(dbName, operatorAccountId, operatorIpAddress)

	databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
	databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

	err = useCase.DeleteDatabase(
		databaseQueryRepo, databaseCmdRepo, liaison.activityRecordCmdRepo, deleteDto,
	)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, "DatabaseDeleted")
}

func (liaison *DatabaseLiaison) CreateUser(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"dbType", "dbName", "username", "password"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, sharedHelper.ServiceUnavailableError)
	}

	dbName, err := valueObject.NewDatabaseName(untrustedInput["dbName"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbUsername, err := valueObject.NewDatabaseUsername(untrustedInput["username"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbPassword, err := tkValueObject.NewPassword(untrustedInput["password"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbPrivileges := []valueObject.DatabasePrivilege{
		valueObject.DatabasePrivilege("ALL"),
	}
	if untrustedInput["privileges"] != nil {
		var assertOk bool
		dbPrivileges, assertOk = untrustedInput["privileges"].([]valueObject.DatabasePrivilege)
		if !assertOk {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, "InvalidDatabasePrivileges")
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusCreated, "DatabaseUserCreated")
}

func (liaison *DatabaseLiaison) DeleteUser(
	untrustedInput map[string]any,
) tkPresentation.LiaisonResponse {
	requiredParams := []string{"dbType", "dbName", "dbUser"}
	err := tkPresentation.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbType, err := valueObject.NewDatabaseType(untrustedInput["dbType"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	serviceName, err := valueObject.NewServiceName(dbType.String())
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}
	if !liaison.availabilityInspector.IsAvailable(serviceName) {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, sharedHelper.ServiceUnavailableError)
	}

	dbName, err := valueObject.NewDatabaseName(untrustedInput["dbName"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	dbUsername, err := valueObject.NewDatabaseUsername(untrustedInput["dbUser"])
	if err != nil {
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = tkValueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = tkValueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusUserError, err.Error())
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
		return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusInfraError, err.Error())
	}

	return tkPresentation.NewLiaisonResponseNoMessage(tkPresentation.LiaisonResponseStatusSuccess, "DatabaseUserDeleted")
}
