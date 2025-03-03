package config

import (
	"errors"
	"fmt"
	validator "github.com/asaskevich/govalidator"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

const (
	defaultLogLevel = "DEBUG"
	defaultLookback = "1.2h"
)

type Config struct {
	Log struct {
		Level string `yaml:"level" env:"LOG_LEVEL"`
	} `yaml:"log"`

	Tailscale struct {
		ClientID     string        `yaml:"client_id" env:"TS_CLIENTID" valid:"minstringlength(3)"`
		ClientSecret string        `yaml:"client_secret" env:"TS_CLIENT_SECRET" valid:"minstringlength(3)"`
		TailnetName  string        `yaml:"tailnet" env:"TS_TAILNET" valid:"minstringlength(3)"`
		Lookback     time.Duration `yaml:"lookback" env:"TS_LOOKBACK"`
	} `yaml:"tailscale"`

	Microsoft struct {
		AppID          string `yaml:"app_id" env:"MS_APP_ID" valid:"minstringlength(3)"`
		SecretKey      string `yaml:"secret_key" env:"MS_SECRET_KEY" valid:"minstringlength(3)"`
		TenantID       string `yaml:"tenant_id" env:"MS_TENANT_ID" valid:"minstringlength(3)"`
		SubscriptionID string `yaml:"subscription_id" env:"MS_SUB_ID" valid:"minstringlength(3)"`

		Audit struct {
			DataCollection struct {
				Endpoint   string `yaml:"endpoint" env:"MS_AD_DCR_ENDPOINT" valid:"minstringlength(3)"`
				RuleID     string `yaml:"rule_id" env:"MS_AD_DCR_RULE" valid:"minstringlength(3)"`
				StreamName string `yaml:"stream_name" env:"MS_AD_DCR_STREAM" valid:"minstringlength(3)"`
			} `yaml:"dcr"`

			ResourceGroup string `yaml:"resource_group" env:"MS_AD_RSG_ID" valid:"minstringlength(3)"`
			WorkspaceName string `yaml:"workspace_name" env:"MS_AD_WS_NAME" valid:"minstringlength(3)"`

			RetentionDays uint32 `yaml:"retention_days" env:"MS_AD_RETENTION_DAYS"`
			UpdateTable   bool   `yaml:"update_table" env:"MS_AD_UPDATE_TABLE"`
		} `yaml:"audit_output"`

		Network struct {
			DataCollection struct {
				Endpoint   string `yaml:"endpoint" env:"MS_NW_DCR_ENDPOINT" valid:"minstringlength(3)"`
				RuleID     string `yaml:"rule_id" env:"MS_NW_DCR_RULE" valid:"minstringlength(3)"`
				StreamName string `yaml:"stream_name" env:"MS_NW_DCR_STREAM" valid:"minstringlength(3)"`
			} `yaml:"dcr"`

			ResourceGroup string `yaml:"resource_group" env:"MS_NW_RSG_ID" valid:"minstringlength(3)"`
			WorkspaceName string `yaml:"workspace_name" env:"MS_NW_WS_NAME" valid:"minstringlength(3)"`

			RetentionDays uint32 `yaml:"retention_days" env:"MS_NW_RETENTION_DAYS"`
			UpdateTable   bool   `yaml:"update_table" env:"MS_NW_UPDATE_TABLE"`
		} `yaml:"network_output"`
	} `yaml:"microsoft"`
}

func (c *Config) Validate() error {
	if c.Log.Level == "" {
		c.Log.Level = defaultLogLevel
	}

	if c.Tailscale.Lookback.Seconds() == 0 {
		var err error
		c.Tailscale.Lookback, err = time.ParseDuration(defaultLookback)
		if err != nil {
			logrus.WithError(err).WithField("defaultLookback", defaultLookback).Fatal("could not parse lookback")
		}
	}

	if c.Tailscale.ClientID == "" {
		return errors.New("no clientid provided")
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
