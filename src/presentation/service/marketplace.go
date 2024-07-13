package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	marketplaceInfra "github.com/speedianet/os/src/infra/marketplace"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
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

func (service *MarketplaceService) ReadCatalog() ServiceOutput {
	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(service.persistentDbSvc)
	itemsList, err := useCase.ReadMarketplaceCatalog(marketplaceQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, itemsList)
}

func (service *MarketplaceService) InstallCatalogItem(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"hostname"}

	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
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

	marketplaceQueryRepo := marketplaceInfra.NewMarketplaceQueryRepo(service.persistentDbSvc)
	marketplaceCmdRepo := marketplaceInfra.NewMarketplaceCmdRepo(service.persistentDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)
	vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(service.persistentDbSvc)

	dto := dto.NewInstallMarketplaceCatalogItem(
		hostname, idPtr, slugPtr, urlPathPtr, dataFields,
	)

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
		shouldUninstallServices, err = serviceHelper.ParseBoolParam(
			input["shouldUninstallServices"],
		)
		if err != nil {
			shouldUninstallServices = false
		}
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
