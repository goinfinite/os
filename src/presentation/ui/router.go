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
	"github.com/speedianet/os/src/presentation/api"
	"golang.org/x/net/websocket"
)

type Router struct {
	baseRoute *echo.Group
}

func NewRouter(baseRoute *echo.Group) *Router {
	return &Router{baseRoute: baseRoute}
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
