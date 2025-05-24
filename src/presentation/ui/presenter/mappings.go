package presenter

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	componentMappings "github.com/goinfinite/os/src/presentation/ui/component/mappings"
	"github.com/goinfinite/os/src/presentation/ui/layout"
	uiPage "github.com/goinfinite/os/src/presentation/ui/page/mappings"
	uiForm "github.com/goinfinite/ui/src/form"
	"github.com/labstack/echo/v4"
)

type MappingsPresenter struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewMappingsPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MappingsPresenter {
	return &MappingsPresenter{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *MappingsPresenter) readVirtualHostWithMappings() []dto.VirtualHostWithMappings {
	virtualHostService := service.NewVirtualHostService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	readVirtualHostsServiceOutput := virtualHostService.ReadWithMappings(map[string]interface{}{
		"itemsPerPage": 1000,
	})
	if readVirtualHostsServiceOutput.Status != service.Success {
		slog.Debug("ReadMappingsServiceOutputBadStatus")
		return nil
	}

	readVirtualHostsResponse, assertOk := readVirtualHostsServiceOutput.Body.(dto.ReadVirtualHostsResponse)
	if !assertOk {
		slog.Debug("ReadMappingsServiceOutputBodyAssertionFailed")
		return nil
	}

	return readVirtualHostsResponse.VirtualHostWithMappings
}

func (presenter *MappingsPresenter) extractVirtualHostHostnames(
	vhostsWithMappings []dto.VirtualHostWithMappings,
) []string {
	vhostsHostnames := []string{}
	for _, vhostWithMappings := range vhostsWithMappings {
		vhostsHostnames = append(vhostsHostnames, vhostWithMappings.Hostname.String())
	}

	return vhostsHostnames
}

func (presenter *MappingsPresenter) readInstalledServiceNames() []string {
	servicesService := service.NewServicesService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	installedServicesResponseOutput := servicesService.ReadInstalledItems(
		map[string]interface{}{"itemsPerPage": 1000},
	)
	if installedServicesResponseOutput.Status != service.Success {
		slog.Debug("ReadInstalledItemsFailed", slog.Any("output", installedServicesResponseOutput))
		return nil
	}

	installedServicesResponseDto, assertOk := installedServicesResponseOutput.Body.(dto.ReadInstalledServicesItemsResponse)
	if !assertOk {
		slog.Debug("ReadInstalledItemsResponseDtoAssertionFailed")
		return nil
	}

	servicesNames := []string{}
	for _, serviceEntity := range installedServicesResponseDto.InstalledServices {
		if len(serviceEntity.PortBindings) == 0 {
			continue
		}
		servicesNames = append(servicesNames, serviceEntity.Name.String())
	}
	slices.Sort(servicesNames)

	return servicesNames
}

func (presenter *MappingsPresenter) readSecRulesLabelValueOptions() []uiForm.SelectLabelValueOption {
	virtualHostService := service.NewVirtualHostService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	readSecRulesServiceOutput := virtualHostService.ReadMappingSecurityRules(map[string]interface{}{
		"itemsPerPage": 1000,
	})

	if readSecRulesServiceOutput.Status != service.Success {
		slog.Debug("ReadSecRulesServiceOutputBadStatus")
		return nil
	}

	readSecRulesResponse, assertOk := readSecRulesServiceOutput.Body.(dto.ReadMappingSecurityRulesResponse)
	if !assertOk {
		slog.Debug("ReadSecRulesServiceOutputBodyAssertionFailed")
		return nil
	}

	secRulesLabelValueOptions := []uiForm.SelectLabelValueOption{}
	for _, secRuleEntity := range readSecRulesResponse.MappingSecurityRules {
		secRulesLabelValueOptions = append(
			secRulesLabelValueOptions,
			uiForm.SelectLabelValueOption{
				Label:     secRuleEntity.Name.String() + " (#" + secRuleEntity.Id.String() + ")",
				LabelHtml: componentMappings.MappingSecurityRuleSummary(secRuleEntity),
				Value:     secRuleEntity.Id.String(),
			})
	}

	return secRulesLabelValueOptions
}

func (presenter *MappingsPresenter) Handler(c echo.Context) error {
	vhostWithMappings := presenter.readVirtualHostWithMappings()

	pageContent := uiPage.MappingsIndex(
		vhostWithMappings, presenter.extractVirtualHostHostnames(vhostWithMappings),
		presenter.readInstalledServiceNames(),
		presenter.readSecRulesLabelValueOptions(),
	)
	return layout.Renderer(layout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
