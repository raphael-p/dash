package logger

import (
	"cmp"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	MaxFileCount                  int    = 10
	MaxFileBytes, MaxMessageBytes int64  = 10 * 1024 * 1024, 10 * 1024
	DirName, Extension, Prefix    string = "logs", ".log", "server"
)

type MyLogger struct {
	stdOutLogger *log.Logger
	fileLogger   *log.Logger
	logLevel     logLevel
	cumBytes     int64
	dir          string
}

type logLevel uint8

const (
	logLevelTrace logLevel = iota
	logLevelDebug
	logLevelInfo
	logLevelWarning
	logLevelError
)

func newLogger(out io.Writer) *log.Logger {
	return log.New(out, "", 0)
}

var Logger *MyLogger = &MyLogger{nil, nil, logLevel(0), 0, "."}

func setDir(dir string) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(fmt.Sprintf("failed to create log directory: %s", err))
	}
	Logger.dir = dir
}

func setFileLogger() {
	parsedNames, err := listNames()
	if err != nil {
		panic(fmt.Sprintf("failed to list log files: %s", err))
	}
	file, err := openLogFile(parsedNames)
	if err != nil {
		panic(fmt.Sprintf("failed to open log file: %s", err))
	}

	Logger.fileLogger = newLogger(file)

	go func() {
		if err := prune(); err != nil {
			Error(fmt.Sprintf("failed to prune log files: %s", err))
		}
	}()
}

type parsedName struct {
	name    string
	dateStr string
	count   uint8
}

// Expected format: <prefix>-20060102-1<Extension>
func parseName(name string) (parsedName, bool) {
	parsedName := parsedName{name: name}
	if !strings.HasPrefix(name, Prefix) || !strings.HasSuffix(name, Extension) {
		return parsedName, false
	}
	splitN := strings.SplitN(name, "-", 3)
	if len(splitN) != 3 || splitN[0] != Prefix {
		return parsedName, false
	}
	if _, err := time.Parse("20060102", splitN[1]); err != nil {
		return parsedName, false
	} else {
		parsedName.dateStr = splitN[1]
	}
	if _index, err := strconv.ParseUint(strings.TrimSuffix(splitN[2], Extension), 10, 8); err != nil {
		return parsedName, false
	} else {
		parsedName.count = uint8(_index)
	}
	return parsedName, true
}

func listNames() ([]parsedName, error) {
	entries, err := os.ReadDir(Logger.dir)
	if err != nil {
		return []parsedName{}, err
	}

	var names []parsedName
	for _, entry := range entries {
		if entry.Type().IsRegular() {
			if name, ok := parseName(entry.Name()); ok {
				names = append(names, name)
			}
		}
	}
	slices.SortFunc(names, func(a, b parsedName) int {
		return cmp.Or(
			cmp.Compare(a.dateStr, b.dateStr),
			cmp.Compare(a.count, b.count),
		)
	})
	return names, nil
}

func getNextLogPath(parsedNames []parsedName) string {
	dateStr := time.Now().UTC().Format("20060102")
	var latestName parsedName
	for _, p := range parsedNames {
		if p.dateStr == dateStr {
			if p.count > latestName.count {
				latestName = p
			}
		}
	}

	if latestName.name != "" {
		latestPath := filepath.Join(Logger.dir, latestName.name)
		stats, err := os.Stat(latestPath)
		if err == nil && stats.Size() < MaxFileBytes {
			return latestPath
		}
	}

	return filepath.Join(Logger.dir, fmt.Sprintf("%s-%s-%d%s", Prefix, dateStr, latestName.count+1, Extension))
}

func openLogFile(parsedNames []parsedName) (*os.File, error) {
	path := getNextLogPath(parsedNames)
	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logFile, err = os.Create(path)
		if err != nil {
			return nil, err
		}
	}
	return logFile, nil
}

func prune() error {
	parsedNames, err := listNames()
	if err != nil {
		return err
	}
	for len(parsedNames) > MaxFileCount {
		path := filepath.Join(Logger.dir, parsedNames[0].name)
		if err := os.Remove(path); err != nil {
			return err
		}
		Debugf("pruned log file: %s", path)
		parsedNames = parsedNames[1:]
	}
	return nil
}

// Create a logger which manages its logs files and logs to the console.
func Create(workingDir string, hasConsoleLog bool) {
	Logger.logLevel = logLevel(0)
	if hasConsoleLog {
		Logger.stdOutLogger = newLogger(os.Stdout)
	}
	setDir(filepath.Join(workingDir, DirName))
	setFileLogger()
	Info("file logger initialised")
	Debugf("logger directory: %s", Logger.dir)
}

// Close the log file
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
	time := time.Now().UTC().Format("2006-01-02 15:04:05.000")
	if int64(len(message)) > MaxMessageBytes {
		message = message[:MaxMessageBytes]
	}
	if Logger.stdOutLogger != nil {
		Logger.stdOutLogger.Printf("%s %s[%s]%s %s", time, ansiColour, level, reset, message)
	}
	if Logger.fileLogger != nil {
		if Logger.cumBytes += int64(len(message)); Logger.cumBytes > MaxFileBytes {
			rollover()
		}
		Logger.fileLogger.Printf("%s [%s] %s", time, level, message)
	}
}

func rollover() {
	if file, ok := Logger.fileLogger.Writer().(*os.File); ok {
		if stats, err := file.Stat(); err == nil && stats.Size() > MaxFileBytes {
			Logger.cumBytes = 0
			Info(fmt.Sprintf(
				"reached maximum log file size (%d bytes), rolling over",
				MaxFileBytes,
			)) // can cause infinite loop if MAX_FILE_BYTES is too low
			Close()
			setFileLogger()
			Debug("new log file opened")
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
