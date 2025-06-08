package uiPresenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type DatabasesPresenter struct {
	databaseLiaison *liaison.DatabaseLiaison
}

func NewDatabasesPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *DatabasesPresenter {
	return &DatabasesPresenter{
		databaseLiaison: liaison.NewDatabaseLiaison(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *DatabasesPresenter) databaseOverviewFactory(
	rawDatabaseType string,
) (databaseOverview DatabaseOverview, err error) {
	databaseType, err := valueObject.NewDatabaseType(rawDatabaseType)
	if err != nil {
		return databaseOverview, err
	}

	isInstalled := false
	databaseEntities := []entity.Database{}
	databaseOverview = DatabaseOverview{
		databaseType, isInstalled, databaseEntities,
	}

	requestBody := map[string]interface{}{
		"dbType":       databaseType.String(),
		"itemsPerPage": 1000,
	}
	responseOutput := presenter.databaseLiaison.Read(requestBody)
	if responseOutput.Status != liaison.Success {
		return databaseOverview, err
	}

	responseDto, assertOk := responseOutput.Body.(dto.ReadDatabasesResponse)
	if assertOk {
		databaseOverview.IsInstalled = true
		databaseOverview.Databases = responseDto.Databases
	}

	return databaseOverview, nil
}

func (presenter *DatabasesPresenter) Handler(c echo.Context) error {
	rawDatabaseType := "mariadb"
	if c.QueryParam("dbType") != "" {
		rawDatabaseType = c.QueryParam("dbType")
	}

	selectedDatabaseOverview, err := presenter.databaseOverviewFactory(
		rawDatabaseType,
	)
	if err != nil {
		slog.Error("DatabaseOverviewFactoryError", slog.String("err", err.Error()))
		return nil
	}

	pageContent := DatabasesIndex(selectedDatabaseOverview)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
