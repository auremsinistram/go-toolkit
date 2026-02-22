package logger

import (
	"github.com/auremsinistram/go-errors"
	"github.com/auremsinistram/go-toolkit/tools"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
	var config zap.Config

	if tools.GetenvBool("DEBUG", false) {
		config = zap.NewDevelopmentConfig()

		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()

		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return nil, errors.Wrap(err, "logger - New - #1")
	}

	return logger, nil
}
