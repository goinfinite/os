package cliMiddleware

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
		if strings.Contains(line, "spf13/cobra") {
			break
		}
	}

	filteredStackTrace := strings.Join(filteredLines, "\n")

	return &StackTrace{trace: filteredStackTrace}
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

	logger.Printf("Panic: %s", err.Error())
	logger.Println("Stack trace:")
	logger.Printf("%s\n", stackTrace.String())
}

func PanicHandler() {
	rec := recover()
	if rec == nil {
		return
	}

	err := convertPanicMessageToError(rec)
	stackTrace := getStackTrace()

	errorMsg := "FatalError: " +
		err.Error() +
		". Please check panic log for more details."
	fmt.Println(errorMsg)

	logPanic(err, stackTrace)
	os.Exit(1)
}
