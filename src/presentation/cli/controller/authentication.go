package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	accountInfra "github.com/speedianet/os/src/infra/account"
	authInfra "github.com/speedianet/os/src/infra/auth"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type AuthenticationController struct {
}

func (controller *AuthenticationController) Login() *cobra.Command {
	var usernameStr string
	var passwordStr string
	var ipAddressStr string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login",
		Run: func(cmd *cobra.Command, args []string) {
			username := valueObject.NewUsernamePanic(usernameStr)
			password := valueObject.NewPasswordPanic(passwordStr)
			ipAddress := valueObject.NewIpAddressPanic(ipAddressStr)

			loginDto := dto.NewLogin(username, password, ipAddress)

			authQueryRepo := authInfra.AuthQueryRepo{}
			authCmdRepo := authInfra.AuthCmdRepo{}
			accQueryRepo := accountInfra.AccQueryRepo{}

			accessToken, err := useCase.GetSessionToken(
				authQueryRepo,
				authCmdRepo,
				accQueryRepo,
				loginDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, accessToken)
		},
	}

	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringVarP(&ipAddressStr, "ip-address", "i", "", "IpAddress")
	cmd.MarkFlagRequired("ip-address")
	return cmd
}
