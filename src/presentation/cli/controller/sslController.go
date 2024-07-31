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
	transientDbSvc  *internalDbInfra.TransientDatabaseService
}

func NewSslController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *SslController {
	return &SslController{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
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
				virtualHost, err := valueObject.NewFqdn(vhost)
				if err != nil {
					continue
				}
				virtualHosts = append(virtualHosts, virtualHost)
			}

			certificateFilePath, err := valueObject.NewUnixFilePath(certificateFilePathStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, "InvalidCertificateFilePath")
			}

			certificateContentStr, err := infraHelper.GetFileContent(
				certificateFilePath.String(),
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, "OpenSslCertificateFileError")
			}
			sslCertificateContent, err := valueObject.NewSslCertificateContent(certificateContentStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			privateKeyContentStr, err := infraHelper.GetFileContent(keyFilePathStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, "OpenPrivateKeyFileError")
			}

			sslCertificate := entity.NewSslCertificatePanic(sslCertificateContent)
			sslPrivateKey, err := valueObject.NewSslPrivateKey(privateKeyContentStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			createSslDto := dto.NewCreateSslPair(
				virtualHosts,
				sslCertificate,
				sslPrivateKey,
			)

			sslCmdRepo := sslInfra.NewSslCmdRepo(
				controller.persistentDbSvc, controller.transientDbSvc,
			)
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
			cronCmdRepo := sslInfra.NewSslCmdRepo(
				controller.persistentDbSvc, controller.transientDbSvc,
			)

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
				virtualHost, err := valueObject.NewFqdn(vhost)
				if err != nil {
					continue
				}
				virtualHosts = append(virtualHosts, virtualHost)
			}

			dto := dto.NewDeleteSslPairVhosts(sslPairId, virtualHosts)

			sslQueryRepo := sslInfra.SslQueryRepo{}
			sslCmdRepo := sslInfra.NewSslCmdRepo(
				controller.persistentDbSvc, controller.transientDbSvc,
			)

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
