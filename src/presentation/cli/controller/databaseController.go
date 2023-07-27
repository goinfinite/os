package cliController

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetDatabasesController() *cobra.Command {
	var dbTypeStr string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetDatabases",
		Run: func(cmd *cobra.Command, args []string) {
			dbType := valueObject.NewDatabaseTypePanic(dbTypeStr)

			databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)

			databasesList, err := useCase.GetDatabases(databaseQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
			}

			cliHelper.ResponseWrapper(true, databasesList)
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	return cmd
}

func AddDatabaseController() *cobra.Command {
	var dbTypeStr string
	var dbNameStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewDatabase",
		Run: func(cmd *cobra.Command, args []string) {
			dbType := valueObject.NewDatabaseTypePanic(dbTypeStr)
			dbName := valueObject.NewDatabaseNamePanic(dbNameStr)

			addDatabaseDto := dto.NewAddDatabase(dbName)

			databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)
			databaseCmdRepo := infra.NewDatabaseCmdRepo(dbType)

			err := useCase.AddDatabase(
				databaseQueryRepo,
				databaseCmdRepo,
				addDatabaseDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
			}

			cliHelper.ResponseWrapper(true, "DatabaseAdded")
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

			databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)
			databaseCmdRepo := infra.NewDatabaseCmdRepo(dbType)

			err := useCase.DeleteDatabase(
				databaseQueryRepo,
				databaseCmdRepo,
				dbName,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
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

func AddDatabaseUserController() *cobra.Command {
	var dbTypeStr string
	var dbNameStr string
	var dbUserStr string
	var dbPassStr string
	var privilegesSlice []string

	cmd := &cobra.Command{
		Use:   "addUser",
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

			addDatabaseUserDto := dto.NewAddDatabaseUser(
				dbName,
				dbUser,
				dbPass,
				privileges,
			)

			databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)
			databaseCmdRepo := infra.NewDatabaseCmdRepo(dbType)

			err := useCase.AddDatabaseUser(
				databaseQueryRepo,
				databaseCmdRepo,
				addDatabaseUserDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
			}

			cliHelper.ResponseWrapper(true, "DatabaseUserAdded")
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
		Use:   "deleteUser",
		Short: "DeleteDatabaseUser",
		Run: func(cmd *cobra.Command, args []string) {
			dbType := valueObject.NewDatabaseTypePanic(dbTypeStr)
			dbName := valueObject.NewDatabaseNamePanic(dbNameStr)
			dbUser := valueObject.NewDatabaseUsernamePanic(dbUserStr)

			databaseQueryRepo := infra.NewDatabaseQueryRepo(dbType)
			databaseCmdRepo := infra.NewDatabaseCmdRepo(dbType)

			err := useCase.DeleteDatabaseUser(
				databaseQueryRepo,
				databaseCmdRepo,
				dbName,
				dbUser,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
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
