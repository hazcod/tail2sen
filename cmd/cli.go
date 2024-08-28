package main

import (
	"context"
	"flag"
	"github.com/hazcod/miro2sentinel/config"
	"github.com/hazcod/miro2sentinel/pkg/miro"
	msSentinel "github.com/hazcod/miro2sentinel/pkg/sentinel"
	"github.com/hazcod/miro2sentinel/pkg/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	confFile := flag.String("config", "config.yml", "The YAML configuration file.")
	flag.Parse()

	conf := config.Config{}
	if err := conf.Load(*confFile); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("failed to load configuration")
	}

	if err := conf.Validate(); err != nil {
		logger.WithError(err).WithField("config", *confFile).Fatal("invalid configuration")
	}

	logrusLevel, err := logrus.ParseLevel(conf.Log.Level)
	if err != nil {
		logger.WithError(err).Error("invalid log level provided")
		logrusLevel = logrus.InfoLevel
	}
	logger.SetLevel(logrusLevel)

	//

	miroClient, err := miro.New(logger, conf.Miro.AccessToken)
	if err != nil {
		logger.WithError(err).Fatal("could not create Miro client")
	}

	//

	sentinel, err := msSentinel.New(logger, msSentinel.Credentials{
		TenantID:       conf.Microsoft.TenantID,
		ClientID:       conf.Microsoft.AppID,
		ClientSecret:   conf.Microsoft.SecretKey,
		SubscriptionID: conf.Microsoft.SubscriptionID,
		ResourceGroup:  conf.Microsoft.ResourceGroup,
		WorkspaceName:  conf.Microsoft.WorkspaceName,
	})
	if err != nil {
		logger.WithError(err).Fatal("could not create audit MS Sentinel client")
	}

	logger.Info("fetching Miro audit logs")
	auditLogs, err := miroClient.GetAccessLogs(conf.Miro.LookbackDays)
	if err != nil {
		logger.WithError(err).Fatal("failed to fetch audit logs")
	}

	logger.WithField("total", len(auditLogs)).Info("fetched all audit logs")

	//

	convertedLogs, err := utils.ConvertMirAuditLogoToMap(logger, auditLogs)
	if err != nil {
		logger.WithError(err).Fatal("could not convert Miro audit logs")
	}

	if err := sentinel.SendLogs(ctx, logger,
		conf.Microsoft.DataCollection.Endpoint,
		conf.Microsoft.DataCollection.RuleID,
		conf.Microsoft.DataCollection.StreamName,
		convertedLogs); err != nil {
		logger.WithError(err).Fatal("could not ship audit logs to sentinel")
	}

	logger.Info("finished ingesting")
}
