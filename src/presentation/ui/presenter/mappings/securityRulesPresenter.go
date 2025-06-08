package uiPresenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type MappingSecurityRulesPresenter struct {
	virtualHostLiaison *liaison.VirtualHostLiaison
}

func NewMappingSecurityRulesPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MappingSecurityRulesPresenter {
	return &MappingSecurityRulesPresenter{
		virtualHostLiaison: liaison.NewVirtualHostLiaison(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *MappingSecurityRulesPresenter) Handler(c echo.Context) error {
	requestBody := map[string]interface{}{
		"itemsPerPage": 1000,
	}

	secRulesServiceResponse := presenter.virtualHostLiaison.ReadMappingSecurityRules(requestBody)
	if secRulesServiceResponse.Status != liaison.Success {
		slog.Debug("SecRulesServiceBadOutput", slog.Any("output", secRulesServiceResponse))
		return nil
	}

	secRulesReadResponse, assertOk := secRulesServiceResponse.Body.(dto.ReadMappingSecurityRulesResponse)
	if !assertOk {
		slog.Debug("SecRulesServiceResponseAssertionFailed")
		return nil
	}

	pageContent := MappingSecurityRulesIndex(secRulesReadResponse.MappingSecurityRules)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
