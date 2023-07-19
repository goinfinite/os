package cliController

import (
	"fmt"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	"github.com/spf13/cobra"
)

func AddUserController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewUser",
		Run: func(cmd *cobra.Command, args []string) {
			username := valueObject.NewUsernamePanic(
				cmd.Flags().Lookup("username").Value.String(),
			)
			password := valueObject.NewPasswordPanic(
				cmd.Flags().Lookup("password").Value.String(),
			)

			addUserDto := dto.NewAddUser(username, password)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err := useCase.AddUser(
				accQueryRepo,
				accCmdRepo,
				addUserDto,
			)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.Flags().StringP("username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringP("password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	return cmd
}

func DeleteUserController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteUser",
		Run: func(cmd *cobra.Command, args []string) {
			userId := valueObject.NewUserIdFromStringPanic(
				cmd.Flags().Lookup("user-id").Value.String(),
			)

			accQueryRepo := infra.AccQueryRepo{}
			accCmdRepo := infra.AccCmdRepo{}

			err := useCase.DeleteUser(
				accQueryRepo,
				accCmdRepo,
				userId,
			)
			if err != nil {
				fmt.Println(err)
			}
		},
	}

	cmd.Flags().StringP("user-id", "u", "", "UserId")
	cmd.MarkFlagRequired("user-id")
	return cmd
}

func UpdateUserController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "UpdateUser (pass or apiKey)",
		Run: func(cmd *cobra.Command, args []string) {
			userId := valueObject.NewUserIdFromStringPanic(
				cmd.Flags().Lookup("user-id").Value.String(),
			)

			var passPtr *valueObject.Password
			if cmd.Flags().Lookup("password").Value.String() != "" {
				password := valueObject.NewPasswordPanic(
					cmd.Flags().Lookup("password").Value.String(),
				)
				passPtr = &password
			}

			var shouldUpdateApiKeyPtr *bool
			if cmd.Flags().Lookup("update-api-key") != nil {
				shouldUpdateApiKey := cmd.Flags().
					Lookup("update-api-key").
					Value.String() == "true"
				shouldUpdateApiKeyPtr = &shouldUpdateApiKey
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

			if updateUserDto.ShouldUpdateApiKey != nil && *updateUserDto.ShouldUpdateApiKey {
				newKey, err := useCase.UpdateUserApiKey(
					accQueryRepo,
					accCmdRepo,
					updateUserDto,
				)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(newKey)
			}
		},
	}

	cmd.Flags().StringP("user-id", "u", "", "UserId")
	cmd.MarkFlagRequired("user-id")
	cmd.Flags().StringP("password", "p", "", "Password")
	cmd.Flags().BoolP("update-api-key", "k", false, "ShouldUpdateApiKey")
	return cmd
}
