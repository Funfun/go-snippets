package config

import (
	"fmt"

	"time"

	"os"

	"github.com/elastic/beats/v7/libbeat/common"
)

type Config struct {
	Project         string `config:"project_id" validate:"required"`
	Topic           string `config:"topic" validate:"required"`
	CredentialsFile string `config:"credentials_file"`
	Subscription    struct {
		Name                string        `config:"name" validate:"required"`
		RetainAckedMessages bool          `config:"retain_acked_messages"`
		RetentionDuration   time.Duration `config:"retention_duration"`
		Create              bool          `config:"create"`

		// Settings for the Pub/Sub receiver
		ConnectionPoolSize int `config:"connection_pool_size"`
	}
	Json struct {
		Enabled               bool   `config:"enabled"`
		AddErrorKey           bool   `config:"add_error_key"`
		FieldsUnderRoot       bool   `config:"fields_under_root"`
		FieldsUseTimestamp    bool   `config:"fields_use_timestamp"`
		FieldsTimestampName   string `config:"fields_timestamp_name"`
		FieldsTimestampFormat string `config:"fields_timestamp_format"`
	}
}

func GetDefaultConfig() Config {
	config := Config{}
	config.Subscription.ConnectionPoolSize = 1
	config.Subscription.Create = true
	config.Json.FieldsTimestampName = "@timestamp"
	return config
}

var DefaultConfig = GetDefaultConfig()

func GetAndValidateConfig(cfg *common.Config) (*Config, error) {
	c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	if d, _ := time.ParseDuration("10m"); c.Subscription.RetentionDuration < d {
		return nil, fmt.Errorf("retention_duration cannot be shorter than 10 minutes")
	}

	if d, _ := time.ParseDuration("168h"); c.Subscription.RetentionDuration > d {
		return nil, fmt.Errorf("retention_duration cannot be longer than 7 days")
	}

	if cxns := c.Subscription.ConnectionPoolSize; cxns < 1 {
		return nil, fmt.Errorf("Connection pool size must be >= 1, got: %d", cxns)
	}

	if c.CredentialsFile != "" {
		if _, err := os.Stat(c.CredentialsFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot find the credentials_file %q", c.CredentialsFile)
		}
	}

	return &c, nil
}
