package cliInit

import internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"

func TransientDatabaseService() *internalDbInfra.TransientDatabaseService {
	transientDbSvc, err := internalDbInfra.NewTransientDatabaseService()
	if err != nil {
		panic("TransientDatabaseConnectionError:" + err.Error())
	}

	return transientDbSvc
}

func PersistentDatabaseService() *internalDbInfra.PersistentDatabaseService {
	persistentDbSvc, err := internalDbInfra.NewPersistentDatabaseService()
	if err != nil {
		panic("PersistentDatabaseConnectionError:" + err.Error())
	}

	return persistentDbSvc
}
