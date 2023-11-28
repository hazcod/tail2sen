package sentinel

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	insights "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	"github.com/sirupsen/logrus"
	"time"
)

func (s *Sentinel) CreateNetworkTable(ctx context.Context, l *logrus.Logger, tableName string, retentionDays uint32) error {
	logger := l.WithField("module", "sentinel_network")

	tablesClient, err := insights.NewTablesClient(s.creds.SubscriptionID, s.azCreds, nil)
	if err != nil {
		return fmt.Errorf("could not create ms graph table client: %v", err)
	}

	retention := int32(retentionDays)

	logger.WithField("table_name", tableName).Info("creating or updating table")

	if _, err = tablesClient.Migrate(ctx, s.creds.ResourceGroup, s.creds.WorkspaceName, tableName, nil); err != nil {
		logger.WithError(err).Debug("could not migrate table")
	}

	poller, err := tablesClient.BeginCreateOrUpdate(ctx,
		s.creds.ResourceGroup, s.creds.WorkspaceName, tableName,
		insights.Table{
			Properties: &insights.TableProperties{
				RetentionInDays:      &retention,
				TotalRetentionInDays: to.Ptr[int32](retention * 2),
				Schema: &insights.Schema{
					Columns: []*insights.Column{
						{
							Name: to.Ptr[string]("TimeGenerated"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDateTime),
						},
						{
							Name: to.Ptr[string]("NodeID"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Start"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDateTime),
						},
						{
							Name: to.Ptr[string]("End"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDateTime),
						},
						{
							Name: to.Ptr[string]("Index"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumInt),
						},
						{
							Name: to.Ptr[string]("Protocol"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Src"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Dst"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Bytes"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumInt),
						},
						{
							Name: to.Ptr[string]("Packets"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumInt),
						},
					},
					Name:        to.Ptr[string](tableName),
					Description: to.Ptr[string]("Table that contains events ingested from 1Password."),
				},
			},
		}, nil)
	if err != nil {
		return fmt.Errorf("could not create table '%s': %v", tableName, err)
	}

	_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		return fmt.Errorf("could not poll table creation: %v", err)
	}

	logger.WithField("table_name", tableName).Info("created table")

	return nil
}

func (s *Sentinel) CreateAuditTable(ctx context.Context, l *logrus.Logger, tableName string, retentionDays uint32) error {
	logger := l.WithField("module", "sentinel_audit")

	tablesClient, err := insights.NewTablesClient(s.creds.SubscriptionID, s.azCreds, nil)
	if err != nil {
		return fmt.Errorf("could not create ms graph table client: %v", err)
	}

	retention := int32(retentionDays)

	logger.WithField("table_name", tableName).Info("creating or updating table")

	if _, err = tablesClient.Migrate(ctx, s.creds.ResourceGroup, s.creds.WorkspaceName, tableName, nil); err != nil {
		logger.WithError(err).Debug("could not migrate table")
	}

	poller, err := tablesClient.BeginCreateOrUpdate(ctx,
		s.creds.ResourceGroup, s.creds.WorkspaceName, tableName,
		insights.Table{
			Properties: &insights.TableProperties{
				RetentionInDays:      &retention,
				TotalRetentionInDays: to.Ptr[int32](retention * 2),
				Schema: &insights.Schema{
					Columns: []*insights.Column{
						{
							Name: to.Ptr[string]("TimeGenerated"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumDateTime),
						},
						{
							Name: to.Ptr[string]("Action"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Type"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Origin"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Actor"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Target"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("Old"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
						{
							Name: to.Ptr[string]("New"),
							Type: to.Ptr[insights.ColumnTypeEnum](insights.ColumnTypeEnumString),
						},
					},
					Name:        to.Ptr[string](tableName),
					Description: to.Ptr[string]("Table that contains events ingested from 1Password."),
				},
			},
		}, nil)
	if err != nil {
		return fmt.Errorf("could not create table '%s': %v", tableName, err)
	}

	_, err = poller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{Frequency: time.Second})
	if err != nil {
		return fmt.Errorf("could not poll table creation: %v", err)
	}

	logger.WithField("table_name", tableName).Info("created table")

	return nil
}
