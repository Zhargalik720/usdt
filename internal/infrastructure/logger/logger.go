package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"usdt/config"
)

const (
	Debug  = "debug"
	Info   = "info"
	Warn   = "warn"
	Error  = "error"
	Dpanic = "dpanic"
	Panic  = "panic"
	Fatal  = "fatal"
)

func NewLogger(conf config.Config) (*zap.Logger, *os.File) {
	levels := map[string]zapcore.Level{
		Debug:  zapcore.DebugLevel,
		Info:   zapcore.InfoLevel,
		Warn:   zapcore.WarnLevel,
		Error:  zapcore.ErrorLevel,
		Dpanic: zapcore.DPanicLevel,
		Panic:  zapcore.PanicLevel,
		Fatal:  zapcore.FatalLevel,
	}

	zapConf := zap.NewProductionConfig()
	zapConf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConf.EncoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder

	atom := zap.NewAtomicLevelAt(levels[conf.LogLvl])
	zapConf.Level = atom

	consoleEncoder := zapcore.NewConsoleEncoder(zapConf.EncoderConfig)
	consoleSyncer := zapcore.AddSync(os.Stdout)

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		zap.L().Fatal("Failed to open log file", zap.Error(err))
	}
	fileEncoder := zapcore.NewJSONEncoder(zapConf.EncoderConfig)
	fileSyncer := zapcore.AddSync(file)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(consoleSyncer), atom),
		zapcore.NewCore(fileEncoder, fileSyncer, atom),
	)

	logger := zap.New(core)

	return logger.Named(conf.AppName), file
}
