package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger levels
const (
	DEBUG string = `DEBUG`
	INFO  string = `INFO`
	WARN  string = `WARN`
	ERROR string = `ERROR`
	FATAL string = `FATAL`
)

// NewLogger returns a logger with level and formatting set from func params.
// If the level is not set or is passed an empty string
// automatically sets itself on WARN level.
func NewLogger(level string, isJSON bool) *zap.SugaredLogger {
	config := zap.Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
	}

	if isJSON {
		config.Encoding = "json"
	} else {
		config.Encoding = "console"
	}

	switch level {
	case DEBUG:
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	case INFO:
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case WARN:
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case ERROR:
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case FATAL:
		config.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		log.Printf("Unknown value of level, setting up to WARN\n")
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	}

	lg, err := config.Build()
	if err != nil {
		log.Fatalf("error initializing logger: %v", err)
	}
	logger := lg.Sugar()
	return logger
}
