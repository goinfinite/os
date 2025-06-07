package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	"github.com/spf13/cobra"
)

type AuthenticationController struct {
	authenticationLiaison *liaison.AuthenticationLiaison
}

func NewAuthenticationController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthenticationController {
	return &AuthenticationController{
		authenticationLiaison: liaison.NewAuthenticationLiaison(
			persistentDbSvc, trailDbSvc,
		),
	}
}

func (controller *AuthenticationController) Login() *cobra.Command {
	var usernameStr, passwordStr, ipAddressStr string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"username":          usernameStr,
				"password":          passwordStr,
				"operatorIpAddress": ipAddressStr,
			}

			cliHelper.LiaisonResponseWrapper(
				controller.authenticationLiaison.Login(requestBody),
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
