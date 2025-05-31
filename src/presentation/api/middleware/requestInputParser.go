package apiMiddleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"github.com/labstack/echo/v4"
)

func bodyParserJson(
	echoContext echo.Context,
	rawRequestInput map[string]any,
) error {
	bodyBytes, err := io.ReadAll(echoContext.Request().Body)
	if err != nil {
		slog.Debug(
			"ReadJsonBodyFailed", slog.String("error", err.Error()),
		)
		return echo.NewHTTPError(http.StatusBadRequest, "ReadJsonBodyFailed")
	}
	echoContext.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

	if err := json.Unmarshal(bodyBytes, &rawRequestInput); err != nil {
		slog.Debug(
			"InvalidJsonBody", slog.String("error", err.Error()),
		)
		return echo.NewHTTPError(http.StatusBadRequest, "InvalidJsonBody")
	}

	return nil
}

func stringDotNotationToHierarchicalMap(
	hierarchicalMap map[string]any,
	remainingKeys []string,
	finalValue []string,
) map[string]any {
	if len(remainingKeys) == 1 {
		hierarchicalMap[remainingKeys[0]] = finalValue
		return hierarchicalMap
	}

	parentKey := remainingKeys[0]
	nextKeys := remainingKeys[1:]

	if _, exists := hierarchicalMap[parentKey]; !exists {
		hierarchicalMap[parentKey] = make(map[string]any)
	}

	if _, assertOk := hierarchicalMap[parentKey].(map[string]any); !assertOk {
		return hierarchicalMap
	}

	hierarchicalMap[parentKey] = stringDotNotationToHierarchicalMap(
		hierarchicalMap[parentKey].(map[string]any), nextKeys, finalValue,
	)

	return hierarchicalMap
}

func bodyParserFormData(
	echoContext echo.Context,
	rawRequestInput map[string]any,
) error {
	formData, err := echoContext.FormParams()
	if err != nil {
		slog.Debug(
			"InvalidFormData", slog.String("error", err.Error()),
		)
		return echo.NewHTTPError(http.StatusBadRequest, "InvalidFormData")
	}

	for formKey, keyValues := range formData {
		rawRequestInput[formKey] = keyValues

		isNestedKey := strings.Contains(formKey, ".")
		if isNestedKey {
			keyParts := strings.Split(formKey, ".")
			if len(keyParts) < 2 {
				continue
			}

			rawRequestInput = stringDotNotationToHierarchicalMap(
				rawRequestInput, keyParts, keyValues,
			)
		}
	}
	return nil
}

func bodyParserMultipartFormData(
	echoContext echo.Context,
	rawRequestInput map[string]any,
) error {
	multipartForm, err := echoContext.MultipartForm()
	if err != nil {
		slog.Debug(
			"InvalidMultipartFormData", slog.String("error", err.Error()),
		)
		return echo.NewHTTPError(http.StatusBadRequest, "InvalidMultipartFormData")
	}

	for formKey, keyValues := range multipartForm.Value {
		rawRequestInput[formKey] = keyValues
	}

	if len(multipartForm.File) > 0 {
		requestFileHeaders := map[string]*multipart.FileHeader{}
		for fileKey, fileHandlers := range multipartForm.File {
			isSingleFile := len(fileHandlers) == 1
			if isSingleFile {
				requestFileHeaders[fileKey] = fileHandlers[0]
				continue
			}

			for fileIndex, fileHandler := range fileHandlers {
				adjustedFileName := fileKey + "_" + strconv.Itoa(fileIndex)
				requestFileHeaders[adjustedFileName] = fileHandler
			}
		}
		rawRequestInput["files"] = requestFileHeaders
	}
	return nil
}

func RequestInputParser() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			rawRequestInput := map[string]any{}

			for _, paramName := range echoContext.ParamNames() {
				rawRequestInput[paramName] = echoContext.Param(paramName)
			}

			for queryParamName, queryParamValues := range echoContext.QueryParams() {
				rawRequestInput[queryParamName] = queryParamValues
			}

			switch contentType := echoContext.Request().Header.Get("Content-Type"); {
			case strings.HasPrefix(contentType, "application/json"):
				if err := bodyParserJson(echoContext, rawRequestInput); err != nil {
					return err
				}
			case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
				if err := bodyParserFormData(echoContext, rawRequestInput); err != nil {
					return err
				}
			case strings.HasPrefix(contentType, "multipart/form-data"):
				if err := bodyParserMultipartFormData(echoContext, rawRequestInput); err != nil {
					return err
				}
			default:
				return echo.NewHTTPError(http.StatusBadRequest, "UnsupportedContentType")
			}

			rawRequestInput["operatorIpAddress"] = echoContext.RealIP()
			if echoContext.Get("operatorAccountId") != nil {
				rawRequestInput["operatorAccountId"] = echoContext.Get("operatorAccountId")
			}

			requestInputParsed := tkPresentation.RequestInputParser(
				tkPresentation.RequestInputSettings{
					RawRequestInput:        rawRequestInput,
					KnownParamConstructors: sharedHelper.KnownParamConstructors,
				},
			)
			echoContext.Set("RequestInputParsed", requestInputParsed)

			return next(echoContext)
		}
	}
}
