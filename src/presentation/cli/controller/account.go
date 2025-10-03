package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/spf13/cobra"
)

type AccountController struct {
	accountLiaison *liaison.AccountLiaison
}

func NewAccountController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountController {
	return &AccountController{
		accountLiaison: liaison.NewAccountLiaison(persistentDbSvc, trailDbSvc),
	}
}

func (controller *AccountController) Read() *cobra.Command {
	var accountIdUint64 uint64
	var accountUsernameStr, shouldIncludeSecureAccessPublicKeysStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadAccounts",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"shouldIncludeSecureAccessPublicKeys": shouldIncludeSecureAccessPublicKeysStr,
			}

			if accountIdUint64 != 0 {
				requestBody["id"] = accountIdUint64
			}

			if accountUsernameStr != "" {
				requestBody["username"] = accountUsernameStr
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

			cliHelper.LiaisonResponseWrapper(
				controller.accountLiaison.Read(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint64, "account-id", "i", 0, "AccountId")
	cmd.Flags().StringVarP(
		&accountUsernameStr, "account-username", "n", "", "AccountUsername",
	)
	cmd.Flags().StringVarP(
		&shouldIncludeSecureAccessPublicKeysStr,
		"should-include-secure-access-public-keys", "s", "false",
		"ShouldIncludeSecureAccessPublicKeys",
	)
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "p", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "m", 0,
		"ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "r", "",
		"SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *AccountController) Create() *cobra.Command {
	var usernameStr, passwordStr, isSuperAdminStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewAccount",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"username":     usernameStr,
				"password":     passwordStr,
				"isSuperAdmin": isSuperAdminStr,
			}

			cliHelper.LiaisonResponseWrapper(
				controller.accountLiaison.Create(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringVarP(
		&isSuperAdminStr, "is-super-admin", "s", "false", "IsSuperAdmin",
	)
	return cmd
}

func (controller *AccountController) Update() *cobra.Command {
	var accountIdUint64 uint64
	var usernameStr, passwordStr, isSuperAdminStr, shouldUpdateApiKeyStr string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateAccount (pass or apiKey)",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"shouldUpdateApiKey": shouldUpdateApiKeyStr,
			}

			if accountIdUint64 != 0 {
				requestBody["accountId"] = accountIdUint64
			}

			if usernameStr != "" {
				requestBody["accountUsername"] = usernameStr
			}

			if passwordStr != "" {
				requestBody["password"] = passwordStr
			}

			if isSuperAdminStr != "" {
				requestBody["isSuperAdmin"] = isSuperAdminStr
			}

			cliHelper.LiaisonResponseWrapper(
				controller.accountLiaison.Update(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint64, "account-id", "i", 0, "AccountId")
	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.Flags().StringVarP(&isSuperAdminStr, "is-super-admin", "s", "", "IsSuperAdmin")
	cmd.Flags().StringVarP(
		&shouldUpdateApiKeyStr, "update-api-key", "k", "false", "ShouldUpdateApiKey",
	)
	return cmd
}

func (controller *AccountController) Delete() *cobra.Command {
	var accountIdUint64 uint64

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteAccount",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint64,
			}

			cliHelper.LiaisonResponseWrapper(
				controller.accountLiaison.Delete(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint64, "account-id", "i", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}

func (controller *AccountController) CreateSecureAccessPublicKey() *cobra.Command {
	var accountIdUint64 uint64
	var keyNameStr, keyContentStr string

	cmd := &cobra.Command{
		Use:   "create-public-key",
		Short: "CreateSecureAccessPublicKey",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint64,
				"content":   keyContentStr,
			}

			if keyNameStr != "" {
				requestBody["name"] = keyNameStr
			}

			cliHelper.LiaisonResponseWrapper(
				controller.accountLiaison.CreateSecureAccessPublicKey(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint64, "account-id", "u", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(
		&keyNameStr, "public-key-name", "n", "", "SecureAccessPublicKeyName",
	)
	cmd.Flags().StringVarP(
		&keyContentStr, "public-key-content", "c", "", "SecureAccessPublicKeyContent",
	)
	cmd.MarkFlagRequired("public-key-content")
	return cmd
}

func (controller *AccountController) DeleteSecureAccessPublicKey() *cobra.Command {
	var accountIdUint64 uint64
	var keyIdUint16 uint16

	cmd := &cobra.Command{
		Use:   "delete-key",
		Short: "DeleteSecureAccessPublicKey",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint64,
				"id":        keyIdUint16,
			}

			cliHelper.LiaisonResponseWrapper(
				controller.accountLiaison.DeleteSecureAccessPublicKey(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(
		&accountIdUint64, "account-id", "u", 0, "AccountId",
	)
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().Uint16VarP(
		&keyIdUint16, "public-key-id", "i", 0, "SecureAccessPublicKeyId",
	)
	cmd.MarkFlagRequired("public-key-id")
	return cmd
}
