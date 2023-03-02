package zaplog

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type SugaredLogger struct {
	*zap.SugaredLogger
}

func TestInfo(t *testing.T) {
	type args struct {
		c   context.Context
		msg string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "test info",
			args: args{
				c:   context.Background(),
				msg: "This is an info log",
			},
		},
		{
			name: "Error  producton",
			args: args{
				c:   context.Background(),
				msg: "This is an info log",
			},
		},
		{
			name: "NEW CASE",
			args: args{
				c:   context.Background(),
				msg: "This is an info log",
			},
		},
	}
	for _, tt := range tests {
		os.Setenv("ENVIRONMENT_NAME", "production")

		_ = gomonkey.ApplyFunc(os.Getenv,
			func(s string) string {
				log.Println("--------- gomonkey call --------")
				return "production"
			})
		t.Run(tt.name, func(t *testing.T) {

			// patchZapProd := gomonkey.ApplyFunc(zap.NewProduction, func(options ...zap.Option) (*zap.Logger, error) {
			// 	return &zap.Logger{}, nil
			// })

			// defer patchZapProd.Reset()

			//patchDaos := gomonkey.ApplyFunc(InitLogger,
			//	func() *zap.SugaredLogger {
			//		var zapLogger *zap.Logger
			//		zapLogger, _ = zap.NewProduction()
			//		return zapLogger.Sugar()
			//	})
			//defer patchDaos.Reset()

			observedZapCore, observedLogs := observer.New(zap.InfoLevel)
			observedLogger := zap.New(observedZapCore).Sugar()

			_ = SetLogger(observedLogger)
			Info(tt.args.c, tt.args.msg)
			assert.Equal(t, 1, observedLogs.Len())
			log := observedLogs.All()[0]
			assert.Equal(t, fmt.Sprintf("\nRequest-ID: <nil>\n[%s]\n", tt.args.msg), log.Message)
			assert.Equal(t, zapcore.Level(0), log.Level)
		})
	}
}

func TestDebug(t *testing.T) {
	type args struct {
		c   context.Context
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test info",
			args: args{
				c:   context.Background(),
				msg: "This is a debug log",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			observedZapCore, observedLogs := observer.New(zap.DebugLevel)
			observedLogger := zap.New(observedZapCore).Sugar()
			_ = SetLogger(observedLogger)
			Debug(tt.args.c, tt.args.msg)
			assert.Equal(t, 1, observedLogs.Len())
			log := observedLogs.All()[0]
			assert.Equal(t, fmt.Sprintf("\nRequest-ID: <nil>\n[%s]\n", tt.args.msg), log.Message)
			assert.Equal(t, zapcore.Level(-1), log.Level)
		})
	}
}
