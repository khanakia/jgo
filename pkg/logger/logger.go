package logger

import (
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
	HOW TO USE
	logger.SugarLogger.Infof("Failed to fetch URL: %s", "test")
	logger.SugarLogger("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", "test",
		"attempt", 3,
		"backoff", time.Second,
	)
*/
// var SugarLogger *zap.SugaredLogger

type Logger struct {
	SugarLogger *zap.SugaredLogger
}

func (logger Logger) Version() string {
	return "0.01"
}

func New() Logger {
	var SugarLogger *zap.SugaredLogger
	var UID string

	writerSyncer := getLogWriter()
	encoder := getEncoder()

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	// logger.With(
	// 	zap.Namespace("metrics"),
	// 	zap.Int("counter", 1),
	// )

	UID = uuid.New().String()
	SugarLogger = logger.With(
		zap.String("uid", UID),
		zap.String("lang", "go"),
	).Sugar()

	return Logger{
		SugarLogger: SugarLogger,
	}
}

func getEncoder() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	return zapcore.NewJSONEncoder(cfg)
	// return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	currentTime := time.Now()
	fileName := currentTime.Format("01-02-2006") + "_scripts.log"

	os.MkdirAll("./logs", os.ModePerm)

	file, _ := os.Create("./logs/" + fileName)
	// file, _ := os.OpenFile("./test1.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return zapcore.AddSync(file)
}
