package uiMiddleware

import (
	"embed"

	"github.com/labstack/echo/v4"
)

type EmbedKeyFs struct {
	EmbedKey string
	EmbedFs  embed.FS
}

func Embed(embeds []EmbedKeyFs) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, embed := range embeds {
				c.Set(embed.EmbedKey, embed.EmbedFs)
			}
			return next(c)
		}
	}
}
