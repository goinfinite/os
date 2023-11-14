package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetSslPairsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetSslPairs",
		Run: func(cmd *cobra.Command, args []string) {
			sslQueryRepo := infra.SslQueryRepo{}
			sslPairsList, err := useCase.GetSslPairs(sslQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, sslPairsList)
		},
	}

	return cmd
}

func AddSslPairController() *cobra.Command {
	var hostnameStr string
	var certificateFilePathStr string
	var keyFilePathStr string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewSslPair",
		Run: func(cmd *cobra.Command, args []string) {
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

			addSslDto := dto.NewAddSslPair(
				valueObject.NewFqdnPanic(hostnameStr),
				sslCertificate,
				sslPrivateKey,
			)

			sslCmdRepo := infra.SslCmdRepo{}

			err = useCase.AddSslPair(
				sslCmdRepo,
				addSslDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "SslPairAdded")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "t", "", "Hostname")
	cmd.MarkFlagRequired("hostname")
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

			cronQueryRepo := infra.SslQueryRepo{}
			cronCmdRepo := infra.SslCmdRepo{}

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
