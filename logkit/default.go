package logkit

import "go.uber.org/zap/zapcore"

// Fatal outputs a message at fatal level.
func Fatal(msg string, fields ...zapcore.Field) {
	GetZapLogger().Fatal(msg, fields...)
}

// Error outputs a message at error level.
func Error(msg string, fields ...zapcore.Field) {
	GetZapLogger().Error(msg, fields...)
}

// Info outputs a message at info level.
func Info(msg string, fields ...zapcore.Field) {
	GetZapLogger().Info(msg, fields...)
}

// Warn outputs a message at warn level.
func Warn(msg string, fields ...zapcore.Field) {
	GetZapLogger().Warn(msg, fields...)
}

// Debug outputs a message at debug level.
func Debug(msg string, fields ...zapcore.Field) {
	GetZapLogger().Debug(msg, fields...)
}

// Fatalf outputs a message at fatal level.
func Fatalf(fmt string, args ...interface{}) {
	GetLogger().sugar.Fatalf(fmt, args...)
}

// Errorf outputs a message at error level.
func Errorf(fmt string, args ...interface{}) {
	GetLogger().sugar.Errorf(fmt, args...)
}

// Infof outputs a message at info level.
func Infof(fmt string, args ...interface{}) {
	GetLogger().sugar.Infof(fmt, args...)
}

// Warnf outputs a message at warn level.
func Warnf(fmt string, args ...interface{}) {
	GetLogger().sugar.Warnf(fmt, args...)
}

// Debugf outputs a message at debug level.
func Debugf(fmt string, args ...interface{}) {
	GetLogger().sugar.Debugf(fmt, args...)
}
