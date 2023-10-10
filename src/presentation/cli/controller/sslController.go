package cliController

import (
	"strings"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/useCase"
	"github.com/speedianet/sam/src/domain/valueObject"
	"github.com/speedianet/sam/src/infra"
	cliHelper "github.com/speedianet/sam/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

func GetSslsController() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetSsls",
		Run: func(cmd *cobra.Command, args []string) {
			sslQueryRepo := infra.NewSslQueryRepo()
			sslPairsList, err := useCase.GetSsls(sslQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, sslPairsList)
		},
	}

	return cmd
}

func AddSslControler() *cobra.Command {
	var hostnameStr string
	var certificateStr string
	var keyStr string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "AddNewSSl",
		Run: func(cmd *cobra.Command, args []string) {
			parsedCertificateStr := strings.Replace(certificateStr, "\\n", "\n", -1)
			parsedKeyStr := strings.Replace(keyStr, "\\n", "\n", -1)

			sslCertificate, err := entity.NewSslCertificate(parsedCertificateStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			sslPrivateKey, err := entity.NewSslPrivateKey(parsedKeyStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			addSslDto := dto.NewAddSsl(
				valueObject.NewFqdnPanic(hostnameStr),
				sslCertificate,
				sslPrivateKey,
			)

			sslCmdRepo := infra.SslCmdRepo{}

			err = useCase.AddSsl(
				sslCmdRepo,
				addSslDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "SslAdded")
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "virtualHostname", "v", "", "Virtual Hostname")
	cmd.MarkFlagRequired("virtualHostname")
	cmd.Flags().StringVarP(&certificateStr, "certificate", "c", "", "Certificate")
	cmd.MarkFlagRequired("certificate")
	cmd.Flags().StringVarP(&keyStr, "key", "k", "", "Key")
	cmd.MarkFlagRequired("key")
	return cmd
}

func DeleteSslController() *cobra.Command {
	var sslSerialNumberStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteSsl",
		Run: func(cmd *cobra.Command, args []string) {
			sslSerialNumber := valueObject.NewSslSerialNumberPanic(sslSerialNumberStr)

			cronQueryRepo := infra.NewSslQueryRepo()
			cronCmdRepo := infra.SslCmdRepo{}

			err := useCase.DeleteSsl(
				cronQueryRepo,
				cronCmdRepo,
				sslSerialNumber,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "SslDeleted")
		},
	}

	cmd.Flags().StringVarP(&sslSerialNumberStr, "id", "i", "", "SslSerialNumber")
	cmd.MarkFlagRequired("id")
	return cmd
}
