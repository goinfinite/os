package databaseInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestMysqlDatabaseCmdRepo(t *testing.T) {
	t.Skip("SkipMysqlDatabaseCmdRepoTest")
	testHelpers.LoadEnvVars()

	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "mysqld_safe",
		Args:    []string{"&"},
	}).Run()
	if err != nil {
		t.Error("Error starting command")
	}

	dbName, _ := valueObject.NewDatabaseName("testing")
	dbUsername, _ := valueObject.NewDatabaseUsername("testing")
	dbPassword, _ := tkValueObject.NewPassword("Testing@1")
	dbPrivilege, _ := valueObject.NewDatabasePrivilege("ALL")
	dbPrivileges := []valueObject.DatabasePrivilege{dbPrivilege}

	ipAddress := tkValueObject.IpAddressLocal
	operatorAccountId, _ := tkValueObject.NewAccountId(0)

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
			operatorAccountId, ipAddress,
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
