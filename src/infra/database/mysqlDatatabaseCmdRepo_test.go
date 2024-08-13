package databaseInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

func TestMysqlDatabaseCmdRepo(t *testing.T) {
	t.Skip("SkipMysqlDatabaseCmdRepoTest")
	testHelpers.LoadEnvVars()

	_, err := infraHelper.RunCmd("mysqld_safe", "&")
	if err != nil {
		t.Error("Error starting command")
	}

	dbName, _ := valueObject.NewDatabaseName("testing")
	dbUsername, _ := valueObject.NewDatabaseUsername("testing")
	dbPassword, _ := valueObject.NewPassword("testing")
	dbPrivilege, _ := valueObject.NewDatabasePrivilege("ALL")
	dbPrivileges := []valueObject.DatabasePrivilege{dbPrivilege}

	mysqlDatabaseCmdRepo := MysqlDatabaseCmdRepo{}

	t.Run("CreateDatabase", func(t *testing.T) {
		err := mysqlDatabaseCmdRepo.Create(dbName)
		if err != nil {
			t.Error("Error creating database")
		}
	})

	t.Run("CreateDatabaseUser", func(t *testing.T) {
		createDatabaseUserDto := dto.NewCreateDatabaseUser(
			dbName, dbUsername, dbPassword, dbPrivileges,
		)

		err := mysqlDatabaseCmdRepo.CreateUser(createDatabaseUserDto)
		if err != nil {
			t.Error("Error creating database user")
		}
	})

	t.Run("DeleteDatabaseUser", func(t *testing.T) {
		err := mysqlDatabaseCmdRepo.DeleteUser(dbName, dbUsername)
		if err != nil {
			t.Error("Error removing database user")
		}
	})

	t.Run("DeleteDatabase", func(t *testing.T) {
		err := mysqlDatabaseCmdRepo.Delete(dbName)
		if err != nil {
			t.Error("Error removing database")
		}
	})
}
