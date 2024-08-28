package config

import (
	"errors"
	"fmt"
	validator "github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	defaultLogLevel = "DEBUG"
	defaultLookback = 30
)

type Config struct {
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	} `yaml:"log"`

	Miro struct {
		AccessToken  string `yaml:"access_token" env:"MIRO_TOKEN" valid:"minstringlength(3)"`
		LookbackDays uint   `yaml:"lookback_days" env:"MIRO_LOOKBACK_DAYS"`
	} `yaml:"miro"`

	Microsoft struct {
		AppID          string `yaml:"app_id" env:"MS_APP_ID" valid:"minstringlength(3)"`
		SecretKey      string `yaml:"secret_key" env:"MS_SECRET_KEY" valid:"minstringlength(3)"`
		TenantID       string `yaml:"tenant_id" env:"MS_TENANT_ID" valid:"minstringlength(3)"`
		SubscriptionID string `yaml:"subscription_id" env:"MS_SUB_ID" valid:"minstringlength(3)"`

		DataCollection struct {
			Endpoint   string `yaml:"endpoint" env:"MS_DCR_ENDPOINT" valid:"minstringlength(3)"`
			RuleID     string `yaml:"rule_id" env:"MS_DCR_RULE" valid:"minstringlength(3)"`
			StreamName string `yaml:"stream_name" env:"MS_DCR_STREAM" valid:"minstringlength(3)"`
		} `yaml:"dcr"`

		ResourceGroup string `yaml:"resource_group" env:"MS_RSG_ID" valid:"minstringlength(3)"`
		WorkspaceName string `yaml:"workspace_name" env:"MS_WS_NAME" valid:"minstringlength(3)"`

		RetentionDays uint32 `yaml:"retention_days" env:"MS_RETENTION_DAYS"`
		UpdateTable   bool   `yaml:"update_table" env:"MS_UPDATE_TABLE"`
	} `yaml:"microsoft"`
}

func (c *Config) Validate() error {
	if c.Log.Level == "" {
		c.Log.Level = defaultLogLevel
	}

	if c.Miro.LookbackDays == 0 {
		c.Miro.LookbackDays = defaultLookback
	}

	if c.Miro.AccessToken == "" {
		return errors.New("no Miro access token provided")
	}

	if valid, err := validator.ValidateStruct(c); !valid || err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}

	return nil
}

func (c *Config) Load(path string) error {
	if path != "" {
		configBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to load configuration file at '%s': %v", path, err)
		}

		if err = yaml.Unmarshal(configBytes, c); err != nil {
			return fmt.Errorf("failed to parse configuration: %v", err)
		}
	}

	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("could not load environment: %v", err)
	}

	return nil
}
