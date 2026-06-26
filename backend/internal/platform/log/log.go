package log

import (
	"os"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger creates a JSON-formatted zap.SugaredLogger.
// The minimum log level is Info; stack traces are attached to all logs at that level or above.
// Regular logs go to stdout; error output goes to both stdout and stderr.
// appName is set as the logger name (the "logger" field in each log entry).
func NewZapLogger(appName string) (*zap.SugaredLogger, error) {
	config := zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout", "stderr"},
		DisableStacktrace: false,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	logger, err := config.Build(zap.AddCaller(), zap.AddStacktrace(zap.InfoLevel))
	if err != nil {
		return nil, err
	}

	return logger.Sugar().Named(appName), nil
}

// NewZerologLogger creates a JSON-formatted zerolog.Logger.
// Note: this sets zerolog's global log level to Info as a process-wide side effect.
// The following fields are automatically attached to every log entry:
//   - timestamp (ISO 8601)
//   - caller (source file and line number)
//   - appName
//   - environment
func NewZerologLogger(appName, env string) zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str("appName", appName).
		Str("environment", env).
		Logger()
}
