package logging

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugaredLogger *zap.SugaredLogger
var zapLogger *zap.Logger

// Init configures the global logger once per process.
func Init(appDebug bool) *zap.SugaredLogger {
	if sugaredLogger != nil {
		return sugaredLogger
	}

	logLevel := zapcore.InfoLevel
	if appDebug {
		logLevel = zapcore.DebugLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "lvl",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	outputPaths := []string{"stdout"}
	logDir := filepath.Join("storage", "logs")
	if err := os.MkdirAll(logDir, 0o755); err == nil {
		outputPaths = append(outputPaths, filepath.Join(logDir, "app.log"))
	} else {
		fmt.Fprintf(os.Stderr, "logging: failed to create log dir %s: %v\n", logDir, err)
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(logLevel),
		Development:      appDebug,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := cfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		panic(fmt.Sprintf("failed to build logger: %v", err))
	}

	zapLogger = logger
	zap.RedirectStdLog(zapLogger)
	zap.ReplaceGlobals(zapLogger)
	sugaredLogger = zapLogger.Sugar()
	return sugaredLogger
}

// Logger returns the initialized sugared logger.
func Logger() *zap.SugaredLogger {
	return sugaredLogger
}

// Sync flushes buffered log entries.
func Sync() {
	if zapLogger != nil {
		_ = zapLogger.Sync()
	}
}
