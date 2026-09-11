package cliInit

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
)

func TransientDatabaseService() *tkInfraDb.TransientDatabaseService {
	transientDbSvc, err := tkInfraDb.NewTransientDatabaseService()
	if err != nil {
		panic("TransientDatabaseConnectionError: " + err.Error())
	}

	return transientDbSvc
}

func PersistentDatabaseService() *internalDbInfra.PersistentDatabaseService {
	persistentDbSvc, err := internalDbInfra.NewPersistentDatabaseService()
	if err != nil {
		panic("PersistentDatabaseConnectionError: " + err.Error())
	}

	return persistentDbSvc
}

func TrailDatabaseService() *internalDbInfra.TrailDatabaseService {
	trailDbSvc, err := internalDbInfra.NewTrailDatabaseService()
	if err != nil {
		panic("TrailDatabaseConnectionError:" + err.Error())
	}

	return trailDbSvc
}
