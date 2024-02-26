package cliInit

import databaseInfra "github.com/speedianet/os/src/infra/database"

func InternalDatabaseService() *databaseInfra.InternalDatabaseService {
	internalDbSvc, err := databaseInfra.NewInternalDatabaseService()
	if err != nil {
		panic("InternalDatabaseConnectionError:" + err.Error())
	}

	return internalDbSvc
}
