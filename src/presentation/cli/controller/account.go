package cliController

import (
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type AccountController struct {
	accountService service.AccountService
}

func NewAccountController() AccountController {
	return AccountController{
		accountService: service.NewAccountService(),
	}
}

func (controller AccountController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetAccounts",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.accountService.Read())
		},
	}

	return cmd
}

func (controller AccountController) Create() *cobra.Command {
	var usernameStr string
	var passwordStr string

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

func (controller AccountController) Update() *cobra.Command {
	var accountIdStr string
	var usernameStr string
	var passwordStr string
	shouldUpdateApiKeyBool := false

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateAccount (pass or apiKey)",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"shouldUpdateApiKey": shouldUpdateApiKeyBool,
			}

			if accountIdStr != "" {
				requestBody["id"] = accountIdStr
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

	cmd.Flags().StringVarP(&accountIdStr, "account-id", "i", "", "AccountId")
	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.Flags().BoolVarP(
		&shouldUpdateApiKeyBool, "update-api-key", "k", false, "ShouldUpdateApiKey",
	)
	return cmd
}

func (controller AccountController) Delete() *cobra.Command {
	var accountIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteAccount",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": accountIdStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.accountService.Delete(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "account-id", "i", "", "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}
