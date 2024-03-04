package cliInit

import internalDatabaseInfra "github.com/speedianet/os/src/infra/internalDatabase"

func TransientDatabaseService() *internalDatabaseInfra.TransientDatabaseService {
	transientDbSvc, err := internalDatabaseInfra.NewTransientDatabaseService()
	if err != nil {
		panic("PersistentDatabaseConnectionError:" + err.Error())
	}

	return transientDbSvc
}
