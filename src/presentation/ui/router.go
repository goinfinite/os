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
	"github.com/goinfinite/os/src/presentation/api"
	"github.com/goinfinite/os/src/presentation/ui/presenter"
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

	accountsPresenter := presenter.NewAccountsPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	accountGroup.GET("/", accountsPresenter.Handler)
}

func (router *Router) cronsRoutes() {
	cronsGroup := router.baseRoute.Group("/crons")

	cronsPresenter := presenter.NewCronsPresenter(router.trailDbSvc)
	cronsGroup.GET("/", cronsPresenter.Handler)
}

func (router *Router) databasesRoutes() {
	databaseGroup := router.baseRoute.Group("/databases")

	databasesPresenter := presenter.NewDatabasesPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	databaseGroup.GET("/", databasesPresenter.Handler)
}

func (router *Router) fileManagerRoutes() {
	fileManagerGroup := router.baseRoute.Group("/file-manager")

	fileManagerPresenter := presenter.NewFileManagerPresenter()
	fileManagerGroup.GET("/", fileManagerPresenter.Handler)
}

func (router *Router) loginRoutes() {
	loginGroup := router.baseRoute.Group("/login")

	loginPresenter := presenter.NewLoginPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	loginGroup.GET("/", loginPresenter.Handler)
	loginGroup.HEAD("/", loginPresenter.Handler)
}

func (router *Router) mappingsRoutes() {
	mappingsGroup := router.baseRoute.Group("/mappings")

	mappingsPresenter := presenter.NewMappingsPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	mappingsGroup.GET("/", mappingsPresenter.Handler)

	secRulesGroup := mappingsGroup.Group("/security-rules")
	secRulesPresenter := presenter.NewMappingSecurityRulesPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	secRulesGroup.GET("/", secRulesPresenter.Handler)
}

func (router *Router) marketplaceRoutes() {
	marketplaceGroup := router.baseRoute.Group("/marketplace")

	marketplacePresenter := presenter.NewMarketplacePresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	marketplaceGroup.GET("/", marketplacePresenter.Handler)
}

func (router *Router) overviewRoutes() {
	overviewGroup := router.baseRoute.Group("/overview")

	overviewPresenter := presenter.NewOverviewPresenter(
		router.persistentDbSvc, router.transientDbSvc, router.trailDbSvc,
	)
	overviewGroup.GET("/", overviewPresenter.Handler)
	overviewGroup.HEAD("/", overviewPresenter.Handler)
}

func (router *Router) runtimesRoutes() {
	runtimesGroup := router.baseRoute.Group("/runtimes")

	runtimesPresenter := presenter.NewRuntimesPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	runtimesGroup.GET("/", runtimesPresenter.Handler)
}

func (router *Router) setupRoutes() {
	setupGroup := router.baseRoute.Group("/setup")

	setupPresenter := presenter.NewSetupPresenter(
		router.persistentDbSvc, router.trailDbSvc,
	)
	setupGroup.GET("/", setupPresenter.Handler)
	setupGroup.HEAD("/", setupPresenter.Handler)
}

func (router *Router) sslsRoutes() {
	sslsGroup := router.baseRoute.Group("/ssls")

	sslsPresenter := presenter.NewSslsPresenter(
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

	footerPresenter := presenter.NewFooterPresenter(
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

	router.baseRoute.RouteNotFound("/*", func(c echo.Context) error {
		urlPath := c.Request().URL.Path
		isApi := strings.HasPrefix(urlPath, api.ApiBasePath)
		if isApi {
			return c.NoContent(http.StatusNotFound)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/overview/")
	})
}
