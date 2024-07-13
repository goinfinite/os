package cliController

import (
	"errors"
	"strconv"
	"strings"

	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/speedianet/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type MarketplaceController struct {
	persistentDbSvc    *internalDbInfra.PersistentDatabaseService
	marketplaceService *service.MarketplaceService
}

func NewMarketplaceController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceController {
	return &MarketplaceController{
		persistentDbSvc:    persistentDbSvc,
		marketplaceService: service.NewMarketplaceService(persistentDbSvc),
	}
}

func (controller *MarketplaceController) ReadCatalog() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-catalog",
		Short: "ReadCatalogItems",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(controller.marketplaceService.ReadCatalog())
		},
	}

	return cmd
}

func (controller *MarketplaceController) parseDataFields(
	rawDataFields []string,
) ([]valueObject.MarketplaceInstallableItemDataField, error) {
	dataFields := []valueObject.MarketplaceInstallableItemDataField{}

	for fieldIndex, rawDataField := range rawDataFields {
		errPrefix := "[index " + strconv.Itoa(fieldIndex) + "] "

		dataFieldsParts := strings.Split(rawDataField, ":")
		if len(dataFieldsParts) < 2 {
			return dataFields, errors.New(errPrefix + "InvalidDataFields")
		}

		fieldName, err := valueObject.NewDataFieldName(dataFieldsParts[0])
		if err != nil {
			return dataFields, errors.New(errPrefix + "InvalidDataFieldName")
		}

		fieldValue, err := valueObject.NewDataFieldValue(dataFieldsParts[1])
		if err != nil {
			return dataFields, errors.New(errPrefix + "InvalidDataFieldValue")
		}

		dataField := valueObject.NewMarketplaceInstallableItemDataFieldPanic(
			fieldName, fieldValue,
		)

		dataFields = append(dataFields, dataField)
	}

	return dataFields, nil
}

func (controller *MarketplaceController) InstallCatalogItem() *cobra.Command {
	var hostnameStr string
	var catalogIdInt int
	var slugStr string
	var urlPath string
	var dataFieldsStr []string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "InstallCatalogItem",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"hostname": hostnameStr,
			}

			if catalogIdInt != 0 {
				requestBody["id"] = catalogIdInt
			}

			if slugStr != "" {
				requestBody["slug"] = slugStr
			}

			if urlPath != "" {
				requestBody["urlPath"] = urlPath
			}

			dataFields, err := controller.parseDataFields(dataFieldsStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}
			requestBody["dataFields"] = dataFields

			cliHelper.ServiceResponseWrapper(
				controller.marketplaceService.InstallCatalogItem(requestBody),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostName")
	cmd.MarkFlagRequired("hostname")
	cmd.Flags().IntVarP(&catalogIdInt, "catalogId", "i", 0, "CatalogItemId")
	cmd.Flags().StringVarP(&slugStr, "slug", "s", "", "CatalogItemSlug")
	cmd.Flags().StringVarP(&urlPath, "urlPath", "d", "", "UrlPath")
	cmd.Flags().StringSliceVarP(
		&dataFieldsStr, "dataFields", "f", []string{}, "InstallationDataFields (key:value)",
	)
	return cmd
}

func (controller *MarketplaceController) ReadInstalledItems() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "ReadInstalledItems",
		Run: func(cmd *cobra.Command, args []string) {
			cliHelper.ServiceResponseWrapper(
				controller.marketplaceService.ReadInstalledItems(),
			)
		},
	}

	return cmd
}

func (controller *MarketplaceController) DeleteInstalledItem() *cobra.Command {
	var installedIdInt int
	var shouldUninstallServices bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteInstalledItem",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"installedId": installedIdInt,
			}

			if !shouldUninstallServices {
				requestBody["shouldUninstallServices"] = false
			}

			cliHelper.ServiceResponseWrapper(
				controller.marketplaceService.DeleteInstalledItem(requestBody),
			)
		},
	}

	cmd.Flags().IntVarP(&installedIdInt, "installedId", "i", 0, "InstalledItemId")
	cmd.MarkFlagRequired("installedId")
	cmd.Flags().BoolVarP(
		&shouldUninstallServices, "shouldUninstallServices", "s", true,
		"ShouldUninstallUnusedServices",
	)
	return cmd
}
