package cliController

import (
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	marketplaceInfra "github.com/speedianet/os/src/infra/marketplace"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	cliHelper "github.com/speedianet/os/src/presentation/cli/helper"
	"github.com/spf13/cobra"
)

type MarketplaceController struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceController {
	return &MarketplaceController{
		persistentDbSvc: persistentDbSvc,
	}
}

func (controller *MarketplaceController) ReadInstalled() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "ReadInstalledItems",
		Run: func(cmd *cobra.Command, args []string) {
			marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(
				controller.persistentDbSvc,
			)

			installedItems, err := useCase.ReadMarketplaceInstalledItems(marketplaceQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, installedItems)
		},
	}
	return cmd
}

func parseDataFields(
	dataFieldsStr []string,
) []valueObject.MarketplaceInstallableItemDataField {
	dataFields := []valueObject.MarketplaceInstallableItemDataField{}

	for _, dataFieldStr := range dataFieldsStr {
		dataFieldsParts := strings.Split(dataFieldStr, ":")
		if len(dataFieldsParts) < 2 {
			panic("InvalidDataFields")
		}

		dataField := valueObject.NewMarketplaceInstallableItemDataFieldPanic(
			valueObject.NewDataFieldNamePanic(dataFieldsParts[0]),
			valueObject.NewDataFieldValuePanic(dataFieldsParts[1]),
		)
		dataFields = append(dataFields, dataField)
	}

	return dataFields
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
			hostname := valueObject.NewFqdnPanic(hostnameStr)

			var catalogIdPtr *valueObject.MarketplaceItemId
			if catalogIdInt != 0 {
				catalogId := valueObject.NewMarketplaceItemIdPanic(catalogIdInt)
				catalogIdPtr = &catalogId
			}

			var slugPtr *valueObject.MarketplaceItemSlug
			if slugStr != "" {
				slug := valueObject.NewMarketplaceItemSlugPanic(slugStr)
				slugPtr = &slug
			}

			var urlPathPtr *valueObject.UrlPath
			if urlPath != "" {
				urlPath := valueObject.NewUrlPathPanic(urlPath)
				urlPathPtr = &urlPath
			}

			// Format: key:value,key:value
			dataFields := parseDataFields(dataFieldsStr)

			marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)
			marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(controller.persistentDbSvc)
			vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(controller.persistentDbSvc)
			vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(controller.persistentDbSvc)

			dto := dto.NewInstallMarketplaceCatalogItem(
				catalogIdPtr,
				slugPtr,
				hostname,
				urlPathPtr,
				dataFields,
			)
			err := useCase.InstallMarketplaceCatalogItem(
				marketplaceQueryRepo,
				marketplaceCmdRepo,
				vhostQueryRepo,
				vhostCmdRepo,
				dto,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MarketplaceCatalogItemInstalled")
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

func (controller *MarketplaceController) DeleteInstalledItem() *cobra.Command {
	var installedIdInt int
	var shouldUninstallServices bool
	var shouldRemoveFiles bool

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteInstalledItem",
		Run: func(cmd *cobra.Command, args []string) {
			installedId := valueObject.NewMarketplaceItemIdPanic(installedIdInt)

			deleteMarketplaceInstalledItem := dto.NewDeleteMarketplaceInstalledItem(
				installedId, shouldUninstallServices, shouldRemoveFiles,
			)

			marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(controller.persistentDbSvc)
			marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(controller.persistentDbSvc)

			err := useCase.DeleteMarketplaceInstalledItem(
				marketplaceQueryRepo,
				marketplaceCmdRepo,
				deleteMarketplaceInstalledItem,
			)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, "MarketplaceInstalledItemDeleted")
		},
	}

	cmd.Flags().IntVarP(&installedIdInt, "installedId", "i", 0, "InstalledItemId")
	cmd.MarkFlagRequired("installedId")
	cmd.Flags().BoolVarP(
		&shouldUninstallServices, "shouldUninstallServices", "s", true,
		"ShouldUninstallInstalledItemServices",
	)
	cmd.Flags().BoolVarP(
		&shouldRemoveFiles, "shouldRemoveFiles", "f", true,
		"ShouldRemoveInstalledItemFiles",
	)
	return cmd
}

func (controller *MarketplaceController) ReadCatalog() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-catalog",
		Short: "ReadCatalogItems",
		Run: func(cmd *cobra.Command, args []string) {
			marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(
				controller.persistentDbSvc,
			)

			catalogItems, err := useCase.ReadMarketplaceCatalog(marketplaceQueryRepo)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}

			cliHelper.ResponseWrapper(true, catalogItems)
		},
	}
	return cmd
}
