package util

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"smlm_automation/pkg/config" // Adjust import path
)

// Logger is the application's central logger.
var Logger *logrus.Logger
var once sync.Once

// InitLogger initializes the global logger based on the application configuration.
func InitLogger() {
	once.Do(func() {
		Logger = logrus.New()
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		// Configure screen output
		screenLevel, err := logrus.ParseLevel(config.AppConfig.Log.ScreenLevel)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid screen log level '%s', defaulting to info: %v\n", config.AppConfig.Log.ScreenLevel, err)
			screenLevel = logrus.InfoLevel
		}
		Logger.SetOutput(os.Stdout) // Default output to screen
		Logger.SetLevel(screenLevel)

		// Configure file output if a file path is provided
		if config.AppConfig.Log.FilePath != "" {
			fileLevel, err := logrus.ParseLevel(config.AppConfig.Log.FileLevel)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid file log level '%s', defaulting to debug: %v\n", config.AppConfig.Log.FileLevel, err)
				fileLevel = logrus.DebugLevel
			}

			// Ensure the directory exists
			logDir := ""
			lastSlash := strings.LastIndex(config.AppConfig.Log.FilePath, "/")
			if lastSlash != -1 {
				logDir = config.AppConfig.Log.FilePath[:lastSlash]
				if err := os.MkdirAll(logDir, 0755); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create log directory '%s': %v\n", logDir, err)
					return // Don't proceed with file logging if directory creation fails
				}
			}

			logFile, err := os.OpenFile(config.AppConfig.Log.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to open log file '%s': %v\n", config.AppConfig.Log.FilePath, err)
				return // Don't proceed with file logging
			}

			// Create a multi-writer to send logs to both screen and file
			mw := io.MultiWriter(os.Stdout, logFile)
			Logger.SetOutput(mw)

			// Set the highest log level among screen and file for the logger itself.
			// Logrus only supports one global level. We'll filter the output for
			// screen and file individually by setting hooks or by creating
			// separate loggers if more granular control is needed.
			// For simplicity, we'll use the higher of the two levels.
			if fileLevel < screenLevel {
				Logger.SetLevel(fileLevel)
			} else {
				Logger.SetLevel(screenLevel)
			}

			// To achieve separate log levels for screen and file with a single Logrus instance,
			// you'd typically use a hook.
			// This example demonstrates how to set up the multi-writer and a single level.
			// For truly separate levels, you'd need custom hooks for each output.
			// Alternatively, you could have two separate logrus instances, one for screen and one for file.
			// For this example, we'll demonstrate a simplified approach.

			// If you truly need separate levels, you'd do something like this (more complex):
			// Logger.AddHook(&ScreenLevelHook{Levels: []logrus.Level{logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}})
			// Logger.AddHook(&FileLevelHook{Levels: []logrus.Level{logrus.DebugLevel, logrus.InfoLevel, logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}})
		}
	})
}

// Logrus levels: Trace, Debug, Info, Warn, Error, Fatal, Panic
// We'll expose the common ones.

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

// Fprintf variations for structured logging
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
