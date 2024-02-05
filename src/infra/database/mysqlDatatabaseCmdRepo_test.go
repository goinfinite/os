package databaseInfra

import (
	"testing"

	testHelpers "github.com/speedianet/os/src/devUtils"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

func TestMysqlDatabaseCmdRepo(t *testing.T) {
	t.Skip("SkipMysqlDatabaseCmdRepoTest")
	testHelpers.LoadEnvVars()

	err := servicesInfra.AddInstallableSimplified("mariadb")
	if err != nil {
		t.Errorf("InstallDependenciesFail: %v", err)
		return
	}

	_, err = infraHelper.RunCmd("mysqld_safe", "&")
	if err != nil {
		t.Error("Error starting command")
	}

	mysqlDatabaseCmdRepo := MysqlDatabaseCmdRepo{}

	t.Run("CreateDatabase", func(t *testing.T) {
		err := mysqlDatabaseCmdRepo.Add("testing")
		if err != nil {
			t.Error("Error adding database")
		}
	})

	t.Run("CreateDatabaseUser", func(t *testing.T) {
		addDatabaseUserDto := dto.NewCreateDatabaseUser(
			valueObject.NewDatabaseNamePanic("testing"),
			valueObject.NewDatabaseUsernamePanic("testing"),
			valueObject.NewPasswordPanic("testing"),
			[]valueObject.DatabasePrivilege{
				valueObject.NewDatabasePrivilegePanic("ALL"),
			},
		)

		err := mysqlDatabaseCmdRepo.AddUser(addDatabaseUserDto)
		if err != nil {
			t.Error("Error adding database user")
		}
	})

	t.Run("DeleteDatabaseUser", func(t *testing.T) {
		dbName := valueObject.NewDatabaseNamePanic("testing")
		dbUsername := valueObject.NewDatabaseUsernamePanic("testing")

		err := mysqlDatabaseCmdRepo.DeleteUser(dbName, dbUsername)
		if err != nil {
			t.Error("Error removing database user")
		}
	})

	t.Run("DeleteDatabase", func(t *testing.T) {
		err := mysqlDatabaseCmdRepo.Delete("testing")
		if err != nil {
			t.Error("Error removing database")
		}
	})

	err = servicesInfra.Uninstall(
		valueObject.NewServiceNamePanic("mysql"),
	)
	if err != nil {
		t.Error("Error uninstalling service")
	}
}
