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
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	marketplaceInfra "github.com/goinfinite/os/src/infra/marketplace"
	scheduledTaskInfra "github.com/goinfinite/os/src/infra/scheduledTask"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type MarketplaceService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewMarketplaceService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplaceService {
	return &MarketplaceService{
		persistentDbSvc: persistentDbSvc,
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
		sortDirection, err := valueObject.NewPaginationSortDirection(input["sortDirection"])
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
		Pagination: paginationDto,
		ItemId:     idPtr,
		ItemSlug:   slugPtr,
		ItemName:   namePtr,
		ItemType:   typePtr,
	}

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(service.persistentDbSvc)
	itemsList, err := useCase.ReadMarketplaceCatalog(marketplaceQueryRepo, readDto)
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
			installParams = append(installParams, "--urlPath", urlPathPtr.String())
		}

		for _, dataField := range dataFields {
			escapedField := shellescape.Quote(dataField.String())
			installParams = append(installParams, "--dataFields", escapedField)
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

	dto := dto.NewInstallMarketplaceCatalogItem(
		hostname, idPtr, slugPtr, urlPathPtr, dataFields,
	)

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(service.persistentDbSvc)
	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(service.persistentDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)
	vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(service.persistentDbSvc)

	err = useCase.InstallMarketplaceCatalogItem(
		marketplaceQueryRepo, marketplaceCmdRepo, vhostQueryRepo, vhostCmdRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MarketplaceCatalogItemInstalled")
}

func (service *MarketplaceService) ReadInstalledItems() ServiceOutput {
	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(service.persistentDbSvc)
	itemsList, err := useCase.ReadMarketplaceInstalledItems(marketplaceQueryRepo)
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

	if shouldSchedule {
		cliCmd := infraEnvs.InfiniteOsBinary + " mktplace delete"
		installParams := []string{
			"--installed-id", installedId.String(),
			"--should-uninstall-services", strconv.FormatBool(shouldUninstallServices),
		}

		cliCmd += " " + strings.Join(installParams, " ")

		scheduledTaskCmdRepo := scheduledTaskInfra.NewScheduledTaskCmdRepo(service.persistentDbSvc)
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
		installedId, shouldUninstallServices,
	)

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(service.persistentDbSvc)
	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteMarketplaceInstalledItem(
		marketplaceQueryRepo, marketplaceCmdRepo, deleteMarketplaceInstalledItem,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "MarketplaceInstalledItemDeleted")
}
