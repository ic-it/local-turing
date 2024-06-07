package internal

import "go.uber.org/zap"

var (
	logger *zap.SugaredLogger
)

func SetupLogger(verbose bool) {
	l, _ := zap.NewDevelopment(zap.IncreaseLevel(zap.WarnLevel))
	logger = l.Sugar()
	if verbose {
		l, _ := zap.NewDevelopment()
		logger = l.Sugar()
	}
}

func GetLogger() *zap.SugaredLogger {
	return logger
}
