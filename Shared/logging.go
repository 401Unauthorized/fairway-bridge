package Shared

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

// NewLogger creates a new logger instance with the specified configuration.
func NewLogger(config Config) (*zap.Logger, *os.File, error) {
	file, err := os.OpenFile(config.Bridge.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Only write to the file if the launch monitor is "VIRTUAL"
	var writer io.Writer
	if config.LaunchMonitor.Name == "VIRTUAL" {
		writer = io.MultiWriter(file)
	} else {
		writer = io.MultiWriter(os.Stdout, file)
	}

	syncer := zapcore.AddSync(writer)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "msg"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var encoder zapcore.Encoder
	switch config.Bridge.LogType {
	case "JSON":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "CONSOLE":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return nil, nil, fmt.Errorf("invalid log type: %s", config.Bridge.LogType)
	}

	level := zapcore.InfoLevel
	core := zapcore.NewCore(encoder, syncer, level)
	return zap.New(core), file, nil
}
