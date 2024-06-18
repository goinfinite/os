package cliController

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	sslInfra "github.com/speedianet/os/src/infra/ssl"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type SslController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewSslController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *SslController {
	return &SslController{
		persistentDbSvc: persistentDbSvc,
	}
}

func (controller *SslController) Read() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "GetSslPairs",
		Run: func(cmd *cobra.Command, args []string) {
			sslQueryRepo := sslInfra.SslQueryRepo{}
			sslPairsList, err := useCase.ReadSslPairs(sslQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, sslPairsList)
		},
	}

	return cmd
}

func (controller *SslController) Create() *cobra.Command {
	var virtualHostsSlice []string
	var certificateFilePathStr string
	var keyFilePathStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "CreateSslPair",
		Run: func(cmd *cobra.Command, args []string) {
			var virtualHosts []valueObject.Fqdn
			for _, vhost := range virtualHostsSlice {
				virtualHosts = append(virtualHosts, valueObject.NewFqdnPanic(vhost))
			}

			certificateFilePath := valueObject.NewUnixFilePathPanic(certificateFilePathStr)
			certificateContentStr, err := infraHelper.GetFileContent(
				certificateFilePath.String(),
			)
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

			createSslDto := dto.NewCreateSslPair(
				virtualHosts,
				sslCertificate,
				sslPrivateKey,
			)

			sslCmdRepo := sslInfra.NewSslCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)

			err = useCase.CreateSslPair(
				sslCmdRepo,
				vhostQueryRepo,
				createSslDto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "SslPairCreated")
		},
	}

	cmd.Flags().StringSliceVarP(&virtualHostsSlice, "virtualHosts", "v", []string{}, "VirtualHosts")
	cmd.MarkFlagRequired("virtualHosts")
	cmd.Flags().StringVarP(&certificateFilePathStr, "certFilePath", "c", "", "CertificateFilePath")
	cmd.MarkFlagRequired("certFilePath")
	cmd.Flags().StringVarP(&keyFilePathStr, "keyFilePath", "k", "", "KeyFilePath")
	cmd.MarkFlagRequired("keyFilePath")
	return cmd
}

func (controller *SslController) Delete() *cobra.Command {
	var sslPairIdStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteSslPair",
		Run: func(cmd *cobra.Command, args []string) {
			sslId := valueObject.NewSslIdPanic(sslPairIdStr)

			cronQueryRepo := sslInfra.SslQueryRepo{}
			cronCmdRepo := sslInfra.NewSslCmdRepo(controller.persistentDbSvc)

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

	cmd.Flags().StringVarP(&sslPairIdStr, "id", "i", "", "SslPairId")
	cmd.MarkFlagRequired("sslPairId")
	return cmd
}

func (controller *SslController) DeleteVhosts() *cobra.Command {
	var sslPairIdStr string
	var virtualHostsSlice []string

	cmd := &cobra.Command{
		Use:   "remove-vhosts",
		Short: "RemoveSslPairVhosts",
		Run: func(cmd *cobra.Command, args []string) {
			sslPairId := valueObject.NewSslIdPanic(sslPairIdStr)

			var virtualHosts []valueObject.Fqdn
			for _, vhost := range virtualHostsSlice {
				virtualHosts = append(virtualHosts, valueObject.NewFqdnPanic(vhost))
			}

			dto := dto.NewDeleteSslPairVhosts(sslPairId, virtualHosts)

			sslQueryRepo := sslInfra.SslQueryRepo{}
			sslCmdRepo := sslInfra.NewSslCmdRepo(controller.persistentDbSvc)

			err := useCase.DeleteSslPairVhosts(
				sslQueryRepo,
				sslCmdRepo,
				dto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "SslPairVhostsRemoved")
		},
	}

	cmd.Flags().StringVarP(&sslPairIdStr, "id", "i", "", "SslPairId")
	cmd.MarkFlagRequired("sslPairId")
	cmd.Flags().StringSliceVarP(&virtualHostsSlice, "virtualHosts", "v", []string{}, "VirtualHosts")
	cmd.MarkFlagRequired("virtualHosts")
	return cmd
}
