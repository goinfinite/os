package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetSslPairsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetSslPairs",
		Run: func(cmd *cobra.Command, args []string) {
			sslQueryRepo := sslInfra.SslQueryRepo{}
			sslPairsList, err := useCase.GetSslPairs(sslQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, sslPairsList)
		},
	}

	return cmd
}

func CreateSslPairController() *cobra.Command {
	var virtualHostsSlice []string
	var certificateFilePathStr string
	var keyFilePathStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewSslPair",
		Run: func(cmd *cobra.Command, args []string) {
			var virtualHosts []valueObject.Fqdn
			for _, vhost := range virtualHostsSlice {
				virtualHosts = append(virtualHosts, valueObject.NewFqdnPanic(vhost))
			}

			certificateContentStr, err := infraHelper.GetFileContent(certificateFilePathStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, "FailedToOpenSslCertificateFile")
			}
			sslCertificateContent := valueObject.NewSslCertificateContentPanic(certificateContentStr)

			privateKeyContentStr, err := infraHelper.GetFileContent(keyFilePathStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, "FailedToOpenPrivateKeyFile")
			}

			sslCertificate := entity.NewSslCertificatePanic(sslCertificateContent)
			sslPrivateKey := valueObject.NewSslPrivateKeyPanic(privateKeyContentStr)

			addSslDto := dto.NewCreateSslPair(
				virtualHosts,
				sslCertificate,
				sslPrivateKey,
			)

			sslCmdRepo := sslInfra.NewSslCmdRepo()
			vhostQueryRepo := vhostInfra.VirtualHostQueryRepo{}

			err = useCase.CreateSslPair(
				sslCmdRepo,
				vhostQueryRepo,
				addSslDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "SslPairCreated")
		},
	}

	cmd.Flags().StringSliceVarP(&virtualHostsSlice, "virtualHosts", "t", []string{}, "VirtualHosts")
	cmd.MarkFlagRequired("virtualHosts")
	cmd.Flags().StringVarP(&certificateFilePathStr, "certFilePath", "c", "", "CertificateFilePath")
	cmd.MarkFlagRequired("certFilePath")
	cmd.Flags().StringVarP(&keyFilePathStr, "keyFilePath", "k", "", "KeyFilePath")
	cmd.MarkFlagRequired("keyFilePath")
	return cmd
}

func DeleteSslPairController() *cobra.Command {
	var sslPairIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteSslPair",
		Run: func(cmd *cobra.Command, args []string) {
			sslId := valueObject.NewSslIdPanic(sslPairIdStr)

			cronQueryRepo := sslInfra.SslQueryRepo{}
			cronCmdRepo := sslInfra.NewSslCmdRepo()

			err := useCase.DeleteSslPair(
				cronQueryRepo,
				cronCmdRepo,
				sslId,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "SslPairDeleted")
		},
	}

	cmd.Flags().StringVarP(&sslPairIdStr, "sslPairId", "s", "", "SslPairId")
	cmd.MarkFlagRequired("sslPairId")
	return cmd
}
