package logkit

import (
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logwrapper *logWrapper
	level      zapcore.Level
	once       sync.Once
)

type logWrapper struct {
	log   *zap.Logger
	sugar *zap.SugaredLogger
}

func NewLogWrapper(log *zap.Logger) *logWrapper {
	return &logWrapper{log: log, sugar: log.Sugar()}
}

func (l *logWrapper) Debug(args ...interface{}) {
	if !l.log.Core().Enabled(zapcore.DebugLevel) {
		return
	}
	l.sugar.Debug(args...)
}

func (l *logWrapper) Debugf(format string, args ...interface{}) {
	if !l.log.Core().Enabled(zapcore.DebugLevel) {
		return
	}
	l.sugar.Debugf(format, args...)
}

func (l *logWrapper) Info(args ...interface{}) {
	if !l.log.Core().Enabled(zapcore.InfoLevel) {
		return
	}
	l.sugar.Info(args...)
}

func (l *logWrapper) Infoln(args ...interface{}) {
	if !l.log.Core().Enabled(zapcore.InfoLevel) {
		return
	}
	l.sugar.Info(args...)
}

func (l *logWrapper) Infof(format string, args ...interface{}) {
	if !l.log.Core().Enabled(zapcore.InfoLevel) {
		return
	}
	l.sugar.Infof(format, args...)
}

func (l *logWrapper) Warning(args ...interface{}) { l.sugar.Warn(args...) }

func (l *logWrapper) Warningln(args ...interface{}) { l.sugar.Warn(args...) }

func (l *logWrapper) Warningf(format string, args ...interface{}) { l.sugar.Warnf(format, args...) }

func (l *logWrapper) Error(args ...interface{}) { l.sugar.Error(args...) }

func (l *logWrapper) Errorln(args ...interface{}) { l.sugar.Error(args...) }

func (l *logWrapper) Errorf(format string, args ...interface{}) { l.sugar.Errorf(format, args...) }

func (l *logWrapper) Fatal(args ...interface{}) { l.sugar.Fatal(args...) }

func (l *logWrapper) Fatalln(args ...interface{}) { l.sugar.Fatal(args...) }

func (l *logWrapper) Fatalf(format string, args ...interface{}) { l.sugar.Fatalf(format, args...) }

func (l *logWrapper) Printf(format string, args ...interface{}) { l.sugar.Infof(format, args...) }

func (l *logWrapper) V(v int) bool {
	if v <= 0 {
		return !l.log.Core().Enabled(zapcore.DebugLevel)
	}
	return true
}

func (l logWrapper) Log(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l logWrapper) Logf(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func Init(opts ...Option) error {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	err := level.UnmarshalText([]byte(options.level))
	if err != nil {
		return err
	}

	infoFileHandler := zapcore.AddSync(&lumberjack.Logger{
		Filename:   options.path + "/info.log",
		MaxSize:    options.maxSize,
		MaxBackups: options.maxBackups,
		MaxAge:     options.maxAge,
	})
	errFileHandler := zapcore.AddSync(&lumberjack.Logger{
		Filename:   options.path + "/error.log",
		MaxSize:    options.maxSize,
		MaxBackups: options.maxBackups,
		MaxAge:     options.maxAge,
	})
	debugFileHandler := zapcore.AddSync(&lumberjack.Logger{
		Filename:   options.path + "/debug.log",
		MaxSize:    options.maxSize,
		MaxBackups: options.maxBackups,
		MaxAge:     options.maxAge,
	})

	liveEncoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	devEncoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	zapCores := []zapcore.Core{
		zapcore.NewCore(zapcore.NewJSONEncoder(liveEncoder), errFileHandler, NewLevelEnabler(&level, zapcore.ErrorLevel)),
		zapcore.NewCore(zapcore.NewJSONEncoder(liveEncoder), infoFileHandler, NewLevelEnabler(&level, zapcore.InfoLevel)),
		zapcore.NewCore(zapcore.NewJSONEncoder(devEncoder), debugFileHandler, NewLevelEnabler(&level, zapcore.DebugLevel)),
	}

	if options.enableConsole {
		consoleHandler := zapcore.Lock(os.Stdout)
		zapCores = append(zapCores, zapcore.NewCore(zapcore.NewConsoleEncoder(devEncoder), consoleHandler, zapcore.DebugLevel))
	}
	// create options with priority for our opts
	defaultOptions := []zap.Option{}
	if options.enableCaller {
		defaultOptions = append(
			defaultOptions,
			zap.AddCaller(),
			//zap.AddStacktrace(level),
			zap.AddCallerSkip(1),
		)
	}

	core := zapcore.NewTee(
		zapCores...,
	)

	logger := zap.New(core, defaultOptions...)

	logwrapper = NewLogWrapper(logger)

	return err
}

func GetZapLogger() *zap.Logger {
	return GetLogger().log
}

func GetLogger() *logWrapper {
	if logwrapper == nil {
		Init()
	}
	return logwrapper
}

func Sync() error {
	return GetLogger().log.Sync()
}

func SetLevel(l string) error {
	// init logger
	GetLogger()

	err := level.UnmarshalText([]byte(l))
	if err != nil {
		return err
	}
	return nil
}
