package handler

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Hook 前后置方法
	Hook func(*gin.Context) []zap.Field

	// Options 中间件配置参数
	Options struct {
		RequestID         bool
		RequestBodyLimit  int
		RequestQueryLimit int
		ResponseBodyLimit int
		Before            Hook
		After             Hook
	}
)

// Logger gin logger
func Logger(logger *zap.Logger, opts Options) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if opts.Before == nil {
			opts.Before = defaultHook
		}
		if opts.After == nil {
			opts.After = defaultHook
		}

		start := time.Now()
		requestID, _ := shortid.Generate()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		userAgent := ctx.Request.UserAgent()
		clientIP := ctx.ClientIP()
		body := readBody(ctx)

		if opts.RequestQueryLimit < len(query) {
			query = query[:opts.RequestQueryLimit]
		}
		if opts.RequestBodyLimit < len(body) {
			body = body[:opts.RequestBodyLimit]
		}

		ctx.Set("requestId", requestID)

		logger.Info(uid(">>>", ctx), append([]zap.Field{
			zap.String("requestId", requestID),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.String("query", query),
			zap.String("body", string(body)),
		}, opts.Before(ctx)...)...)

		bodyWriter := newBodyWriter(ctx)
		ctx.Writer = bodyWriter

		ctx.Next()

		end := time.Now()
		latency := end.Sub(start)
		resp := bodyWriter.body.String()

		if opts.ResponseBodyLimit < len(resp) {
			resp = resp[:opts.ResponseBodyLimit]
		}

		for _, err := range ctx.Errors {
			logger.Info(uid("xxx", ctx),
				zap.String("requestId", requestID),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("ip", clientIP),
				zap.String("user-agent", userAgent),
				zap.String("query", query),
				zap.String("body", string(body)),
				zap.Int("status", ctx.Writer.Status()),
				zap.String("err", err.Error()),
			)
		}

		logger.Info(uid("<<<", ctx),
			append([]zap.Field{
				zap.String("requestId", requestID),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("ip", clientIP),
				zap.String("user-agent", userAgent),
				zap.String("query", query),
				zap.String("body", string(body)),
				zap.Int("status", ctx.Writer.Status()),
				zap.Duration("latency", latency),
				zap.String("resp", resp),
			}, opts.After(ctx)...)...)
	}
}

// NewProductionLogger 新建一个logger
func NewProductionLogger() (*zap.Logger, error) {
	cfg := newProductionConfig()

	return cfg.Build()
}

func newProductionConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    newProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func newProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "logger",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}
}

func defaultHook(*gin.Context) []zap.Field {
	return []zap.Field{}
}

func uid(prefix string, ctx *gin.Context) string {
	var b strings.Builder

	b.WriteString(prefix)
	b.WriteString("[")
	b.WriteString(ctx.GetString("uid"))
	b.WriteString("]")

	return b.String()
}
