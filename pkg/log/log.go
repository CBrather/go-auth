package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Initialize(logLevel string) error {
	parsedLevel, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		parsedLevel = zap.InfoLevel
	}

	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(parsedLevel)
	setEncoderConfig(&config)

	logger, err := config.Build()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(logger)

	return nil
}

func setEncoderConfig(cfg *zap.Config) {
	cfg.EncoderConfig.EncodeTime = timeEncoder
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.StacktraceKey = "callstack"
	cfg.EncoderConfig.TimeKey = "timestamp"
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format(time.RFC3339Nano))
}
