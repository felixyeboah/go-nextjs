package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}

// ZapLogger implements the Logger interface using zap
type ZapLogger struct {
	logger *zap.SugaredLogger
}

// NewZapLogger creates a new ZapLogger
func NewZapLogger(level string, development bool) (*ZapLogger, error) {
	var zapLevel zapcore.Level
	err := zapLevel.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	var config zap.Config
	if development {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &ZapLogger{
		logger: logger.Sugar(),
	}, nil
}

// Debug logs a debug message
func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

// Info logs an info message
func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

// Warn logs a warning message
func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

// Error logs an error message
func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message and exits
func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

// With returns a logger with the specified key-value pairs
func (l *ZapLogger) With(keysAndValues ...interface{}) Logger {
	return &ZapLogger{
		logger: l.logger.With(keysAndValues...),
	}
}

// NewFileLogger creates a new logger that writes to a file
func NewFileLogger(level string, development bool, filePath string) (*ZapLogger, error) {
	var zapLevel zapcore.Level
	err := zapLevel.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	// Create log file
	// #nosec G304 - This is a controlled file path for logging
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
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
	}

	if development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(io.MultiWriter(os.Stdout, file)),
		zap.NewAtomicLevelAt(zapLevel),
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		logger: logger.Sugar(),
	}, nil
}

// NewConsoleLogger creates a new logger that writes to the console
func NewConsoleLogger(level string, development bool) (*ZapLogger, error) {
	var zapLevel zapcore.Level
	err := zapLevel.UnmarshalText([]byte(level))
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
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
	}

	if development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create core
	var encoder zapcore.Encoder
	if development {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zapLevel),
	)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		logger: logger.Sugar(),
	}, nil
}

// DefaultLogger returns a default logger
func DefaultLogger() Logger {
	logger, err := NewConsoleLogger("info", true)
	if err != nil {
		// If we can't create a logger, use a simple fallback
		return &fallbackLogger{}
	}
	return logger
}

// fallbackLogger is a simple logger used as a fallback
type fallbackLogger struct{}

func (l *fallbackLogger) log(level, msg string, keysAndValues ...interface{}) {
	fmt.Printf("[%s] %s | %s", time.Now().Format(time.RFC3339), level, msg)
	if len(keysAndValues) > 0 {
		fmt.Print(" |")
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				fmt.Printf(" %v=%v", keysAndValues[i], keysAndValues[i+1])
			} else {
				fmt.Printf(" %v=?", keysAndValues[i])
			}
		}
	}
	fmt.Println()
}

func (l *fallbackLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.log("DEBUG", msg, keysAndValues...)
}

func (l *fallbackLogger) Info(msg string, keysAndValues ...interface{}) {
	l.log("INFO", msg, keysAndValues...)
}

func (l *fallbackLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.log("WARN", msg, keysAndValues...)
}

func (l *fallbackLogger) Error(msg string, keysAndValues ...interface{}) {
	l.log("ERROR", msg, keysAndValues...)
}

func (l *fallbackLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.log("FATAL", msg, keysAndValues...)
	os.Exit(1)
}

func (l *fallbackLogger) With(keysAndValues ...interface{}) Logger {
	// The fallback logger doesn't support With, so just return itself
	return l
}
