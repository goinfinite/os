package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/liaison"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
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

			tkPresentation.LiaisonCliResponseRenderer(
				controller.authenticationLiaison.Login(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&usernameStr, "username", "u", "", "Username")
	cliHelper.MarkRequiredFlag(cmd, "username")
	cmd.Flags().StringVarP(&passwordStr, "password", "p", "", "Password")
	cliHelper.MarkRequiredFlag(cmd, "password")
	cmd.Flags().StringVarP(&ipAddressStr, "ip-address", "i", "", "IpAddress")
	cliHelper.MarkRequiredFlag(cmd, "ip-address")
	return cmd
}
