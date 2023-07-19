package databaseInfra

import (
	"testing"

	testHelpers "github.com/speedianet/sam/src/devUtils"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
	servicesInfra "github.com/speedianet/sam/src/infra/services"
)

func TestMysqlDatabaseCmdRepo(t *testing.T) {
	t.Skip("Skip mysql database cmd repo test")
	testHelpers.LoadEnvVars()

	err := servicesInfra.Install(
		valueObject.NewServiceNamePanic("mysql"),
		nil,
	)
	if err != nil {
		t.Error("Error installing service")
	}
	_, err = infraHelper.RunCmd("mysqld_safe", "&")
	if err != nil {
		t.Error("Error starting command")
	}

	mysqlDatabaseCmdRepo := MysqlDatabaseCmdRepo{}

	t.Run("AddDatabase", func(t *testing.T) {
		err := mysqlDatabaseCmdRepo.Add("testing")
		if err != nil {
			t.Error("Error adding database")
		}
	})

	t.Run("AddDatabaseUser", func(t *testing.T) {
		addDatabaseUserDto := dto.NewAddDatabaseUser(
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
