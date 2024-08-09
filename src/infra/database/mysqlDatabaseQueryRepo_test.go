package databaseInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
)

func TestMysqlDatabaseQueryRepo(t *testing.T) {
	t.Skip("SkipMysqlDatabaseQueryRepoTest")
	testHelpers.LoadEnvVars()

	t.Run("GetDatabases", func(t *testing.T) {
		databasesQueryRepo := MysqlDatabaseQueryRepo{}
		databasesList, err := databasesQueryRepo.Read()
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if len(databasesList) == 0 {
			t.Errorf("Expected: %v, got: %v", "a list of databases", databasesList)
		}
	})
}
