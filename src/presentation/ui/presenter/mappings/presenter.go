package uiPresenter

import (
	"log/slog"
	"net/http"
	"slices"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	componentMappings "github.com/goinfinite/os/src/presentation/ui/component/mappings"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
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
	virtualHostLiaison := liaison.NewVirtualHostLiaison(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	readVirtualHostsLiaisonOutput := virtualHostLiaison.ReadWithMappings(map[string]interface{}{
		"itemsPerPage": 1000,
	})
	if readVirtualHostsLiaisonOutput.Status != liaison.Success {
		slog.Debug("ReadMappingsLiaisonOutputBadStatus")
		return nil
	}

	readVirtualHostsResponse, assertOk := readVirtualHostsLiaisonOutput.Body.(dto.ReadVirtualHostsResponse)
	if !assertOk {
		slog.Debug("ReadMappingsLiaisonOutputBodyAssertionFailed")
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
	servicesLiaison := liaison.NewServicesLiaison(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	installedServicesResponseOutput := servicesLiaison.ReadInstalledItems(
		map[string]interface{}{"itemsPerPage": 1000},
	)
	if installedServicesResponseOutput.Status != liaison.Success {
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
	virtualHostLiaison := liaison.NewVirtualHostLiaison(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	readSecRulesLiaisonOutput := virtualHostLiaison.ReadMappingSecurityRules(map[string]interface{}{
		"itemsPerPage": 1000,
	})

	if readSecRulesLiaisonOutput.Status != liaison.Success {
		slog.Debug("ReadSecRulesLiaisonOutputBadStatus")
		return nil
	}

	readSecRulesResponse, assertOk := readSecRulesLiaisonOutput.Body.(dto.ReadMappingSecurityRulesResponse)
	if !assertOk {
		slog.Debug("ReadSecRulesLiaisonOutputBodyAssertionFailed")
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

	pageContent := MappingsIndex(
		vhostWithMappings, presenter.extractVirtualHostHostnames(vhostWithMappings),
		presenter.readInstalledServiceNames(),
		presenter.readSecRulesLabelValueOptions(),
	)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
