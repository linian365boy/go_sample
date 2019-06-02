package zapLog

import (
	"go.uber.org/zap/zapcore"
	"time"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"go.uber.org/zap"
)

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey: "T",
		LevelKey: "L",
		NameKey: "N",
		CallerKey: "C",
		MessageKey: "M",
		StacktraceKey: "S",
		LineEnding: zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime: TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}


func InitLog() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),
			w),
		zap.DebugLevel,
	)
	logger := zap.New(core, zap.AddCaller())
	logger.Info("zap log init info...")
}

