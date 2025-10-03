package ui

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	presenterAccounts "github.com/goinfinite/os/src/presentation/ui/presenter/accounts"
	presenterCrons "github.com/goinfinite/os/src/presentation/ui/presenter/crons"
	presenterDatabases "github.com/goinfinite/os/src/presentation/ui/presenter/databases"
	presenterFileManager "github.com/goinfinite/os/src/presentation/ui/presenter/fileManager"
	presenterFooter "github.com/goinfinite/os/src/presentation/ui/presenter/footer"
	presenterLogin "github.com/goinfinite/os/src/presentation/ui/presenter/login"
	presenterMappings "github.com/goinfinite/os/src/presentation/ui/presenter/mappings"
	presenterMarketplace "github.com/goinfinite/os/src/presentation/ui/presenter/marketplace"
	presenterOverview "github.com/goinfinite/os/src/presentation/ui/presenter/overview"
	presenterRuntimes "github.com/goinfinite/os/src/presentation/ui/presenter/runtimes"
	presenterSetup "github.com/goinfinite/os/src/presentation/ui/presenter/setup"
	presenterSsls "github.com/goinfinite/os/src/presentation/ui/presenter/ssls"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type Router struct {
	baseRoute       *echo.Group
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewRouter(
	baseRoute *echo.Group,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *Router {
	return &Router{
		baseRoute:       baseRoute,
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (router *Router) assetsRoute() {
	router.baseRoute.GET("/assets/*", func(c echo.Context) error {
		assetsFiles, assertOk := c.Get("assets").(embed.FS)
		if !assertOk {
			slog.Error("AssertAssetsFilesFailed")
			os.Exit(1)
		}

		assetsFs, err := fs.Sub(assetsFiles, "assets")
		if err != nil {
			slog.Error("ReadAssetsFilesError", slog.String("err", err.Error()))
			os.Exit(1)
		}
		assetsFileServer := http.FileServer(http.FS(assetsFs))

		http.StripPrefix("/assets/", assetsFileServer).
			ServeHTTP(c.Response(), c.Request())
		return nil
	})
}

func (router *Router) accountsRoutes() {
	accountGroup := router.baseRoute.Group("/accounts")

	accountsPresenter := presenterAccounts.NewAccountsPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	accountGroup.GET("/", accountsPresenter.Handler)
}

func (router *Router) cronsRoutes() {
	cronsGroup := router.baseRoute.Group("/crons")

	cronsPresenter := presenterCrons.NewCronsPresenter(router.trailDbSvc)
	cronsGroup.GET("/", cronsPresenter.Handler)
}

func (router *Router) databasesRoutes() {
	databaseGroup := router.baseRoute.Group("/databases")

	databasesPresenter := presenterDatabases.NewDatabasesPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	databaseGroup.GET("/", databasesPresenter.Handler)
}

func (router *Router) fileManagerRoutes() {
	fileManagerGroup := router.baseRoute.Group("/file-manager")

	fileManagerPresenter := presenterFileManager.NewFileManagerPresenter()
	fileManagerGroup.GET("/", fileManagerPresenter.Handler)
}

func (router *Router) loginRoutes() {
	loginGroup := router.baseRoute.Group("/login")

	loginPresenter := presenterLogin.NewLoginPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	loginGroup.GET("/", loginPresenter.Handler)
	loginGroup.HEAD("/", loginPresenter.Handler)
}

func (router *Router) mappingsRoutes() {
	mappingsGroup := router.baseRoute.Group("/mappings")

	mappingsPresenter := presenterMappings.NewMappingsPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	mappingsGroup.GET("/", mappingsPresenter.Handler)

	secRulesGroup := mappingsGroup.Group("/security-rules")
	secRulesPresenter := presenterMappings.NewMappingSecurityRulesPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	secRulesGroup.GET("/", secRulesPresenter.Handler)
}

func (router *Router) marketplaceRoutes() {
	marketplaceGroup := router.baseRoute.Group("/marketplace")

	marketplacePresenter := presenterMarketplace.NewMarketplacePresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	marketplaceGroup.GET("/", marketplacePresenter.Handler)
}

func (router *Router) overviewRoutes() {
	overviewGroup := router.baseRoute.Group("/overview")

	overviewPresenter := presenterOverview.NewOverviewPresenter(
		router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
	)
	overviewGroup.GET("/", overviewPresenter.Handler)
	overviewGroup.HEAD("/", overviewPresenter.Handler)
}

func (router *Router) runtimesRoutes() {
	runtimesGroup := router.baseRoute.Group("/runtimes")

	runtimesPresenter := presenterRuntimes.NewRuntimesPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	runtimesGroup.GET("/", runtimesPresenter.Handler)
}

func (router *Router) setupRoutes() {
	setupGroup := router.baseRoute.Group("/setup")

	setupPresenter := presenterSetup.NewSetupPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	setupGroup.GET("/", setupPresenter.Handler)
	setupGroup.HEAD("/", setupPresenter.Handler)
}

func (router *Router) sslsRoutes() {
	sslsGroup := router.baseRoute.Group("/ssls")

	sslsPresenter := presenterSsls.NewSslsPresenter(
		router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
	)
	sslsGroup.GET("/", sslsPresenter.Handler)
}

func (router *Router) devRoutes() {
	devGroup := router.baseRoute.Group("/dev")
	devGroup.GET("/hot-reload", func(c echo.Context) error {
		websocket.Handler(func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				err := websocket.Message.Send(ws, "WS Hot Reload Activated!")
				if err != nil {
					break
				}

				msgReceived := ""
				err = websocket.Message.Receive(ws, &msgReceived)
				if err != nil {
					break
				}
			}
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	})
}

func (router *Router) fragmentRoutes() {
	fragmentGroup := router.baseRoute.Group("/fragment")

	footerPresenter := presenterFooter.NewFooterPresenter(
		router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
	)
	fragmentGroup.GET("/footer", footerPresenter.Handler)
}

func (router *Router) RegisterRoutes() {
	router.assetsRoute()
	router.accountsRoutes()
	router.cronsRoutes()
	router.databasesRoutes()
	router.fileManagerRoutes()
	router.loginRoutes()
	router.mappingsRoutes()
	router.marketplaceRoutes()
	router.overviewRoutes()
	router.runtimesRoutes()
	router.setupRoutes()
	router.sslsRoutes()

	if isDevMode, _ := voHelper.InterfaceToBool(os.Getenv("DEV_MODE")); isDevMode {
		router.devRoutes()
	}

	router.fragmentRoutes()

	router.baseRoute.RouteNotFound("/*", func(echoContext echo.Context) error {
		apiBasePath, assertOk := echoContext.Get("apiBasePath").(string)
		if !assertOk {
			slog.Error("AssertApiBasePathFailed")
			return echoContext.NoContent(http.StatusInternalServerError)
		}

		urlPath := echoContext.Request().URL.Path
		isApi := strings.HasPrefix(urlPath, apiBasePath)
		if isApi {
			return echoContext.NoContent(http.StatusNotFound)
		}

		slog.Debug("RouteNotFound", slog.String("urlPath", urlPath))

		uiBasePath, assertOk := echoContext.Get("uiBasePath").(string)
		if !assertOk {
			slog.Error("AssertUiBasePathFailed")
			return echoContext.NoContent(http.StatusInternalServerError)
		}

		baseHref, assertOk := echoContext.Get("baseHref").(string)
		if !assertOk {
			slog.Error("AssertBaseHrefFailed")
			return echoContext.NoContent(http.StatusInternalServerError)
		}
		if len(baseHref) > 0 {
			baseHrefNoTrailing := strings.TrimSuffix(baseHref, "/")
			uiBasePath = baseHrefNoTrailing + uiBasePath
		}

		return echoContext.Redirect(http.StatusTemporaryRedirect, uiBasePath+"/overview/")
	})
}
