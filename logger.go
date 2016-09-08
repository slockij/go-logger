package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

//Log level constants
const (
	LogFatal int = iota
	LogError
	LogWarning
	LogInfo
	LogDebug
	LogDebug2
)

type (
	//Logger - log writer
	Logger interface {
		LogLevel() int
		LogLevelString() string
		GetLogWriter() io.Writer
		Rotate()
		Output(calldepth int, s string) error
		Log(logLevel, depth int, args ...interface{})
		Fatal(args ...interface{})
		Error(args ...interface{})
		Warning(args ...interface{})
		Info(args ...interface{})
		Debug(args ...interface{})
		Debug2(args ...interface{})
	}

	logHandler struct {
		logLevel int
		w        *log.Logger
		f        io.Writer
		filename string
	}

	logWrapper struct {
		logLevel int
	}
)

var levels = []string{"FATAL", "ERROR", "WARNING", "INFO", "DEBUG", "DEBUG2"}

//FindLogLevel - find log level from string
func FindLogLevel(logLevel string) int {
	if logLevel == "DEBUG2" {
		return LogDebug2
	}
	if logLevel == "DEBUG" {
		return LogDebug
	}
	if logLevel == "INFO" {
		return LogInfo
	}
	if logLevel == "WARNING" || logLevel == "WARN" {
		return LogWarning
	}
	if logLevel == "ERROR" {
		return LogError
	}
	return LogFatal
}

//GetLogLevelString - find log level string name from log level
func GetLogLevelString(lvl int) string {
	return levels[lvl]
}

//NewLogger - create new logger
func NewLogger(logLevel int, file string) Logger {
	if logLevel > LogDebug {
		logLevel = LogDebug2
	}
	if logLevel < LogFatal {
		logLevel = LogFatal
	}
	f, err := newLog(file)
	if err != nil {
		panic("Cannot instantiate logger with file " + file)
	}
	w := log.New(f, "", log.Lshortfile|log.LUTC|log.LstdFlags)
	return &logHandler{logLevel, w, f, file}
}

//GetLogWriter - retrieve file that is used as log
func (l *logHandler) GetLogWriter() io.Writer {
	return l.f
}

//Rotate - reopen file
func (l *logHandler) Rotate() {
	if l.filename != "" && l.filename != "stdout" {
		f, err := newLog(l.filename)
		if err == nil {
			l.f = f
			l.w = log.New(f, "", log.Lshortfile|log.LUTC|log.LstdFlags)
			l.Output(0, "Rotated log file")
		} else {
			l.Output(0, fmt.Sprint("Could not rotate log file, error was:", err))
		}
	}
}

func newLog(file string) (io.Writer, error) {
	f := os.Stdout
	if file != "" && file != "stdout" {
		var err error
		f, err = os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

//LogLevel - returns log level
func (l *logHandler) LogLevel() int {
	return l.logLevel
}

//LogLevelString - returns log level as string
func (l *logHandler) LogLevelString() string {
	return levels[l.logLevel]
}

//Output - implements output from standard logger
func (l *logHandler) Output(calldepth int, s string) error {
	return l.w.Output(calldepth+1, s)
}

//Log - more complex logging method
func (l *logHandler) Log(logLevel, depth int, args ...interface{}) {
	if logLevel <= l.logLevel {
		args = append([]interface{}{levels[logLevel]}, args...)
		l.w.Output(depth, fmt.Sprintln(args...))
	}
}

//Fatal - log and panic
func (l *logHandler) Fatal(args ...interface{}) {
	l.Log(LogFatal, 3, args...)
	panic("Fatal error occured")
}

//Error - save log message on ERROR level
func (l *logHandler) Error(args ...interface{}) {
	l.Log(LogError, 3, args...)
}

//Warning - save log message on WARNING level
func (l *logHandler) Warning(args ...interface{}) {
	l.Log(LogWarning, 3, args...)
}

//Info - save log message on INFO level
func (l *logHandler) Info(args ...interface{}) {
	l.Log(LogInfo, 3, args...)
}

//Debug - save log message on DEBUG level
func (l *logHandler) Debug(args ...interface{}) {
	l.Log(LogDebug, 3, args...)
}

//Debug2 - save log message on DEBUG2 level (even more debug)
func (l *logHandler) Debug2(args ...interface{}) {
	l.Log(LogDebug2, 3, args...)
}

//WrapLogger - wraps standard log with interface
func WrapLogger(logLevel int) Logger {
	if logLevel > LogDebug {
		logLevel = LogDebug2
	}
	if logLevel < LogFatal {
		logLevel = LogFatal
	}
	return &logWrapper{logLevel}
}

//GetLogWriter - returns nil as there's no file known
func (l *logWrapper) GetLogWriter() io.Writer {
	return nil
}

//Rotate - reopen file
func (l *logWrapper) Rotate() {
	l.Info("Rotate called")
}

//LogLevel - returns log level
func (l *logWrapper) LogLevel() int {
	return l.logLevel
}

//LogLevelString - returns log level as string
func (l *logWrapper) LogLevelString() string {
	return levels[l.logLevel]
}

//Output - implements output from standard logger
func (l *logWrapper) Output(calldepth int, s string) error {
	return log.Output(calldepth+1, s)
}

//Log - more complex logging method
func (l *logWrapper) Log(logLevel, depth int, args ...interface{}) {
	if logLevel <= l.logLevel {
		args = append([]interface{}{levels[logLevel]}, args...)
		log.Output(depth, fmt.Sprintln(args...))
	}
}

//Fatal - log and panic
func (l *logWrapper) Fatal(args ...interface{}) {
	l.Log(LogFatal, 3, args...)
	panic("Fatal error occured")
}

//Error - save log message on ERROR level
func (l *logWrapper) Error(args ...interface{}) {
	l.Log(LogError, 3, args...)
}

//Warning - save log message on WARNING level
func (l *logWrapper) Warning(args ...interface{}) {
	l.Log(LogWarning, 3, args...)
}

//Info - save log message on INFO level
func (l *logWrapper) Info(args ...interface{}) {
	l.Log(LogInfo, 3, args...)
}

//Debug - save log message on DEBUG level
func (l *logWrapper) Debug(args ...interface{}) {
	l.Log(LogDebug, 3, args...)
}

//Debug2 - save log message on DEBUG2 level (even more debug)
func (l *logWrapper) Debug2(args ...interface{}) {
	l.Log(LogDebug2, 3, args...)
}
