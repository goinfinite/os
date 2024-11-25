package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountController {
	return &AccountController{
		accountService: service.NewAccountService(persistentDbSvc, trailDbSvc),
	}
}

func (controller *AccountController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetAccounts",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.accountService.Read())
		},
	}

	return cmd
}

func (controller *AccountController) Create() *cobra.Command {
	var usernameStr, passwordStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateNewAccount",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"username": usernameStr,
				"password": passwordStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.Create(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	return cmd
}

func (controller *AccountController) Update() *cobra.Command {
	var accountIdUint64 uint64
	var usernameStr, passwordStr, shouldUpdateApiKeyStr string

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
				requestBody["username"] = usernameStr
			}

			if passwordStr != "" {
				requestBody["password"] = passwordStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.Update(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint64, "account-id", "i", 0, "AccountId")
	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
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

			cliHelper.ServiceResponseWrapper(
				controller.accountService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint64, "account-id", "i", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}

func (controller *AccountController) ReadSecureAccessKeys() *cobra.Command {
	var accountIdUint64 uint64

	cmd := &cobra.Command{
		Use:   "get-keys",
		Short: "GetSecureAccessKeys",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint64,
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.ReadSecureAccessKey(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint64, "account-id", "u", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}

func (controller *AccountController) CreateSecureAccessKey() *cobra.Command {
	var accountIdUint64 uint64
	var keyNameStr, keyContentStr string

	cmd := &cobra.Command{
		Use:   "create-key",
		Short: "CreateSecureAccessKey",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint64,
				"content":   keyContentStr,
			}

			if keyNameStr != "" {
				requestBody["name"] = keyNameStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.CreateSecureAccessKey(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(&accountIdUint64, "account-id", "u", 0, "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&keyNameStr, "key-name", "n", "", "SecureAccessKeyName")
	cmd.Flags().StringVarP(
		&keyContentStr, "key-content", "c", "", "SecureAccessKeyContent",
	)
	cmd.MarkFlagRequired("key-content")
	return cmd
}

func (controller *AccountController) DeleteSecureAccessKey() *cobra.Command {
	var accountIdUint64 uint64
	var keyIdUint16 uint16

	cmd := &cobra.Command{
		Use:   "delete-key",
		Short: "DeleteSecureAccessKey",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"accountId": accountIdUint64,
				"id":        keyIdUint16,
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.DeleteSecureAccessKey(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(
		&accountIdUint64, "account-id", "u", 0, "AccountId",
	)
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().Uint16VarP(&keyIdUint16, "key-id", "i", 0, "SecureAccessKeyId")
	cmd.MarkFlagRequired("key-id")
	return cmd
}
