package databaseInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
)

func TestMysqlDatabaseQueryRepo(t *testing.T) {
	t.Skip("SkipMysqlDatabaseQueryRepoTest")
	testHelpers.LoadEnvVars()

	t.Run("GetDatabases", func(t *testing.T) {
		databasesQueryRepo := MysqlDatabaseQueryRepo{}
		allDatabases, err := databasesQueryRepo.readAllDatabases()
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if len(allDatabases) == 0 {
			t.Errorf("Expected: %v, got: %v", "a list of databases", allDatabases)
		}
	})
}
