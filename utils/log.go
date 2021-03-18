package utils

import (
	"fmt"

	ucfg "github.com/elastic/go-ucfg"
	"go.uber.org/zap"
)

type LoggerConfig struct {
	Environment string `config:"environment"`
	Level       string `config:"logging.level"`
}

var config LoggerConfig = LoggerConfig{
	Environment: "development",
	Level:       "debug",
}

func ConfigureLogger(cfg *ucfg.Config) error {
	err := cfg.Unpack(&config)
	if err != nil {
		return fmt.Errorf("cannot unpack configuration for logging: %v", err.Error())
	}
	return nil
}

func NewLogger(name string) (*zap.SugaredLogger, error) {
	// set the level of the logger
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(config.Level))
	if err != nil {
		return nil, err
	}

	// set the logger
	var logger *zap.Logger
	if config.Environment == "development" {
		devConfig := zap.NewDevelopmentConfig()
		devConfig.Level.SetLevel(level.Level())
		logger, err = devConfig.Build()
	} else {
		prodConfig := zap.NewProductionConfig()
		prodConfig.Level.SetLevel(level.Level())
		logger, err = prodConfig.Build()
	}
	if err != nil {
		return nil, err
	}

	sugaredLogger := logger.Sugar()
	sugaredLogger = sugaredLogger.Named(name)
	return sugaredLogger, err
}
