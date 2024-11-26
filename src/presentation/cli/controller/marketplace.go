package cliController

import (
	"errors"
	"strconv"
	"strings"

	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	cliHelper "github.com/goinfinite/os/src/presentation/cli/helper"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/spf13/cobra"
)

type MarketplaceController struct {
	persistentDbSvc    *internalDbInfra.PersistentDatabaseService
	marketplaceService *service.MarketplaceService
}

func NewMarketplaceController(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MarketplaceController {
	return &MarketplaceController{
		persistentDbSvc:    persistentDbSvc,
		marketplaceService: service.NewMarketplaceService(persistentDbSvc, trailDbSvc),
	}
}

func (controller *MarketplaceController) ReadCatalog() *cobra.Command {
	var catalogItemIdUint uint64
	var catalogItemSlugStr, catalogItemNameStr, catalogItemTypeStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "list-catalog",
		Short: "ReadCatalogItems",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if catalogItemIdUint != 0 {
				requestBody["id"] = catalogItemIdUint
			}

			if catalogItemSlugStr != "" {
				requestBody["slug"] = catalogItemSlugStr
			}

			if catalogItemNameStr != "" {
				requestBody["name"] = catalogItemNameStr
			}

			if catalogItemTypeStr != "" {
				requestBody["type"] = catalogItemTypeStr
			}

			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}

			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}

			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}

			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}

			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.marketplaceService.ReadCatalog(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(
		&catalogItemIdUint, "catalog-item-id", "i", 0, "CatalogItemId",
	)
	cmd.Flags().StringVarP(
		&catalogItemSlugStr, "catalog-item-slug", "s", "", "CatalogItemSlug",
	)
	cmd.Flags().StringVarP(
		&catalogItemNameStr, "catalog-item-name", "n", "", "CatalogItemName",
	)
	cmd.Flags().StringVarP(
		&catalogItemTypeStr, "catalog-item-type", "t", "", "CatalogItemType",
	)
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "p", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "m", 0,
		"ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "r", "",
		"SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

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

		dataField, err := valueObject.NewMarketplaceInstallableItemDataField(
			fieldName, fieldValue,
		)
		if err != nil {
			return dataFields, errors.New(errPrefix + "InvalidDataField")
		}

		dataFields = append(dataFields, dataField)
	}

	return dataFields, nil
}

func (controller *MarketplaceController) InstallCatalogItem() *cobra.Command {
	var hostnameStr string
	var catalogIdInt int
	var slugStr, urlPathStr string
	var dataFieldsStr []string

	cmd := &cobra.Command{
		Use:   "install",
		Short: "InstallCatalogItem",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if hostnameStr != "" {
				requestBody["hostname"] = hostnameStr
			}

			if catalogIdInt != 0 {
				requestBody["id"] = catalogIdInt
			}

			if slugStr != "" {
				requestBody["slug"] = slugStr
			}

			if urlPathStr != "" {
				requestBody["urlPath"] = urlPathStr
			}

			dataFields, err := controller.parseDataFields(dataFieldsStr)
			if err != nil {
				cliHelper.ResponseWrapper(false, err.Error())
			}
			requestBody["dataFields"] = dataFields

			cliHelper.ServiceResponseWrapper(
				controller.marketplaceService.InstallCatalogItem(requestBody, false),
			)
		},
	}

	cmd.Flags().StringVarP(&hostnameStr, "hostname", "n", "", "VirtualHostName")
	cmd.Flags().IntVarP(&catalogIdInt, "id", "i", 0, "CatalogItemId")
	cmd.Flags().StringVarP(&slugStr, "slug", "s", "", "CatalogItemSlug")
	cmd.Flags().StringVarP(&urlPathStr, "url-path", "d", "", "UrlPath")
	cmd.Flags().StringSliceVarP(
		&dataFieldsStr, "data-fields", "f", []string{},
		"InstallationDataFields (key:value)",
	)
	return cmd
}

func (controller *MarketplaceController) ReadInstalledItems() *cobra.Command {
	var installedItemIdUint uint64
	var installedItemHostnameStr, installedItemTypeStr, installedItemUuidStr string
	var paginationPageNumberUint32 uint32
	var paginationItemsPerPageUint16 uint16
	var paginationSortByStr, paginationSortDirectionStr, paginationLastSeenIdStr string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "ReadInstalledItems",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{}

			if installedItemIdUint != 0 {
				requestBody["id"] = installedItemIdUint
			}

			if installedItemHostnameStr != "" {
				requestBody["hostname"] = installedItemHostnameStr
			}

			if installedItemTypeStr != "" {
				requestBody["type"] = installedItemTypeStr
			}

			if installedItemUuidStr != "" {
				requestBody["installId"] = installedItemUuidStr
			}

			if paginationPageNumberUint32 != 0 {
				requestBody["pageNumber"] = paginationPageNumberUint32
			}

			if paginationItemsPerPageUint16 != 0 {
				requestBody["itemsPerPage"] = paginationItemsPerPageUint16
			}

			if paginationSortByStr != "" {
				requestBody["sortBy"] = paginationSortByStr
			}

			if paginationSortDirectionStr != "" {
				requestBody["sortDirection"] = paginationSortDirectionStr
			}

			if paginationLastSeenIdStr != "" {
				requestBody["lastSeenId"] = paginationLastSeenIdStr
			}

			cliHelper.ServiceResponseWrapper(
				controller.marketplaceService.ReadInstalledItems(requestBody),
			)
		},
	}

	cmd.Flags().Uint64VarP(
		&installedItemIdUint, "installed-item-id", "i", 0, "InstalledItemId",
	)
	cmd.Flags().StringVarP(
		&installedItemHostnameStr, "installed-item-hostname", "n", "",
		"InstalledItemHostname",
	)
	cmd.Flags().StringVarP(
		&installedItemTypeStr, "installed-item-type", "t", "", "InstalledItemType",
	)
	cmd.Flags().StringVarP(
		&installedItemUuidStr, "installed-item-uuid", "u", "", "InstalledItemUuidStr",
	)
	cmd.Flags().Uint32VarP(
		&paginationPageNumberUint32, "page-number", "p", 0, "PageNumber (Pagination)",
	)
	cmd.Flags().Uint16VarP(
		&paginationItemsPerPageUint16, "items-per-page", "m", 0,
		"ItemsPerPage (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortByStr, "sort-by", "y", "", "SortBy (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationSortDirectionStr, "sort-direction", "r", "",
		"SortDirection (Pagination)",
	)
	cmd.Flags().StringVarP(
		&paginationLastSeenIdStr, "last-seen-id", "l", "", "LastSeenId (Pagination)",
	)

	return cmd
}

func (controller *MarketplaceController) DeleteInstalledItem() *cobra.Command {
	var installedIdInt int
	var shouldUninstallServicesStr string

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "DeleteInstalledItem",
		Run: func(cmd *cobra.Command, args []string) {
			requestBody := map[string]interface{}{
				"installedId":             installedIdInt,
				"shouldUninstallServices": shouldUninstallServicesStr,
			}

			cliHelper.ServiceResponseWrapper(
				controller.marketplaceService.DeleteInstalledItem(requestBody, false),
			)
		},
	}

	cmd.Flags().IntVarP(&installedIdInt, "installed-id", "i", 0, "InstalledItemId")
	cmd.MarkFlagRequired("installed-id")
	cmd.Flags().StringVarP(
		&shouldUninstallServicesStr, "should-uninstall-services", "s", "true",
		"ShouldUninstallUnusedServices",
	)
	return cmd
}
