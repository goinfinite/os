package cliController

import (
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthController {
	return &AuthController{
		authService: service.NewAuthService(trailDbSvc),
	}
}

func (controller *AuthController) Login() *cobra.Command {
	var usernameStr, passwordStr, ipAddressStr string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"username":  usernameStr,
				"password":  passwordStr,
				"ipAddress": ipAddressStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.authService.GenerateJwtWithCredentials(requestBody),
			)
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
