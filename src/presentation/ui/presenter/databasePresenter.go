package presenter

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	"github.com/speedianet/os/src/presentation/service"
	uiHelper "github.com/speedianet/os/src/presentation/ui/helper"
	"github.com/speedianet/os/src/presentation/ui/page"
	presenterDto "github.com/speedianet/os/src/presentation/ui/presenter/dto"
)

type DatabasesPresenter struct {
	databaseService *service.DatabaseService
}

func NewDatabasesPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *DatabasesPresenter {
	return &DatabasesPresenter{
		databaseService: service.NewDatabaseService(persistentDbSvc),
	}
}

func (presenter *DatabasesPresenter) getSelectedDatabaseTypeSummary(
	rawDatabaseType string,
) (selectedDbTypeSummary presenterDto.SelectedDatabaseTypeSummary, err error) {
	databaseType, err := valueObject.NewDatabaseType(rawDatabaseType)
	if err != nil {
		return selectedDbTypeSummary, err
	}

	isInstalled := false

	requestBody := map[string]interface{}{"dbType": databaseType.String()}
	responseOutput := presenter.databaseService.Read(requestBody)
	if responseOutput.Status == service.Success {
		isInstalled = true
	}

	databases, assertOk := responseOutput.Body.([]entity.Database)
	if !assertOk {
		isInstalled = false
	}

	return presenterDto.NewSelectedDatabaseTypeSummary(
		databaseType, isInstalled, databases,
	), nil
}

func (presenter *DatabasesPresenter) Handler(c echo.Context) error {
	rawDatabaseType := "mariadb"
	if c.QueryParam("dbType") != "" {
		rawDatabaseType = c.QueryParam("dbType")
	}

	selectedDbTypeSummary, err := presenter.getSelectedDatabaseTypeSummary(
		rawDatabaseType,
	)
	if err != nil {
		slog.Error("GetSelectedDatabaseTypeSummaryError", slog.Any("err", err))
		return nil
	}

	pageContent := page.DatabasesIndex(selectedDbTypeSummary)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
