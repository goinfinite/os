package service

import (
	"errors"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
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
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type MarketplaceService struct {
	marketplaceQueryRepo  *marketplaceInfra.MarketplaceQueryRepo
	marketplaceCmdRepo    *marketplaceInfra.MarketplaceCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MarketplaceService {
	return &MarketplaceService{
		marketplaceQueryRepo:  marketplaceInfra.NewMarketplaceQueryRepo(persistentDbSvc),
		marketplaceCmdRepo:    marketplaceInfra.NewMarketplaceCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
		persistentDbSvc:       persistentDbSvc,
	}
}

func (service *MarketplaceService) ReadCatalog(
	input map[string]interface{},
) ServiceOutput {
	var idPtr *valueObject.MarketplaceItemId
	if input["id"] != nil {
		id, err := valueObject.NewMarketplaceItemId(input["id"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		idPtr = &id
	}

	var slugPtr *valueObject.MarketplaceItemSlug
	if input["slug"] != nil {
		slug, err := valueObject.NewMarketplaceItemSlug(input["slug"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		slugPtr = &slug
	}

	var namePtr *valueObject.MarketplaceItemName
	if input["name"] != nil {
		name, err := valueObject.NewMarketplaceItemName(input["name"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		namePtr = &name
	}

	var typePtr *valueObject.MarketplaceItemType
	if input["type"] != nil {
		itemType, err := valueObject.NewMarketplaceItemType(input["type"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		typePtr = &itemType
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			input["sortDirection"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
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

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(
		service.persistentDbSvc,
	)
	itemsList, err := useCase.ReadMarketplaceCatalogItems(marketplaceQueryRepo, readDto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, itemsList)
}

func (service *MarketplaceService) InstallCatalogItem(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
	hostname, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	if input["hostname"] != nil {
		hostname, err = valueObject.NewFqdn(input["hostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	var idPtr *valueObject.MarketplaceItemId
	if input["id"] != nil {
		id, err := valueObject.NewMarketplaceItemId(input["id"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		idPtr = &id
	}

	var slugPtr *valueObject.MarketplaceItemSlug
	if input["slug"] != nil {
		slug, err := valueObject.NewMarketplaceItemSlug(input["slug"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		slugPtr = &slug
	}

	var urlPathPtr *valueObject.UrlPath
	if input["urlPath"] != nil {
		urlPath, err := valueObject.NewUrlPath(input["urlPath"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		urlPathPtr = &urlPath
	}

	dataFields := []valueObject.MarketplaceInstallableItemDataField{}
	if _, exists := input["dataFields"]; exists {
		var assertOk bool
		dataFields, assertOk = input["dataFields"].([]valueObject.MarketplaceInstallableItemDataField)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidDataFields")
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
			escapedField := shellescape.Quote(dataField.String())
			installParams = append(installParams, "--data-fields", escapedField)
		}

		cliCmd += " " + strings.Join(installParams, " ")

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
		taskName, _ := valueObject.NewScheduledTaskName("InstallMarketplaceCatalogItem")
		taskCmd, _ := valueObject.NewUnixCommand(cliCmd)
		taskTag, _ := valueObject.NewScheduledTaskTag("marketplace")
		taskTags := []valueObject.ScheduledTaskTag{taskTag}
		timeoutSeconds := uint16(600)

		scheduledTaskCreateDto := dto.NewCreateScheduledTask(
			taskName, taskCmd, taskTags, &timeoutSeconds, nil,
		)

		err = useCase.CreateScheduledTask(scheduledTaskCmdRepo, scheduledTaskCreateDto)
		if err != nil {
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "MarketplaceCatalogItemInstallationScheduled")
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	dto := dto.NewInstallMarketplaceCatalogItem(
		hostname, idPtr, slugPtr, urlPathPtr, dataFields, operatorAccountId,
		operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)
	vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(service.persistentDbSvc)

	err = useCase.InstallMarketplaceCatalogItem(
		service.marketplaceQueryRepo, service.marketplaceCmdRepo, vhostQueryRepo,
		vhostCmdRepo, service.activityRecordCmdRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MarketplaceCatalogItemInstalled")
}

func (service *MarketplaceService) ReadInstalledItems(
	input map[string]interface{},
) ServiceOutput {
	var idPtr *valueObject.MarketplaceItemId
	if input["id"] != nil {
		id, err := valueObject.NewMarketplaceItemId(input["id"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		idPtr = &id
	}

	var hostnamePtr *valueObject.Fqdn
	if input["hostname"] != nil {
		hostname, err := valueObject.NewFqdn(input["hostname"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		hostnamePtr = &hostname
	}

	var typePtr *valueObject.MarketplaceItemType
	if input["type"] != nil {
		itemType, err := valueObject.NewMarketplaceItemType(input["type"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		typePtr = &itemType
	}

	var installationUuidPtr *valueObject.MarketplaceInstalledItemUuid
	if input["installationUuid"] != nil {
		installationUuid, err := valueObject.NewMarketplaceInstalledItemUuid(
			input["installationUuid"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		installationUuidPtr = &installationUuid
	}

	paginationDto := useCase.MarketplaceDefaultPagination
	if input["pageNumber"] != nil {
		pageNumber, err := voHelper.InterfaceToUint32(input["pageNumber"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidPageNumber"))
		}
		paginationDto.PageNumber = pageNumber
	}

	if input["itemsPerPage"] != nil {
		itemsPerPage, err := voHelper.InterfaceToUint16(input["itemsPerPage"])
		if err != nil {
			return NewServiceOutput(UserError, errors.New("InvalidItemsPerPage"))
		}
		paginationDto.ItemsPerPage = itemsPerPage
	}

	if input["sortBy"] != nil {
		sortBy, err := valueObject.NewPaginationSortBy(input["sortBy"])
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortBy = &sortBy
	}

	if input["sortDirection"] != nil {
		sortDirection, err := valueObject.NewPaginationSortDirection(
			input["sortDirection"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err)
		}
		paginationDto.SortDirection = &sortDirection
	}

	if input["lastSeenId"] != nil {
		lastSeenId, err := valueObject.NewPaginationLastSeenId(input["lastSeenId"])
		if err != nil {
			return NewServiceOutput(UserError, err)
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

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(
		service.persistentDbSvc,
	)
	itemsList, err := useCase.ReadMarketplaceInstalledItems(
		marketplaceQueryRepo, readDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, itemsList)
}

func (service *MarketplaceService) DeleteInstalledItem(
	input map[string]interface{},
	shouldSchedule bool,
) ServiceOutput {
	requiredParams := []string{"installedId"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	installedId, err := valueObject.NewMarketplaceItemId(input["installedId"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	shouldUninstallServices := true
	if input["shouldUninstallServices"] != nil {
		shouldUninstallServices, err = voHelper.InterfaceToBool(
			input["shouldUninstallServices"],
		)
		if err != nil {
			shouldUninstallServices = false
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
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
			service.persistentDbSvc,
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
			return NewServiceOutput(InfraError, err.Error())
		}

		return NewServiceOutput(Created, "MarketplaceCatalogItemDeletionScheduled")
	}

	deleteMarketplaceInstalledItem := dto.NewDeleteMarketplaceInstalledItem(
		installedId, shouldUninstallServices, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteMarketplaceInstalledItem(
		service.marketplaceQueryRepo, service.marketplaceCmdRepo,
		service.activityRecordCmdRepo, deleteMarketplaceInstalledItem,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "MarketplaceInstalledItemDeleted")
}
