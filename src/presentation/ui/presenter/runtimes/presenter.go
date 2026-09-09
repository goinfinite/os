package uiPresenter

import (
	"log/slog"
	"net/http"

	tkPresentation "github.com/goinfinite/tk/src/presentation"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	"github.com/goinfinite/os/src/presentation/liaison"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
	presenterHelper "github.com/goinfinite/os/src/presentation/ui/presenter/helper"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	"github.com/labstack/echo/v4"
)

const unexpectedPhpConfigsResponseFormat string = "UnexpectedPhpConfigsResponseFormat"

type RuntimesPresenter struct {
	runtimeLiaison  *liaison.RuntimeLiaison
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewRuntimesPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *RuntimesPresenter {
	return &RuntimesPresenter{
		runtimeLiaison:  liaison.NewRuntimeLiaison(persistentDbSvc, trailDbSvc),
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

type phpRuntimeState struct {
	isPhpWebServerInstalled   bool
	isVirtualHostUsingRuntime bool
	configsReadFailure        string
	phpConfigsPtr             *entity.PhpConfigs
}

func (presenter *RuntimesPresenter) phpRuntimeStateFactory(
	responseOutput tkPresentation.LiaisonResponse,
) phpRuntimeState {
	switch responseOutput.Status {
	case tkPresentation.LiaisonResponseStatusSuccess:
		phpConfigs, assertOk := responseOutput.Body.(entity.PhpConfigs)
		if !assertOk {
			return phpRuntimeState{
				isPhpWebServerInstalled: true,
				configsReadFailure:      unexpectedPhpConfigsResponseFormat,
			}
		}
		return phpRuntimeState{
			isPhpWebServerInstalled:   true,
			isVirtualHostUsingRuntime: true,
			phpConfigsPtr:             &phpConfigs,
		}
	case tkPresentation.LiaisonResponseStatusNotFound:
		return phpRuntimeState{isPhpWebServerInstalled: true}
	case tkPresentation.LiaisonResponseStatusServiceUnavailable:
		return phpRuntimeState{}
	default:
		infraError, assertOk := responseOutput.Body.(string)
		if !assertOk {
			infraError = unexpectedPhpConfigsResponseFormat
		}
		return phpRuntimeState{
			isPhpWebServerInstalled: true,
			configsReadFailure:      infraError,
		}
	}
}

func (presenter *RuntimesPresenter) runtimeOverviewFactory(
	runtimeType valueObject.RuntimeType,
	selectedVhostHostname tkValueObject.Fqdn,
) RuntimeOverview {
	phpState := phpRuntimeState{}
	if runtimeType.String() == "php" {
		requestBody := map[string]any{"hostname": selectedVhostHostname.String()}
		responseOutput := presenter.runtimeLiaison.ReadPhpConfigs(requestBody)
		phpState = presenter.phpRuntimeStateFactory(responseOutput)
		if phpState.configsReadFailure != "" {
			slog.Error(
				"PhpConfigsReadFailure",
				slog.String("hostname", selectedVhostHostname.String()),
				slog.String("reason", phpState.configsReadFailure),
			)
		}
	}

	return RuntimeOverview{
		VirtualHostHostname:       selectedVhostHostname,
		Type:                      runtimeType,
		IsInstalled:               phpState.isPhpWebServerInstalled,
		IsVirtualHostUsingRuntime: phpState.isVirtualHostUsingRuntime,
		PhpConfigsReadFailure:     phpState.configsReadFailure,
		PhpConfigs:                phpState.phpConfigsPtr,
	}
}

func (presenter *RuntimesPresenter) Handler(c echo.Context) error {
	rawRuntimeType := "php"
	if c.QueryParam("runtimeType") != "" {
		rawRuntimeType = c.QueryParam("runtimeType")
	}
	runtimeType, err := valueObject.NewRuntimeType(rawRuntimeType)
	if err != nil {
		slog.Error("InvalidRuntimeType", slog.String("err", err.Error()))
		return nil
	}

	primaryVhostHostname, err := vhostInfra.NewVirtualHostHelpers().
		ReadPrimaryVirtualHostHostname()
	if err != nil {
		slog.Error("ReadPrimaryVirtualHost", slog.String("err", err.Error()))
		return nil
	}
	selectedVhostHostname := primaryVhostHostname
	if c.QueryParam("vhostHostname") != "" {
		selectedVhostHostname, err = tkValueObject.NewFqdn(c.QueryParam("vhostHostname"))
		if err != nil {
			slog.Error("InvalidVhostHostname", slog.String("err", err.Error()))
			return nil
		}
	}

	runtimeOverview := presenter.runtimeOverviewFactory(runtimeType, selectedVhostHostname)

	vhostsHostnames, err := presenterHelper.ReadVirtualHostHostnames(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	if err != nil {
		slog.Error("ReadVirtualHostsHostnames", slog.String("err", err.Error()))
		return nil
	}

	pageContent := RuntimesIndex(runtimeOverview, vhostsHostnames)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
