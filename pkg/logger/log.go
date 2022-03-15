package logger

import (
	"context"
	"os"
	"sort"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once
var stdout Logger

const (
	defaultLevel = zapcore.InfoLevel
)

type Field = zapcore.Field
type Level = zapcore.Level
type Encoder = zapcore.Encoder
type EncoderConfig = zapcore.EncoderConfig

// logger interface
type Logger interface {
	Debug(ctx context.Context, message string, fields ...Field)
	Info(ctx context.Context, message string, fields ...Field)
	Warn(ctx context.Context, message string, fields ...Field)
	Error(ctx context.Context, message string, fields ...Field)
	Fatal(ctx context.Context, message string, fields ...Field)
}

// logger init
func InitLogger(name, application, environment, std string) {
	once.Do(func() {
		switch name {
		case "zap":
			stdout = NewZapLogger(application, environment, std)
		default:
			stdout = NewZapLogger(application, environment, std)
		}
	})
}

func switchLevel(environment string) Level {
	switch environment {
	case "product":
		return zapcore.InfoLevel
	case "test":
		return zapcore.DebugLevel
	case "dev":
		return zapcore.DebugLevel
	default:
		return defaultLevel
	}
}

func applyEncoder(std string, enc EncoderConfig) Encoder {
	switch std {
	case "json":
		return zapcore.NewJSONEncoder(enc)
	case "console":
		return zapcore.NewConsoleEncoder(enc)
	default:
		return zapcore.NewJSONEncoder(enc)
	}
}

type ZapLogger struct {
	Logger      *zap.Logger
	Env         string
	App         string
	defaultCode int
}

// new logger base zap
func NewZapLogger(application, environment, std string) Logger {
	// level
	level := switchLevel(environment)
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	enc := zapcore.EncoderConfig{
		TimeKey:  "when",
		LevelKey: "level",
		// NameKey:   "logger",
		CallerKey:  "caller",
		MessageKey: "message",
		// FunctionKey:    "function",
		StacktraceKey:  "traceback",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// set writer
	zapWriterSync := zapcore.AddSync(os.Stdout)
	// set encoder
	encoder := applyEncoder(std, enc)
	// apply param
	zapCore := zapcore.NewCore(
		encoder,
		zapWriterSync,
		atomicLevel,
	)
	// new logger
	logger := zap.New(zapCore, zap.AddCaller(), zap.AddCallerSkip(1))
	defer logger.Sync()
	return &ZapLogger{
		Logger:      logger,
		Env:         environment,
		App:         application,
		defaultCode: 10000,
	}
}

// write message
func (c *ZapLogger) Writer(ctx context.Context, level string, message string, fields ...Field) {
	fields = c.build(ctx, fields...)
	switch level {
	case "DEBUG":
		c.Logger.Debug(message, fields...)
	case "INFO":
		c.Logger.Info(message, fields...)
	case "WARN":
		c.Logger.Warn(message, fields...)
	case "ERROR":
		c.Logger.Error(message, fields...)
	case "FATAL":
		c.Logger.Fatal(message, fields...)
	default:
		c.Logger.Info(message, fields...)
	}
}

// whether the element is in the array
func containString(target string, raw []string) bool {
	sort.Strings(raw)
	index := sort.SearchStrings(raw, target)
	if index < len(raw) && raw[index] == target {
		return true
	}
	return false
}

/*
function of zap logger
*/
func (c *ZapLogger) Debug(ctx context.Context, message string, fields ...Field) {
	c.Writer(ctx, "Debug", message, fields...)
}

func (c *ZapLogger) Info(ctx context.Context, message string, fields ...Field) {
	c.Writer(ctx, "INFO", message, fields...)
}

func (c *ZapLogger) Warn(ctx context.Context, message string, fields ...Field) {
	c.Writer(ctx, "WARN", message, fields...)
}

func (c *ZapLogger) Error(ctx context.Context, message string, fields ...Field) {
	c.Writer(ctx, "ERROR", message, fields...)
}

func (c *ZapLogger) Fatal(ctx context.Context, message string, fields ...Field) {
	c.Writer(ctx, "FATAL", message, fields...)
}

// build message body
func (c *ZapLogger) build(ctx context.Context, fields ...zapcore.Field) []zapcore.Field {
	keys := []string{"trace", "remark", "traceback", "d-lever", "duration", "param", "result", "code", "channel", "remote"}
	message := map[string]zapcore.Field{
		"trace":       zap.String("trace", ctx.Value("trace").(string)),
		"env":         zap.String("env", c.Env),
		"application": zap.String("application", c.App),
		"channel":     zap.String("channel", ""),
		"remark":      zap.String("remark", ""),
		"code":        zap.Int("code", c.defaultCode),
		"traceback":   zap.String("traceback", ""),
		"duration":    zap.Int64("duration", 0),
		"d-lever":     zap.String("d-lever", "info"),
		"param":       zap.Any("param", map[string]interface{}{}),
		"result":      zap.Any("result", map[string]interface{}{}),
		"remote":      zap.Any("remote", map[string]interface{}{}),
	}
	for i := 0; i < len(fields); i++ {
		key := fields[i].Key
		if containString(key, keys) {
			message[key] = fields[i]
		}
	}
	result := []zapcore.Field{}
	for _, field := range message {
		result = append(result, field)
	}
	return result
}

/*
for caller
*/
func Debug(ctx context.Context, message string, fields ...Field) {
	stdout.Debug(ctx, message, fields...)
}

func Info(ctx context.Context, message string, fields ...Field) {
	stdout.Info(ctx, message, fields...)
}

func Warn(ctx context.Context, message string, fields ...Field) {
	stdout.Warn(ctx, message, fields...)
}

func Error(ctx context.Context, message string, fields ...Field) {
	stdout.Error(ctx, message, fields...)
}

func Fatal(ctx context.Context, message string, fields ...Field) {
	stdout.Fatal(ctx, message, fields...)
}
