package apiHelper

import (
	"mime/multipart"
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
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		formData, err := c.FormParams()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidFormData")
		}
		for k, v := range formData {
			if len(v) > 0 {
				requestData[k] = v[0]
			}
		}
	case strings.HasPrefix(contentType, "multipart/form-data"):
		multipartFormData, err := c.MultipartForm()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidMultipartFormData")
		}

		for k, v := range multipartFormData.Value {
			if len(v) > 0 {
				requestData[k] = v[0]
			}
		}

		if len(multipartFormData.File) > 0 {
			requestDataFiles := map[string]*multipart.FileHeader{}

			for k, v := range multipartFormData.File {
				if len(v) > 0 {
					requestDataFiles[k] = v[0]
				}
			}

			requestData["files"] = requestDataFiles
		}
	default:
		return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidContentType")
	}

	return requestData, echo.NewHTTPError(http.StatusBadRequest, "EmptyRequestBody")
}
