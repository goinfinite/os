package apiMiddleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
)

type StackTrace struct {
	trace string
}

func (st *StackTrace) String() string {
	return st.trace
}

func convertPanicMessageToError(rec interface{}) error {
	err, isError := rec.(error)
	if !isError {
		err = fmt.Errorf("%v", rec)
	}
	return err
}

func getStackTrace() *StackTrace {
	buf := make([]byte, 1<<16)
	runtime.Stack(buf, true)
	stackTraceString := string(buf)
	lines := strings.Split(stackTraceString, "\n")

	filteredLines := []string{}
	for _, line := range lines {
		filteredLines = append(filteredLines, line)
		if strings.Contains(line, "created by net/http") {
			break
		}
	}

	filteredStackTrace := strings.Join(filteredLines, "\n")

	return &StackTrace{trace: filteredStackTrace}
}

func isDomainLayerPanic(stackTrace *StackTrace) bool {
	valueObjectPath := filepath.Join("domain", "valueObject")
	entityPath := filepath.Join("domain", "entity")
	useCasePath := filepath.Join("domain", "useCase")

	return strings.Contains(stackTrace.String(), valueObjectPath) ||
		strings.Contains(stackTrace.String(), entityPath) ||
		strings.Contains(stackTrace.String(), useCasePath)
}

func isTrustworthy(c echo.Context) bool {
	currentIp := c.RealIP()
	trustedIps := strings.Split(os.Getenv("TRUSTED_IPS"), ",")

	for _, staffIp := range trustedIps {
		if currentIp == strings.TrimSpace(staffIp) {
			return true
		}
	}

	return false
}

func getFirstChars(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func createLogDirectory() error {
	logDir := "logs"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		return os.Mkdir(logDir, 0755)
	}
	return nil
}

func logPanic(err error, stackTrace *StackTrace) {
	if err := createLogDirectory(); err != nil {
		log.Printf("Error creating log directory: %v", err)
		return
	}

	logFile := filepath.Join("logs", "panic.log")
	file, errOpen := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errOpen != nil {
		log.Printf("Error opening log file: %v", errOpen)
		return
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)

	logger.Printf("Panic recovered: %s", err.Error())
	logger.Println("Stack trace:")
	logger.Printf("%s\n", stackTrace.String())
}

func handlePanic(c echo.Context) {
	rec := recover()
	if rec == nil {
		return
	}

	err := convertPanicMessageToError(rec)
	stackTrace := getStackTrace()
	uri := c.Request().RequestURI
	queryParams := c.QueryParams()

	statusCode := http.StatusInternalServerError
	if isDomainLayerPanic(stackTrace) {
		statusCode = http.StatusBadRequest
	}

	if isTrustworthy(c) {
		c.JSON(statusCode, map[string]interface{}{
			"status": statusCode,
			"body": map[string]interface{}{
				"uri":            uri,
				"queryParams":    queryParams,
				"exceptionCode":  err.Error(),
				"exceptionTrace": stackTrace.String(),
			},
		})
		return
	}

	shortErrMsg := getFirstChars(err.Error(), 150)
	if len(shortErrMsg) == 0 {
		shortErrMsg = "InternalServerError"
	}

	c.JSON(statusCode, map[string]interface{}{
		"status": statusCode,
		"body":   shortErrMsg,
	})

	logPanic(err, stackTrace)
}

func PanicHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer handlePanic(c)
		return next(c)
	}
}
