package logging

import (
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Encoding = "json"
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
