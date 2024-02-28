package cliInit

import databaseInfra "github.com/speedianet/os/src/infra/database"

func TransientDatabaseService() *databaseInfra.TransientDatabaseService {
	transientDbSvc, err := databaseInfra.NewTransientDatabaseService()
	if err != nil {
		panic("PersistentDatabaseConnectionError:" + err.Error())
	}

	return transientDbSvc
}
