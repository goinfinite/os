package cliController

import (
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/service"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/spf13/cobra"
)

type AuthenticationController struct {
	authenticationService *service.AuthenticationService
}

func NewAuthenticationController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AuthenticationController {
	return &AuthenticationController{
		authenticationService: service.NewAuthenticationService(
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
			requestInputParsed := cliHelper.RequestInputParser(map[string]any{
				"username":          usernameStr,
				"password":          passwordStr,
				"operatorIpAddress": ipAddressStr,
			})

			cliHelper.ServiceResponseWrapper(
				controller.authenticationService.Login(requestInputParsed),
			)
		},
	}

	cmd.Flags().StringVarP(
		&usernameStr, "username", "u", tkPresentation.UnsetParameterValueStr, "Username",
	)
	cmd.MarkFlagRequired("username")
	cmd.Flags().StringVarP(
		&passwordStr, "password", "p", tkPresentation.UnsetParameterValueStr, "Password",
	)
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringVarP(
		&ipAddressStr, "ip-address", "i", tkPresentation.UnsetParameterValueStr, "IpAddress",
	)
	cmd.MarkFlagRequired("ip-address")
	return cmd
}
