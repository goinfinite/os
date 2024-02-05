package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	accountInfra "github.com/speedianet/os/src/infra/account"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetAccountsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetAccounts",
		Run: func(cmd *cobra.Command, args []string) {
			accQueryRepo := accountInfra.AccQueryRepo{}
			accsList, err := useCase.GetAccounts(accQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, accsList)
		},
	}

	return cmd
}

func CreateAccountController() *cobra.Command {
	var usernameStr string
	var passwordStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "CreateNewAccount",
		Run: func(cmd *cobra.Command, args []string) {
			username := valueObject.NewUsernamePanic(usernameStr)
			password := valueObject.NewPasswordPanic(passwordStr)

			createAccountDto := dto.NewCreateAccount(username, password)

			accQueryRepo := accountInfra.AccQueryRepo{}
			accCmdRepo := accountInfra.AccCmdRepo{}

			err := useCase.CreateAccount(
				accQueryRepo,
				accCmdRepo,
				createAccountDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "AccountCreated")
		},
	}

	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	return cmd
}

func DeleteAccountController() *cobra.Command {
	var accountIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteAccount",
		Run: func(cmd *cobra.Command, args []string) {
			accountId := valueObject.NewAccountIdFromStringPanic(accountIdStr)

			accQueryRepo := accountInfra.AccQueryRepo{}
			accCmdRepo := accountInfra.AccCmdRepo{}

			err := useCase.DeleteAccount(
				accQueryRepo,
				accCmdRepo,
				accountId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "AccountDeleted")
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "account-id", "u", "", "AccountId")
	cmd.MarkFlagRequired("account-id")
	return cmd
}

func UpdateAccountController() *cobra.Command {
	var accountIdStr string
	var passwordStr string
	shouldUpdateApiKeyBool := false

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateAccount (pass or apiKey)",
		Run: func(cmd *cobra.Command, args []string) {
			accountId := valueObject.NewAccountIdFromStringPanic(accountIdStr)

			var passPtr *valueObject.Password
			if passwordStr != "" {
				password := valueObject.NewPasswordPanic(passwordStr)
				passPtr = &password
			}

			var shouldUpdateApiKeyPtr *bool
			if shouldUpdateApiKeyBool {
				shouldUpdateApiKeyPtr = &shouldUpdateApiKeyBool
			}

			updateAccountDto := dto.NewUpdateAccount(
				accountId,
				passPtr,
				shouldUpdateApiKeyPtr,
			)

			accQueryRepo := accountInfra.AccQueryRepo{}
			accCmdRepo := accountInfra.AccCmdRepo{}

			if updateAccountDto.Password != nil {
				useCase.UpdateAccountPassword(
					accQueryRepo,
					accCmdRepo,
					updateAccountDto,
				)
			}

			if shouldUpdateApiKeyBool {
				newKey, err := useCase.UpdateAccountApiKey(
					accQueryRepo,
					accCmdRepo,
					updateAccountDto,
				)
				if err != nil {
					cliHelper.ResponseWrapper(false, err.Error())
				}

				cliHelper.ResponseWrapper(true, newKey)
			}
		},
	}

	cmd.Flags().StringVarP(&accountIdStr, "account-id", "u", "", "AccountId")
	cmd.MarkFlagRequired("account-id")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.Flags().BoolVarP(
		&shouldUpdateApiKeyBool,
		"update-api-key",
		"k",
		false,
		"ShouldUpdateApiKey",
	)
	return cmd
}
