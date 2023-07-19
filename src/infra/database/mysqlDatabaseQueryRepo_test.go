package databaseInfra

import (
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
)

func TestMysqlDatabaseQueryRepo(t *testing.T) {
	t.Skip("Skip mysql database query repo test")
	testHelpers.LoadEnvVars()

	t.Run("GetDatabases", func(t *testing.T) {
		databasesQueryRepo := MysqlDatabaseQueryRepo{}
		databasesList, err := databasesQueryRepo.Get()
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if len(databasesList) == 0 {
			t.Errorf("Expected: %v, got: %v", "a list of databases", databasesList)
		}
	})
}
