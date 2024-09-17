package ui

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	voHelper "github.com/speedianet/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	"github.com/speedianet/os/src/presentation/api"
	"github.com/speedianet/os/src/presentation/ui/presenter"
	"golang.org/x/net/websocket"
)

type Router struct {
	baseRoute       *echo.Group
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
}

func NewRouter(
	baseRoute *echo.Group,
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *Router {
	return &Router{
		baseRoute:       baseRoute,
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
	}
}

//go:embed dist/*
var previousDashFiles embed.FS

//go:embed assets/*
var assetsFiles embed.FS

func (router *Router) assetsRoute() {
	assetsFs, err := fs.Sub(assetsFiles, "assets")
	if err != nil {
		slog.Error("ReadAssetsFilesError", slog.Any("error", err))
		os.Exit(1)
	}
	assetsFileServer := http.FileServer(http.FS(assetsFs))

	router.baseRoute.GET(
		"/assets/*",
		echo.WrapHandler(http.StripPrefix("/assets/", assetsFileServer)),
	)
}

func (router *Router) databasesRoutes() {
	databaseGroup := router.baseRoute.Group("/databases")

	databasesPresenter := presenter.NewDatabasesPresenter(router.persistentDbSvc)
	databaseGroup.GET("/", databasesPresenter.Handler)
}

func (router *Router) mappingsRoutes() {
	mappingsGroup := router.baseRoute.Group("/mappings")

	mappingsPresenter := presenter.NewMappingsPresenter(router.persistentDbSvc)
	mappingsGroup.GET("/", mappingsPresenter.Handler)
}

func (router *Router) sslsRoutes() {
	sslsGroup := router.baseRoute.Group("/ssls")

	sslsPresenter := presenter.NewSslsPresenter(
		router.persistentDbSvc, router.transientDbSvc,
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

func (router *Router) previousDashboardRoute() {
	dashFilesFs, err := fs.Sub(previousDashFiles, "dist")
	if err != nil {
		slog.Error("ReadPreviousDashFilesError", slog.Any("error", err))
		os.Exit(1)
	}
	dashFileServer := http.FileServer(http.FS(dashFilesFs))

	previousDashGroup := router.baseRoute.Group("/_")
	previousDashGroup.GET(
		"/*", echo.WrapHandler(http.StripPrefix("/_", dashFileServer)),
	)
}

func (router *Router) RegisterRoutes() {
	router.assetsRoute()
	router.databasesRoutes()
	router.mappingsRoutes()
	router.sslsRoutes()

	if isDevMode, _ := voHelper.InterfaceToBool(os.Getenv("DEV_MODE")); isDevMode {
		router.devRoutes()
	}

	router.previousDashboardRoute()

	router.baseRoute.RouteNotFound("/*", func(c echo.Context) error {
		urlPath := c.Request().URL.Path
		isApi := strings.HasPrefix(urlPath, api.ApiBasePath)
		if isApi {
			return c.NoContent(http.StatusNotFound)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/_/")
	})
}
