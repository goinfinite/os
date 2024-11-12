package apiHelper

import (
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func StringDotNotationToHierarchicalMap(
	hierarchicalMap map[string]interface{}, remainingKeys []string, finalValue string,
) map[string]interface{} {
	if len(remainingKeys) == 1 {
		hierarchicalMap[remainingKeys[0]] = finalValue
		return hierarchicalMap
	}

	parentKey := remainingKeys[0]
	nextKeys := remainingKeys[1:]

	if _, exists := hierarchicalMap[parentKey]; !exists {
		hierarchicalMap[parentKey] = make(map[string]interface{})
	}

	hierarchicalMap[parentKey] = StringDotNotationToHierarchicalMap(
		hierarchicalMap[parentKey].(map[string]interface{}), nextKeys, finalValue,
	)

	return hierarchicalMap
}

func ReadRequestBody(c echo.Context) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{}

	contentType := c.Request().Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(contentType, "application/json"):
		if err := c.Bind(&requestBody); err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidJsonBody")
		}
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		formData, err := c.FormParams()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidFormData")
		}

		for formKey, keyValues := range formData {
			if len(keyValues) == 0 {
				continue
			}

			if len(keyValues) > 1 {
				requestBody[formKey] = keyValues
				continue
			}

			keyValue := keyValues[0]
			isNestedKey := strings.Contains(formKey, ".")
			if !isNestedKey {
				requestBody[formKey] = keyValue
				continue
			}

			keyParts := strings.Split(formKey, ".")
			if len(keyParts) < 2 {
				continue
			}

			requestBody = StringDotNotationToHierarchicalMap(requestBody, keyParts, keyValue)
		}
	case strings.HasPrefix(contentType, "multipart/form-data"):
		multipartForm, err := c.MultipartForm()
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidMultipartFormData")
		}

		for formKey, keyValue := range multipartForm.Value {
			if len(keyValue) != 1 {
				continue
			}

			requestBody[formKey] = keyValue[0]
		}

		if len(multipartForm.File) > 0 {
			multipartFiles := map[string]*multipart.FileHeader{}
			for fileKey, fileValues := range multipartForm.File {
				if len(fileValues) != 1 {
					continue
				}

				multipartFiles[fileKey] = fileValues[0]
			}
			requestBody["files"] = multipartFiles
		}
	default:
		return nil, echo.NewHTTPError(http.StatusBadRequest, "InvalidContentType")
	}

	for queryParamName, queryParamValues := range c.QueryParams() {
		requestBody[queryParamName] = queryParamValues[0]
	}

	for _, paramName := range c.ParamNames() {
		requestBody[paramName] = c.Param(paramName)
	}

	requestBody["operatorAccountId"] = c.Get("accountId")
	requestBody["operatorIpAddress"] = c.RealIP()

	return requestBody, nil
}
