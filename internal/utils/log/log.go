package log

import (
	"context"
	"github.com/mayank-sheoran/golang_YT_search_service/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

var (
	rawClient = getRawClient()
)

type genericLog struct {
	Message   string
	Level     string
	Timestamp time.Time
	FlowName  string
}

func getRawClient() *zap.SugaredLogger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	loggerConfig.Encoding = "console"
	loggerClient, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	sugar := loggerClient.Sugar()
	return sugar
}

func generateGenericLog(message string, level string, ctx context.Context) genericLog {
	flowName := utils.GetStringFromContext(ctx, utils.FlowCtx)
	return genericLog{
		Message:   message,
		Level:     level,
		Timestamp: time.Now(),
		FlowName:  flowName,
	}
}

func Info(message string, ctx context.Context) {
	rawClient.Info(generateGenericLog(message, "Info", ctx))
}

func Error(message string, ctx context.Context) {
	rawClient.Error(generateGenericLog(message, "Error", ctx))
}

func Warn(message string, ctx context.Context) {
	rawClient.Warn(generateGenericLog(message, "Warn", ctx))
}

func Debug(message string, ctx context.Context) {
	rawClient.Debug(generateGenericLog(message, "Debug", ctx))
}

func HandleError(err error, ctx context.Context, isPanic bool) {
	if err != nil {
		Error(err.Error()+" | un-expected error", ctx)
		if isPanic {
			panic(err)
		}
	}
}

func HandleErrorWithSuccessMessage(err error, ctx context.Context, successMessage string) {
	if err != nil {
		Error(err.Error()+" | un-expected error", ctx)
	} else {
		Info(successMessage, ctx)
	}
}
