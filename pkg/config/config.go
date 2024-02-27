package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	defaultClientTimeout = 2 * time.Second
	defaultHTTPPort      = 9747
	defaultRetryCount    = 3
)

type ExporterConfig struct {
	HTTPPort      int
	ClientTimeout time.Duration
	RetryCount    int
	Log           *zap.Logger
}

// InitConfig initializes a config and configure viper to receive config from file and environment.
func InitConfig(configFilePath string) (*ExporterConfig, error) {
	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Unable to create logger", zap.Error(err))
	}

	config := &ExporterConfig{
		HTTPPort:      defaultHTTPPort,
		ClientTimeout: defaultClientTimeout,
		RetryCount:    defaultRetryCount,
		Log:           log,
	}

	// If a config file is specified, read it in.
	if configFilePath != "" {
		viper.SetConfigType("yaml")
		viper.SetConfigFile(configFilePath)
		viper.AutomaticEnv()
	} else {
		config.Log.Info("No config file specified, using default values.")
		return config, nil
	}

	// If a config file found, read it in.
	readConfigErr := viper.ReadInConfig()
	if readConfigErr == nil {
		log.Info(fmt.Sprintf("Using config file: %s", viper.ConfigFileUsed()))

		clientTimeout := viper.GetDuration("client_timeout")
		if clientTimeout > 0 {
			config.ClientTimeout = clientTimeout
		}
		httpPort := viper.GetInt("http_port")
		if httpPort > 0 {
			config.HTTPPort = httpPort
		}

		retryCount := viper.GetInt("retry_count")
		if retryCount > 0 {
			config.RetryCount = retryCount
		}

		return config, nil
	}

	log.Error("Unable to read config", zap.Error(readConfigErr))
	return config, readConfigErr
}
