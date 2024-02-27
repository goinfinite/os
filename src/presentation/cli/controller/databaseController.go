package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	databaseInfra "github.com/speedianet/os/src/infra/database"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetDatabasesController() *cobra.Command {
	var dbTypeStr string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetDatabases",
		Run: func(cmd *cobra.Command, args []string) {
			dbType := valueObject.NewDatabaseTypePanic(dbTypeStr)

			databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)

			databasesList, err := useCase.GetDatabases(databaseQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, databasesList)
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	return cmd
}

func CreateDatabaseController() *cobra.Command {
	var dbTypeStr string
	var dbNameStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "CreateNewDatabase",
		Run: func(cmd *cobra.Command, args []string) {
			dbType := valueObject.NewDatabaseTypePanic(dbTypeStr)
			dbName := valueObject.NewDatabaseNamePanic(dbNameStr)

			createDatabaseDto := dto.NewCreateDatabase(dbName)

			databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
			databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

			err := useCase.CreateDatabase(
				databaseQueryRepo,
				databaseCmdRepo,
				createDatabaseDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "DatabaseCreated")
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	cmd.Flags().StringVarP(&dbNameStr, "db-name", "n", "", "DatabaseName")
	cmd.MarkFlagRequired("db-name")
	return cmd
}

func DeleteDatabaseController() *cobra.Command {
	var dbTypeStr string
	var dbNameStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteDatabase",
		Run: func(cmd *cobra.Command, args []string) {
			dbType := valueObject.NewDatabaseTypePanic(dbTypeStr)
			dbName := valueObject.NewDatabaseNamePanic(dbNameStr)

			databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
			databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

			err := useCase.DeleteDatabase(
				databaseQueryRepo,
				databaseCmdRepo,
				dbName,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "DatabaseDeleted")
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	cmd.Flags().StringVarP(&dbNameStr, "db-name", "n", "", "DatabaseName")
	cmd.MarkFlagRequired("db-name")
	return cmd
}

func CreateDatabaseUserController() *cobra.Command {
	var dbTypeStr string
	var dbNameStr string
	var dbUserStr string
	var dbPassStr string
	var privilegesSlice []string

	cmd := &cobra.Command{
		Use:   "add-user",
		Short: "AddNewDatabaseUser",
		Run: func(cmd *cobra.Command, args []string) {
			dbType := valueObject.NewDatabaseTypePanic(dbTypeStr)
			dbName := valueObject.NewDatabaseNamePanic(dbNameStr)
			dbUser := valueObject.NewDatabaseUsernamePanic(dbUserStr)
			dbPass := valueObject.NewPasswordPanic(dbPassStr)

			privileges := []valueObject.DatabasePrivilege{}
			if len(privilegesSlice) > 0 {
				for _, privilege := range privilegesSlice {
					parsedPrivilege := valueObject.NewDatabasePrivilegePanic(privilege)
					privileges = append(privileges, parsedPrivilege)
				}
			}

			createDatabaseUserDto := dto.NewCreateDatabaseUser(
				dbName,
				dbUser,
				dbPass,
				privileges,
			)

			databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
			databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

			err := useCase.CreateDatabaseUser(
				databaseQueryRepo,
				databaseCmdRepo,
				createDatabaseUserDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "DatabaseUserCreated")
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	cmd.Flags().StringVarP(&dbNameStr, "db-name", "n", "", "DatabaseName")
	cmd.MarkFlagRequired("db-name")
	cmd.Flags().StringVarP(&dbUserStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(&dbPassStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringSliceVarP(
		&privilegesSlice,
		"privileges",
		"r",
		[]string{},
		"DatabasePrivileges (Comma-separated)",
	)

	return cmd
}

func DeleteDatabaseUserController() *cobra.Command {
	var dbTypeStr string
	var dbNameStr string
	var dbUserStr string

	cmd := &cobra.Command{
		Use:   "delete-user",
		Short: "DeleteDatabaseUser",
		Run: func(cmd *cobra.Command, args []string) {
			dbType := valueObject.NewDatabaseTypePanic(dbTypeStr)
			dbName := valueObject.NewDatabaseNamePanic(dbNameStr)
			dbUser := valueObject.NewDatabaseUsernamePanic(dbUserStr)

			databaseQueryRepo := databaseInfra.NewDatabaseQueryRepo(dbType)
			databaseCmdRepo := databaseInfra.NewDatabaseCmdRepo(dbType)

			err := useCase.DeleteDatabaseUser(
				databaseQueryRepo,
				databaseCmdRepo,
				dbName,
				dbUser,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "DatabaseUserDeleted")
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	cmd.Flags().StringVarP(&dbNameStr, "db-name", "n", "", "DatabaseName")
	cmd.MarkFlagRequired("db-name")
	cmd.Flags().StringVarP(&dbUserStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	return cmd
}
