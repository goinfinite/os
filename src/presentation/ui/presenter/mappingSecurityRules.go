package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type MappingSecurityRulesPresenter struct {
	virtualHostService *service.VirtualHostService
}

func NewMappingSecurityRulesPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MappingSecurityRulesPresenter {
	return &MappingSecurityRulesPresenter{
		virtualHostService: service.NewVirtualHostService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *MappingSecurityRulesPresenter) Handler(c echo.Context) error {
	requestBody := map[string]interface{}{
		"itemsPerPage": 1000,
	}

	secRulesServiceResponse := presenter.virtualHostService.ReadMappingSecurityRules(requestBody)
	if secRulesServiceResponse.Status != service.Success {
		slog.Debug("SecRulesServiceBadOutput", slog.Any("output", secRulesServiceResponse))
		return nil
	}

	secRulesReadResponse, assertOk := secRulesServiceResponse.Body.(dto.ReadMappingSecurityRulesResponse)
	if !assertOk {
		slog.Debug("SecRulesServiceResponseAssertionFailed")
		return nil
	}

	pageContent := page.MappingSecurityRulesIndex(secRulesReadResponse.MappingSecurityRules)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
