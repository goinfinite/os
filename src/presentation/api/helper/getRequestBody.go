package helper

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetRequestBody(c echo.Context) (map[string]interface{}, error) {
	requestData := map[string]interface{}{}

	contentType := c.Request().Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		if err := c.Bind(&requestData); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidJsonBody")
		}
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"), strings.HasPrefix(contentType, "multipart/form-data"):
		formData, err := c.FormParams()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidFormData")
		}
		for k, v := range formData {
			if len(v) > 0 {
				requestData[k] = v[0]
			}
		}
	default:
		return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidContentType")
	}

	return requestData, echo.NewHTTPError(http.StatusBadRequest, "EmptyRequestBody")
}
