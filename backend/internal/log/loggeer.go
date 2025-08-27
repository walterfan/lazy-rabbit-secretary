package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapLogger *zap.Logger
)

func InitLogger() error {
	// Define encoder configuration for human-readable + machine-friendly output
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create cores: one for stdout, one for file
	consoleEncoder := zapcore.NewJSONEncoder(encoderCfg)
	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)

	// Open the log file
	logFile, err := os.OpenFile("async-llm-agent.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zap.InfoLevel),
	)

	zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))
	defer zapLogger.Sync() // flushes buffer, if any

	return nil
}

func GetLogger() *zap.SugaredLogger {
	return zapLogger.Sugar()
}
