package logger

import (
	"cvs/internal/config"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for logging.
type Logger interface {
	// InitLogger initializes the logger with settings from the configuration.
	InitLogger()

	// Debug logs a debug message.
	Debug(args ...interface{})

	// Debugf logs a formatted debug message.
	Debugf(template string, args ...interface{})

	// Info logs an informational message.
	Info(args ...interface{})

	// Infof logs a formatted informational message.
	Infof(template string, args ...interface{})

	// Warn logs a warning message.
	Warn(args ...interface{})

	// Warnf logs a formatted warning message.
	Warnf(template string, args ...interface{})

	// Error logs an error message.
	Error(args ...interface{})

	// Errorf logs a formatted error message.
	Errorf(template string, args ...interface{})

	// DPanic logs a panic message in development mode.
	DPanic(args ...interface{})

	// DPanicf logs a formatted panic message in development mode.
	DPanicf(template string, args ...interface{})

	// Panic logs a panic message and calls panic().
	Panic(args ...interface{})

	// Panicf logs a formatted panic message and calls panic().
	Panicf(template string, args ...interface{})

	// Fatal logs a fatal message and exits the program.
	Fatal(args ...interface{})

	// Fatalf logs a formatted fatal message and exits the program.
	Fatalf(template string, args ...interface{})
}

// apiLogger is an implementation of Logger using the Zap library.
type apiLogger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
}

// NewApiLogger creates a new instance of apiLogger with the given configuration.
func NewApiLogger(cfg *config.Config) *apiLogger {
	return &apiLogger{cfg: cfg}
}

// customTimeEncoder configures the time format for logs in AM/PM format.
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 03:04:05 PM"))
}

// loggerLevelMap maps configuration logger levels to Zap logger levels.
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// getLoggerLevel returns the logging level based on the configuration.
func (l *apiLogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exist := loggerLevelMap[cfg.Logger.Level]
	if !exist {
		return zapcore.DebugLevel // Default level
	}
	return level
}

// InitLogger initializes the logger with settings from the configuration.
func (l *apiLogger) InitLogger() {
	logLevel := l.getLoggerLevel(l.cfg)

	logWriter := zapcore.AddSync(os.Stderr)

	var encoderCfg zapcore.EncoderConfig
	if l.cfg.ServerMode == "dev" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.EncodeTime = customTimeEncoder
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	var encoder zapcore.Encoder
	if l.cfg.Logger.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.sugarLogger = logger.Sugar()
}

// Logger methods implementations
func (l *apiLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *apiLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *apiLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *apiLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *apiLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *apiLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *apiLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *apiLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *apiLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *apiLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *apiLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *apiLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *apiLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *apiLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}
