package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"tail2sentinel/config"
	msSentinel "tail2sentinel/pkg/sentinel"
	"tail2sentinel/pkg/tailscale"
	"tail2sentinel/pkg/utils"
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

	ts, err := tailscale.New(logger, conf.Tailscale.TailnetName, conf.Tailscale.ClientID, conf.Tailscale.ClientSecret)
	if err != nil {
		logger.WithError(err).Fatal("could not create onepassword client")
	}

	//

	{
		sentinel, err := msSentinel.New(logger, msSentinel.Credentials{
			TenantID:       conf.Microsoft.TenantID,
			ClientID:       conf.Microsoft.AppID,
			ClientSecret:   conf.Microsoft.SecretKey,
			SubscriptionID: conf.Microsoft.SubscriptionID,
			ResourceGroup:  conf.Microsoft.Audit.ResourceGroup,
			WorkspaceName:  conf.Microsoft.Audit.WorkspaceName,
		})
		if err != nil {
			logger.WithError(err).Fatal("could not create audit MS Sentinel client")
		}

		logger.Info("fetching tailscale audit logs")
		auditLogs, err := ts.GetAuditLogs(conf.Tailscale.LookbackDays)
		if err != nil {
			logger.WithError(err).Fatal("failed to fetch audit logs")
		}

		logger.WithField("total", len(auditLogs)).Info("fetched all audit logs")

		//

		convertedLogs, err := utils.ConvertTSAuditToMap(logger, auditLogs)
		if err != nil {
			logger.WithError(err).Fatal("could not convert tailscale audit logs")
		}

		//

		if conf.Microsoft.Audit.UpdateTable {
			if err := sentinel.CreateAuditTable(ctx, logger, "TailscaleAudit", conf.Microsoft.Audit.RetentionDays); err != nil {
				logger.WithError(err).Fatal("failed to create MS Sentinel table for audit logs")
			}
		}

		if err := sentinel.SendLogs(ctx, logger,
			conf.Microsoft.Audit.DataCollection.Endpoint,
			conf.Microsoft.Audit.DataCollection.RuleID,
			conf.Microsoft.Audit.DataCollection.StreamName,
			convertedLogs); err != nil {
			logger.WithError(err).Fatal("could not ship audit logs to sentinel")
		}
	}

	//
	{
		sentinel, err := msSentinel.New(logger, msSentinel.Credentials{
			TenantID:       conf.Microsoft.TenantID,
			ClientID:       conf.Microsoft.AppID,
			ClientSecret:   conf.Microsoft.SecretKey,
			SubscriptionID: conf.Microsoft.SubscriptionID,
			ResourceGroup:  conf.Microsoft.Network.ResourceGroup,
			WorkspaceName:  conf.Microsoft.Network.WorkspaceName,
		})
		if err != nil {
			logger.WithError(err).Fatal("could not create network MS Sentinel client")
		}

		logger.Info("fetching tailscale network logs")
		networkLogs, err := ts.GetNetworkLogs(conf.Tailscale.LookbackDays)
		if err != nil {
			logger.WithError(err).Fatal("failed to fetch network logs")
		}

		logger.WithField("total", len(networkLogs)).Info("fetched all network logs")

		//

		convertedLogs, err := utils.ConvertTSNetworkToMap(logger, networkLogs)
		if err != nil {
			logger.WithError(err).Fatal("could not convert tailscale network logs")
		}

		//

		if conf.Microsoft.Network.UpdateTable {
			if err := sentinel.CreateNetworkTable(ctx, logger, "TailscaleNetwork", conf.Microsoft.Network.RetentionDays); err != nil {
				logger.WithError(err).Fatal("failed to create MS Sentinel table for network logs")
			}
		}

		if err := sentinel.SendLogs(ctx, logger,
			conf.Microsoft.Network.DataCollection.Endpoint,
			conf.Microsoft.Network.DataCollection.RuleID,
			conf.Microsoft.Network.DataCollection.StreamName,
			convertedLogs); err != nil {
			logger.WithError(err).Fatal("could not ship network logs to sentinel")
		}
	}
}
