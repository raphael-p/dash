package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	MAX_FILE_COUNT    = 1000
	MAX_FILE_BYTES    = 10 * 1024 * 1024
	MAX_MESSAGE_BYTES = 10 * 1024
)

type MyLogger struct {
	stdOutLogger *log.Logger
	fileLogger   *log.Logger
	logLevel     logLevel
	cumBytes     int64
	workingDir   string
}

type logLevel uint8

const (
	logLevelTrace logLevel = iota
	logLevelDebug
	logLevelInfo
	logLevelWarning
	logLevelError
)

type nowFunc func() string

var now nowFunc = func() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05.000")
}

func newLogger(out io.Writer) *log.Logger {
	return log.New(out, "", 0)
}

var Logger *MyLogger = &MyLogger{nil, nil, logLevel(0), 0, "."}

func openLogFile() *os.File {
	directory := filepath.Join(Logger.workingDir, "logs")
	path := filepath.Join(directory, generateFilename("server", directory))
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		err = os.MkdirAll(directory, 0755)
		if err != nil {
			panic(fmt.Sprintf("failed to create log directory: %s", err))
		}

		logFile, err = os.Create(path)
		if err != nil {
			panic(fmt.Sprintf("failed to create log file: %s", err))
		}
	}
	return logFile
}

func generateFilename(name, directory string) string {
	filenameArr := strings.SplitN(name, ".", 2)
	basename := filenameArr[0]
	extension := ".log"
	if len(filenameArr) > 1 && filenameArr[1] != "" {
		extension = "." + filenameArr[1]
	}
	basename += "-" + time.Now().UTC().Format("20060102")

	count := 1
	lastUsedFile := ""
	for {
		filename := fmt.Sprint(basename, "-", count, extension)
		path := filepath.Join(directory, filename)
		stats, err := os.Stat(path)
		if err == nil && stats.Size() < MAX_FILE_BYTES {
			lastUsedFile = filename
		} else if err != nil && lastUsedFile != "" {
			return lastUsedFile
		} else if err != nil {
			return filename
		}

		count++
		if count > MAX_FILE_COUNT {
			panic(fmt.Sprint(
				"reached maximum allowed number of log files: ",
				MAX_FILE_COUNT,
			))
		}
	}
}

// Creates a logger which manages its logs files and logs to the console.
func Create(workingDir string, hasConsoleLog bool) {
	Logger.logLevel = logLevel(0)
	if hasConsoleLog {
		Logger.stdOutLogger = newLogger(os.Stdout)
	}
	Logger.workingDir = workingDir
	Logger.fileLogger = newLogger(openLogFile())
	Info("file logger initialised")
}

// Handles the closing of the log file
func Close() {
	if Logger.fileLogger == nil {
		return
	}
	if file, ok := Logger.fileLogger.Writer().(*os.File); ok {
		err := file.Close()
		if err != nil {
			panic(fmt.Sprint("failed to close log file: ", err))
		}
		Info("log file closed")
	}
}

func logMessage(level string, ansiColour string, message string) {
	reset := "\033[0m"
	time := now()
	if len(message) > MAX_MESSAGE_BYTES {
		message = message[:MAX_MESSAGE_BYTES]
	}
	if Logger.stdOutLogger != nil {
		Logger.stdOutLogger.Printf("%s %s[%s]%s %s", time, ansiColour, level, reset, message)
	}
	if Logger.fileLogger != nil {
		if Logger.cumBytes += int64(len(message)); Logger.cumBytes > MAX_FILE_BYTES {
			rollover()
		}
		Logger.fileLogger.Printf("%s [%s] %s", time, level, message)
	}
}

func rollover() {
	if file, ok := Logger.fileLogger.Writer().(*os.File); ok {
		if stats, err := file.Stat(); err == nil && stats.Size() > MAX_FILE_BYTES {
			Logger.cumBytes = 0
			Info(fmt.Sprintf(
				"reached maximum log file size (%d bytes), rolling over",
				MAX_FILE_BYTES,
			)) // can cause infinite loop if MAX_FILE_BYTES is too low
			Close()
			Logger.fileLogger = newLogger(openLogFile())
			Info("new log file opened")
		}
	}
}

func Trace(message string) {
	if Logger.logLevel <= logLevelTrace {
		logMessage("TRACE", "", message)
	}
}

func Debug(message string) {
	if Logger.logLevel <= logLevelDebug {
		logMessage("DEBUG", "\033[34m", message)
	}
}

func Info(message string) {
	if Logger.logLevel <= logLevelInfo {
		logMessage("INFO", "\033[36m", message)
	}
}

func Warning(message string) {
	if Logger.logLevel <= logLevelWarning {
		logMessage("WARNING", "\033[33;1m", message)
	}
}

func Error(message string) {
	if Logger.logLevel <= logLevelError {
		buf := make([]byte, 1<<16)
		n := runtime.Stack(buf, false)
		stackTrace := strings.ReplaceAll(string(buf[:n-1]), "\n", "\n\t")
		message += fmt.Sprintf("\n\tstack trace: %s", stackTrace)
		logMessage("ERROR", "\033[31;1m", message)
	}
}

func Fatal(message string) {
	Error(message)
	if Logger.stdOutLogger == nil {
		os.Stdout.WriteString(message + "\n")
	}
	os.Exit(1)
}

func Tracef(format string, a ...any) {
	Trace(fmt.Sprintf(format, a...))
}

func Debugf(format string, a ...any) {
	Debug(fmt.Sprintf(format, a...))
}

func Infof(format string, a ...any) {
	Info(fmt.Sprintf(format, a...))
}

func Warningf(format string, a ...any) {
	Warning(fmt.Sprintf(format, a...))
}

func Errorf(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
}

func Fatalf(format string, a ...any) {
	Fatal(fmt.Sprintf(format, a...))
}
