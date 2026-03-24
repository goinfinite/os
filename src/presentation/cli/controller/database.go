package cliController

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/spf13/cobra"
)

type DatabaseController struct {
	persistentDbService *internalDbInfra.PersistentDatabaseService
	databaseLiaison     *liaison.DatabaseLiaison
}

func NewDatabaseController(
	persistentDbService *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *DatabaseController {
	return &DatabaseController{
		persistentDbService: persistentDbService,
		databaseLiaison: liaison.NewDatabaseLiaison(
			persistentDbService, trailDbSvc,
		),
	}
}

func (controller *DatabaseController) Read() *cobra.Command {
	var (
		dbTypeStr, dbNameStr, usernameStr                                        string
		paginationPageNumberUint32                                               uint32
		paginationItemsPerPageUint16                                             uint16
		paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadDatabases",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"dbType": dbTypeStr,
			}

			if dbNameStr != "" {
				requestBody["name"] = dbNameStr
			}

			if usernameStr != "" {
				requestBody["username"] = usernameStr
			}

			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}
			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}
			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}
			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}
			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			tkPresentation.LiaisonCliResponseRenderer(
				controller.databaseLiaison.Read(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&dbTypeStr, "db-type", "t", "", "DatabaseType")
	cmd.MarkFlagRequired("db-type")
	cmd.Flags().StringVarP(&dbNameStr, "name", "n", "", "DatabaseName")
	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "DatabaseUsername")
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "o", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "j", 0, "ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "x", "", "SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)
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

			tkPresentation.LiaisonCliResponseRenderer(
				controller.databaseLiaison.Create(requestBody),
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

			tkPresentation.LiaisonCliResponseRenderer(
				controller.databaseLiaison.Delete(requestBody),
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

			tkPresentation.LiaisonCliResponseRenderer(
				controller.databaseLiaison.CreateUser(requestBody),
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

			tkPresentation.LiaisonCliResponseRenderer(
				controller.databaseLiaison.DeleteUser(requestBody),
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
