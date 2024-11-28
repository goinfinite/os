package cliController

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type SslController struct {
	sslService *service.SslService
}

func NewSslController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SslController {
	return &SslController{
		sslService: service.NewSslService(persistentDbSvc, transientDbSvc, trailDbSvc),
	}
}

func (controller *SslController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "ReadSslPairs",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.sslService.Read())
		},
	}

	return cmd
}

func (controller *SslController) Create() *cobra.Command {
	var virtualHostsSlice []string
	var certFilePathStr, keyFilePathStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateSslPair",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"virtualHosts": virtualHostsSlice,
			}

			certFilePath, err := valueObject.NewUnixFilePath(certFilePathStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, "InvalidCertificateFilePath")
			}
			certContentStr, err := infraHelper.GetFileContent(certFilePath.String())
			if err != nil {
				cliHelper.ResponseWrapper(false, "OpenSslCertificateFileError")
			}
			requestBody["certificate"] = certContentStr

			privateKeyFilePath, err := valueObject.NewUnixFilePath(keyFilePathStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, "InvalidSslPrivateKeyFilePath")
			}
			privateKeyContentStr, err := infraHelper.GetFileContent(
				privateKeyFilePath.String(),
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, "OpenSslPrivateKeyFileError")
			}
			requestBody["key"] = privateKeyContentStr

			cliHelper.ServiceResponseWrapper(controller.sslService.Create(requestBody))
		},
	}

	cmd.Flags().StringSliceVarP(
		&virtualHostsSlice, "virtualHosts", "v", []string{}, "VirtualHosts",
	)
	cmd.MarkFlagRequired("virtualHosts")
	cmd.Flags().StringVarP(
		&certFilePathStr, "certFilePath", "c", "", "SslCertificateFilePath",
	)
	cmd.MarkFlagRequired("certFilePath")
	cmd.Flags().StringVarP(&keyFilePathStr, "keyFilePath", "k", "", "SslKeyFilePath")
	cmd.MarkFlagRequired("keyFilePath")
	return cmd
}

func (controller *SslController) Delete() *cobra.Command {
	var sslPairIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteSslPair",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id": sslPairIdStr,
			}

			cliHelper.ServiceResponseWrapper(controller.sslService.Delete(requestBody))
		},
	}

	cmd.Flags().StringVarP(&sslPairIdStr, "pairId", "i", "", "SslPairId")
	cmd.MarkFlagRequired("pairId")
	return cmd
}

func (controller *SslController) DeleteVhosts() *cobra.Command {
	var sslPairIdStr string
	var virtualHostsSlice []string

	cmd := &cobra.Command{
		Use:   "delete-vhosts",
		Short: "RemoveSslPairVhosts",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"id":           sslPairIdStr,
				"virtualHosts": virtualHostsSlice,
			}

			cliHelper.ServiceResponseWrapper(
				controller.sslService.DeleteVhosts(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&sslPairIdStr, "id", "i", "", "SslPairId")
	cmd.MarkFlagRequired("sslPairId")
	cmd.Flags().StringSliceVarP(
		&virtualHostsSlice, "virtualHosts", "v", []string{}, "VirtualHosts",
	)
	cmd.MarkFlagRequired("virtualHosts")
	return cmd
}
