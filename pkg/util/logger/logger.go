package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	model "smlm_automation/pkg/models/inputfile"
	"smlm_automation/pkg/util/constants"
	ri "smlm_automation/pkg/util/readconfig"
	"strings"
	"sync"
	// "smlm_automation/pkg/config" // Adjust import path
)

// Logger is the application's central logger.
var Logger *logrus.Logger
var once sync.Once

// InitLogger initializes the global logger based on the application configuration.
func InitLogger() error {
	var gc model.ConfigGeneral
	err := ri.ReadConfig(constants.GeneralConfigFile, &gc)
	if err != nil {
		return err
	}
	once.Do(func() {
		Logger = logrus.New()
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		// Configure screen output
		screenLevel, err := logrus.ParseLevel(gc.General.Log.ScreenLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid screen log level '%s', defaulting to info: %v\n", gc.General.Log.ScreenLevel, err)
			screenLevel = logrus.InfoLevel
		}
		Logger.SetOutput(os.Stdout) // Default output to screen
		Logger.SetLevel(screenLevel)
		// Configure file output if a file path is provided
		if gc.General.Log.FilePath != "" {
			fileLevel, err := logrus.ParseLevel(gc.General.Log.FileLevel)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid file log level '%s', defaulting to debug: %v\n", gc.General.Log.FileLevel, err)
				fileLevel = logrus.DebugLevel
			}
			logDir := ""
			lastSlash := strings.LastIndex(gc.General.Log.FilePath, "/")
			if lastSlash != -1 {
				logDir = gc.General.Log.FilePath[:lastSlash]
				if err := os.MkdirAll(logDir, 0755); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create log directory '%s': %v\n", logDir, err)
					return
				}
			}
			logFile, err := os.OpenFile(gc.General.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open log file '%s': %v\n", gc.General.Log.FilePath, err)
				return
			}
			mw := io.MultiWriter(os.Stdout, logFile)
			Logger.SetOutput(mw)
			if fileLevel < screenLevel {
				Logger.SetLevel(fileLevel)
			} else {
				Logger.SetLevel(screenLevel)
			}
		}
	})
	return nil
}

func Debug(args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.DebugLevel) {
		Logger.Debug(args...)
	}
}

func Info(args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.InfoLevel) {
		Logger.Info(args...)
	}
}

func Warn(args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.WarnLevel) {
		Logger.Warn(args...)
	}
}

func Error(args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.ErrorLevel) {
		Logger.Error(args...)
	}
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Debugf(format string, args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.DebugLevel) {
		Logger.Debugf(format, args...)
	}
}

func Infof(format string, args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.InfoLevel) {
		Logger.Infof(format, args...)
	}
}

func Warnf(format string, args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.WarnLevel) {
		Logger.Warnf(format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	if Logger.IsLevelEnabled(logrus.ErrorLevel) {
		Logger.Errorf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// NOTE: For truly separate log levels for screen and file,
// you would typically need to implement `logrus.Hook` for each output.
// A simpler approach for the scope of this example is to set the
// *global* log level of the Logrus instance to the most permissive
// of the two (e.g., if screen is INFO and file is DEBUG, the global
// level should be DEBUG). Then, within the hook or before writing
// to the specific output, you would filter based on the desired level
// for that output.
