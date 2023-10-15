package cliController

import (
	"os"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
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
			certificateBytesOutput, err := os.ReadFile(certificateFilePathStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, "FailedToOpenSslCertificateFile")
			}
			certificateOutputStr := string(certificateBytesOutput)

			privateKeyBytesOutput, err := os.ReadFile(keyFilePathStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, "FailedToOpenPrivateKeyFile")
			}
			privateKeyOutputStr := string(privateKeyBytesOutput)

			sslCertificate := entity.NewSslCertificatePanic(certificateOutputStr)
			sslPrivateKey := valueObject.NewSslPrivateKeyPanic(privateKeyOutputStr)

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
	cmd.Flags().StringVarP(&certificateFilePathStr, "certFilePath", "c", "", "Certificate File Path")
	cmd.MarkFlagRequired("certificateFilePath")
	cmd.Flags().StringVarP(&keyFilePathStr, "keyFilePath", "k", "", "Key File Path")
	cmd.MarkFlagRequired("keyFilePath")
	return cmd
}

func DeleteSslPairController() *cobra.Command {
	var sslSerialNumberStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteSslPair",
		Run: func(cmd *cobra.Command, args []string) {
			sslSerialNumber := valueObject.NewSslSerialNumberPanic(sslSerialNumberStr)

			cronQueryRepo := infra.SslQueryRepo{}
			cronCmdRepo := infra.SslCmdRepo{}

			err := useCase.DeleteSslPair(
				cronQueryRepo,
				cronCmdRepo,
				sslSerialNumber,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "SslPairDeleted")
		},
	}

	cmd.Flags().StringVarP(&sslSerialNumberStr, "serialNumber", "s", "", "SslSerialNumber")
	cmd.MarkFlagRequired("serialNumber")
	return cmd
}
