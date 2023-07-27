package cliController

import (
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetUsersController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "GetUsers",
		Run: func(cmd *cobra.Command, args []string) {
			accQueryRepo := infra.AccQueryRepo{}
			usersList, err := useCase.GetUsers(accQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
			}

			cliHelper.ResponseWrapper(true, usersList)
		},
	}

	return cmd
}

func AddUserController() *cobra.Command {
	var usernameStr string
	var passwordStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewUser",
		Run: func(cmd *cobra.Command, args []string) {
			username := valueObject.NewUsernamePanic(usernameStr)
			password := valueObject.NewPasswordPanic(passwordStr)

			addUserDto := dto.NewAddUser(username, password)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err := useCase.AddUser(
				accQueryRepo,
				accCmdRepo,
				addUserDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
			}

			cliHelper.ResponseWrapper(true, "UserAdded")
		},
	}

	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	return cmd
}

func DeleteUserController() *cobra.Command {
	var userIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteUser",
		Run: func(cmd *cobra.Command, args []string) {
			userId := valueObject.NewUserIdFromStringPanic(userIdStr)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err := useCase.DeleteUser(
				accQueryRepo,
				accCmdRepo,
				userId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err)
			}

			cliHelper.ResponseWrapper(true, "UserDeleted")
		},
	}

	cmd.Flags().StringVarP(&userIdStr, "user-id", "u", "", "UserId")
	cmd.MarkFlagRequired("user-id")
	return cmd
}

func UpdateUserController() *cobra.Command {
	var userIdStr string
	var passwordStr string
	shouldUpdateApiKeyBool := false

	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateUser (pass or apiKey)",
		Run: func(cmd *cobra.Command, args []string) {
			userId := valueObject.NewUserIdFromStringPanic(userIdStr)

			var passPtr *valueObject.Password
			if passwordStr != "" {
				password := valueObject.NewPasswordPanic(passwordStr)
				passPtr = &password
			}

			var shouldUpdateApiKeyPtr *bool
			if shouldUpdateApiKeyBool {
				shouldUpdateApiKeyPtr = &shouldUpdateApiKeyBool
			}

			updateUserDto := dto.NewUpdateUser(
				userId,
				passPtr,
				shouldUpdateApiKeyPtr,
			)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			if updateUserDto.Password != nil {
				useCase.UpdateUserPassword(
					accQueryRepo,
					accCmdRepo,
					updateUserDto,
				)
			}

			if shouldUpdateApiKeyBool {
				newKey, err := useCase.UpdateUserApiKey(
					accQueryRepo,
					accCmdRepo,
					updateUserDto,
				)
				if err != nil {
					cliHelper.ResponseWrapper(false, err)
				}

				cliHelper.ResponseWrapper(true, newKey)
			}
		},
	}

	cmd.Flags().StringVarP(&userIdStr, "user-id", "u", "", "UserId")
	cmd.MarkFlagRequired("user-id")
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
