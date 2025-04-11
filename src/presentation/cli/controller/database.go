package cliController

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/service"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/spf13/cobra"
)

type DatabaseController struct {
	persistentDbService *internalDbInfra.PersistentDatabaseService
	dbService           *service.DatabaseService
}

func NewDatabaseController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *DatabaseController {
	return &DatabaseController{
		persistentDbService: persistentDbService,
		dbService: service.NewDatabaseService(
			persistentDbService, trailDbSvc,
		),
	}
}

func (controller *DatabaseController) Read() *cobra.Command {
	var dbTypeStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadDatabases",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"dbType": dbTypeStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.dbService.Read(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	return cmd
}

func (controller *DatabaseController) Create() *cobra.Command {
	var dbTypeStr, dbNameStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewDatabase",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"dbType": dbTypeStr,
				"dbName": dbNameStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.dbService.Create(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	cmd.Flags().StringVarP(&dbNameStr, "db-name", "n", "", "DatabaseName")
	cmd.MarkFlagRequired("db-name")
	return cmd
}

func (controller *DatabaseController) Delete() *cobra.Command {
	var dbTypeStr, dbNameStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteDatabase",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"dbType": dbTypeStr,
				"dbName": dbNameStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.dbService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	cmd.Flags().StringVarP(&dbNameStr, "db-name", "n", "", "DatabaseName")
	cmd.MarkFlagRequired("db-name")
	return cmd
}

func (controller *DatabaseController) CreateUser() *cobra.Command {
	var dbTypeStr, dbNameStr, dbUserStr, dbPassStr string
	var privilegesSlice []string

	cmd := &cobra.Command{
		Use:   "create-user",
		Short: "CreateNewDatabaseUser",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"dbType":   dbTypeStr,
				"dbName":   dbNameStr,
				"username": dbUserStr,
				"password": dbPassStr,
			}

			if len(privilegesSlice) > 0 {
				requestBody["privileges"] = tkPresentation.StringSliceValueObjectParser(
					privilegesSlice, valueObject.NewDatabasePrivilege,
				)
			}

			cliHelper.ServiceResponseWrapper(
				controller.dbService.CreateUser(requestBody),
			)
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
		&privilegesSlice, "privileges", "r", []string{},
		"DatabasePrivileges (Comma or semicolon separated)",
	)

	return cmd
}

func (controller *DatabaseController) DeleteUser() *cobra.Command {
	var dbTypeStr, dbNameStr, dbUsernameStr string

	cmd := &cobra.Command{
		Use:   "delete-user",
		Short: "DeleteDatabaseUser",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"dbType": dbTypeStr,
				"dbName": dbNameStr,
				"dbUser": dbUsernameStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.dbService.DeleteUser(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	cmd.Flags().StringVarP(&dbNameStr, "db-name", "n", "", "DatabaseName")
	cmd.MarkFlagRequired("db-name")
	cmd.Flags().StringVarP(&dbUsernameStr, "db-username", "u", "", "DatabaseUsername")
	cmd.MarkFlagRequired("db-username")
	return cmd
}
