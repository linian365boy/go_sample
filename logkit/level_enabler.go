package logkit

import (
	"go.uber.org/zap/zapcore"
)

type LevelEnabler struct {
	out_level *zapcore.Level
	min_level zapcore.Level
	max_level zapcore.Level
	log_level zapcore.Level
}

func (l *LevelEnabler) Enabled(level zapcore.Level) bool {
	//fmt.Println(level, *l.out_level <= l.min_level, l.min_level <= level && level <= l.max_level, *l.out_level, l.log_level)
	if *l.out_level <= l.min_level && // enable loggger
		l.min_level <= level && level <= l.max_level { // log to right logger
		return true
	}
	return false
}

func NewLevelEnabler(out_level *zapcore.Level, logger_level zapcore.Level) *LevelEnabler {
	var (
		min_level zapcore.Level
		max_level zapcore.Level
	)

	switch logger_level {
	case zapcore.DebugLevel:
		min_level = zapcore.DebugLevel
		max_level = zapcore.DebugLevel
	case zapcore.InfoLevel:
		min_level = zapcore.InfoLevel
		max_level = zapcore.WarnLevel
	case zapcore.ErrorLevel:
		min_level = zapcore.ErrorLevel
		max_level = zapcore.FatalLevel
	}

	return &LevelEnabler{
		out_level: out_level,
		min_level: min_level,
		max_level: max_level,
		log_level: logger_level,
	}
}
