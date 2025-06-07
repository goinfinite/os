package liaison

import (
	"errors"
	"strconv"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	marketplaceInfra "github.com/goinfinite/os/src/infra/marketplace"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	liaisonHelper "github.com/goinfinite/os/src/presentation/liaison/helper"
)

type MarketplaceLiaison struct {
	marketplaceQueryRepo  *marketplaceInfra.MarketplaceQueryRepo
	marketplaceCmdRepo    *marketplaceInfra.MarketplaceCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceLiaison(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MarketplaceLiaison {
	return &MarketplaceLiaison{
		marketplaceQueryRepo:  marketplaceInfra.NewMarketplaceQueryRepo(persistentDbSvc),
		marketplaceCmdRepo:    marketplaceInfra.NewMarketplaceCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
		persistentDbSvc:       persistentDbSvc,
	}
}

func (liaison *MarketplaceLiaison) ReadCatalog(
	untrustedInput map[string]any,
) LiaisonOutput {
	var idPtr *valueObject.MarketplaceItemId
	if untrustedInput["id"] != nil {
		id, err := valueObject.NewMarketplaceItemId(untrustedInput["id"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		idPtr = &id
	}

	var slugPtr *valueObject.MarketplaceItemSlug
	if untrustedInput["slug"] != nil {
		slug, err := valueObject.NewMarketplaceItemSlug(untrustedInput["slug"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		slugPtr = &slug
	}

	var namePtr *valueObject.MarketplaceItemName
	if untrustedInput["name"] != nil {
		name, err := valueObject.NewMarketplaceItemName(untrustedInput["name"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		namePtr = &name
	}

	var typePtr *valueObject.MarketplaceItemType
	if untrustedInput["type"] != nil {
		itemType, err := valueObject.NewMarketplaceItemType(untrustedInput["type"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		typePtr = &itemType
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if untrustedInput["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(untrustedInput["pageNumber"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if untrustedInput["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(untrustedInput["itemsPerPage"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if untrustedInput["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(untrustedInput["sortBy"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if untrustedInput["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			untrustedInput["sortDirection"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if untrustedInput["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(untrustedInput["lastSeenId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadMarketplaceCatalogItemsRequest{
		Pagination:                 paginationDto,
		MarketplaceCatalogItemId:   idPtr,
		MarketplaceCatalogItemSlug: slugPtr,
		MarketplaceCatalogItemName: namePtr,
		MarketplaceCatalogItemType: typePtr,
	}

	itemsList, err := useCase.ReadMarketplaceCatalogItems(liaison.marketplaceQueryRepo, readDto)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, itemsList)
}

func (liaison *MarketplaceLiaison) InstallCatalogItem(
	untrustedInput map[string]any,
	shouldSchedule bool,
) LiaisonOutput {
	hostname, err := infraHelper.ReadPrimaryVirtualHostHostname()
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	if untrustedInput["hostname"] != nil {
		hostname, err = valueObject.NewFqdn(untrustedInput["hostname"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	var idPtr *valueObject.MarketplaceItemId
	if untrustedInput["id"] != nil {
		id, err := valueObject.NewMarketplaceItemId(untrustedInput["id"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		idPtr = &id
	}

	var slugPtr *valueObject.MarketplaceItemSlug
	if untrustedInput["slug"] != nil {
		slug, err := valueObject.NewMarketplaceItemSlug(untrustedInput["slug"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		slugPtr = &slug
	}

	var urlPathPtr *valueObject.UrlPath
	if untrustedInput["urlPath"] != nil {
		urlPath, err := valueObject.NewUrlPath(untrustedInput["urlPath"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
		urlPathPtr = &urlPath
	}

	dataFields := []valueObject.MarketplaceInstallableItemDataField{}
	if _, exists := untrustedInput["dataFields"]; exists {
		var assertOk bool
		dataFields, assertOk = untrustedInput["dataFields"].([]valueObject.MarketplaceInstallableItemDataField)
		if !assertOk {
			return NewLiaisonOutput(UserError, "InvalidDataFields")
		}
	}

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteOsBinary + " mktplace install"
		installParams := []string{
			"--hostname", hostname.String(),
		}

		if idPtr != nil {
			installParams = append(installParams, "--id", idPtr.String())
		}

		if slugPtr != nil {
			installParams = append(installParams, "--slug", slugPtr.String())
		}

		if urlPathPtr != nil {
			installParams = append(installParams, "--url-path", urlPathPtr.String())
		}

		for _, dataField := range dataFields {
			escapedField := infraHelper.ShellEscape{}.Quote(dataField.String())
			installParams = append(installParams, "--data-fields", escapedField)
		}

		cliCmd += " " + strings.Join(installParams, " ")

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(liaison.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("InstallMarketplaceCatalogItem")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("marketplace")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSecs := uint16(1800)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSecs, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewLiaisonOutput(InfraError, err.Error())
		}

		return NewLiaisonOutput(Created, "MarketplaceCatalogItemInstallationScheduled")
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	installDto := dto.NewInstallMarketplaceCatalogItem(
		hostname, idPtr, slugPtr, urlPathPtr, dataFields, operatorAccountId,
		operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(liaison.persistentDbSvc)

	err = useCase.InstallMarketplaceCatalogItem(
		vhostQueryRepo, liaison.marketplaceQueryRepo, liaison.marketplaceCmdRepo,
		liaison.activityRecordCmdRepo, installDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Created, "MarketplaceCatalogItemInstalled")
}

func (liaison *MarketplaceLiaison) ReadInstalledItems(
	untrustedInput map[string]any,
) LiaisonOutput {
	var idPtr *valueObject.MarketplaceItemId
	if untrustedInput["id"] != nil {
		id, err := valueObject.NewMarketplaceItemId(untrustedInput["id"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		idPtr = &id
	}

	var hostnamePtr *valueObject.Fqdn
	if untrustedInput["hostname"] != nil {
		hostname, err := valueObject.NewFqdn(untrustedInput["hostname"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		hostnamePtr = &hostname
	}

	var typePtr *valueObject.MarketplaceItemType
	if untrustedInput["type"] != nil {
		itemType, err := valueObject.NewMarketplaceItemType(untrustedInput["type"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		typePtr = &itemType
	}

	var installationUuidPtr *valueObject.MarketplaceInstalledItemUuid
	if untrustedInput["installationUuid"] != nil {
		installationUuid, err := valueObject.NewMarketplaceInstalledItemUuid(
			untrustedInput["installationUuid"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		installationUuidPtr = &installationUuid
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if untrustedInput["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(untrustedInput["pageNumber"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if untrustedInput["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(untrustedInput["itemsPerPage"])
		if err != nil {
			return NewLiaisonOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if untrustedInput["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(untrustedInput["sortBy"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if untrustedInput["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			untrustedInput["sortDirection"],
		)
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if untrustedInput["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(untrustedInput["lastSeenId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err)
		}
		paginationDto.LastSeenId = &lastSeenId
	}

	readDto := dto.ReadMarketplaceInstalledItemsRequest{
		Pagination:                       paginationDto,
		MarketplaceInstalledItemId:       idPtr,
		MarketplaceInstalledItemHostname: hostnamePtr,
		MarketplaceInstalledItemType:     typePtr,
		MarketplaceInstalledItemUuid:     installationUuidPtr,
	}

	itemsList, err := useCase.ReadMarketplaceInstalledItems(
		liaison.marketplaceQueryRepo, readDto,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, itemsList)
}

func (liaison *MarketplaceLiaison) DeleteInstalledItem(
	untrustedInput map[string]any,
	shouldSchedule bool,
) LiaisonOutput {
	requiredParams := []string{"installedId"}

	err := liaisonHelper.RequiredParamsInspector(untrustedInput, requiredParams)
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	installedId, err := valueObject.NewMarketplaceItemId(untrustedInput["installedId"])
	if err != nil {
		return NewLiaisonOutput(UserError, err.Error())
	}

	shouldUninstallServices := true
	if untrustedInput["shouldUninstallServices"] != nil {
		shouldUninstallServices, err = voHelper.InterfaceToBool(
			untrustedInput["shouldUninstallServices"],
		)
		if err != nil {
			shouldUninstallServices = false
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(untrustedInput["operatorAccountId"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(untrustedInput["operatorIpAddress"])
		if err != nil {
			return NewLiaisonOutput(UserError, err.Error())
		}
	}

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteOsBinary + " mktplace delete"
		installParams := []string{
			"--installed-id", installedId.String(),
			"--should-uninstall-services", strconv.FormatBool(shouldUninstallServices),
		}

		cliCmd += " " + strings.Join(installParams, " ")

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(
			liaison.persistentDbSvc,
		)
		taskName, _ := valueObject.NewScheduledTaskName("DeleteMarketplaceCatalogItem")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("marketplace")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSeconds := uint16(600)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSeconds, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewLiaisonOutput(InfraError, err.Error())
		}

		return NewLiaisonOutput(Created, "MarketplaceCatalogItemDeletionScheduled")
	}

	deleteMarketplaceInstalledItem := dto.NewDeleteMarketplaceInstalledItem(
		installedId, shouldUninstallServices, operatorAccountId, operatorIpAddress,
	)

	mappingQueryRepo := vhostInfra.NewMappingQueryRepo(liaison.persistentDbSvc)
	mappingCmdRepo := vhostInfra.NewMappingCmdRepo(liaison.persistentDbSvc)
	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(liaison.persistentDbSvc)
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(liaison.persistentDbSvc)

	err = useCase.DeleteMarketplaceInstalledItem(
		liaison.marketplaceQueryRepo, liaison.marketplaceCmdRepo,
		mappingQueryRepo, mappingCmdRepo, servicesQueryRepo, servicesCmdRepo,
		liaison.activityRecordCmdRepo, deleteMarketplaceInstalledItem,
	)
	if err != nil {
		return NewLiaisonOutput(InfraError, err.Error())
	}

	return NewLiaisonOutput(Success, "MarketplaceInstalledItemDeleted")
}
