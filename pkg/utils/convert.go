package utils

import (
	"github.com/hazcod/miro2sentinel/pkg/miro"
	"github.com/sirupsen/logrus"
)

const (
	iso8601Format = "2006-01-02T15:04:05Z"
)

func ConvertMirAuditLogoToMap(_ *logrus.Logger, logs []miro.AuditLog) ([]map[string]string, error) {
	output := make([]map[string]string, len(logs))

	for i, log := range logs {
		output[i] = map[string]string{
			"TimeGenerated": log.EventTime.Format(iso8601Format),
			"IPAddress":     log.IPAddress,
			"Organisation":  log.Organisation,
			"UserEmail":     log.UserEmail,
			"Details":       log.Details,
			"Event":         log.Event,
			"Object":        log.Object,
		}
	}

	return output, nil
}
