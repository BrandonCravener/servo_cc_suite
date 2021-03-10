package logger

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/logger"
)

// Logger contains the current logger
var Logger *logger.Logger

// Init sets up logging
func Init() {
	currentTime := time.Now()
	date := strconv.Itoa(currentTime.Year()) + "-" + strconv.Itoa(currentTime.YearDay())
	logPath := "./logs/" + date

	if err := os.MkdirAll(logPath, 0755); err != nil {
		panic(err)
	}
	logFilePath := logPath + "/" + strings.ReplaceAll(strings.Split(strings.ReplaceAll(currentTime.String(), ".", " "), " ")[1], ":", "-") + ".log"
	logFilePath, _ = filepath.Abs(logFilePath)
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

	if err != nil {
		logger.Fatalf("Failed to open log file [%v]", err)
	}

	loggerName := "servo_cc_backend"

	Logger = logger.Init(loggerName, true, true, logFile)

	Logger.Infof("Logger successfully initalized to file and systemlog as [%s]", loggerName)
}
