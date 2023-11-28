package tailscale

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

type auditLogResponse struct {
	Version string     `json:"version"`
	Logs    []AuditLog `json:"logs"`
}

type Actor struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	LoginName   string `json:"loginName"`
	DisplayName string `json:"displayName"`
}

type Target struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Property string `json:"property,omitempty"`
}

type AuditLog struct {
	EventTime    time.Time   `json:"eventTime"`
	Type         string      `json:"type"`
	DeferredAt   time.Time   `json:"deferredAt"`
	EventGroupID string      `json:"eventGroupID"`
	Origin       string      `json:"origin"`
	Actor        Actor       `json:"actor"`
	Target       Target      `json:"target,omitempty"`
	Action       string      `json:"action"`
	Old          interface{} `json:"old,omitempty"`
	New          interface{} `json:"new,omitempty"`
}

func (ts *Tailscale) GetAuditLogs(lookbackDays uint) ([]AuditLog, error) {
	logger := ts.logger.WithField("module", "audit_logs")

	startTimestamp := time.Now().AddDate(0, 0, -1*int(lookbackDays))
	endTimestamp := time.Now()

	logger.WithField("start_time", startTimestamp.Format(tailscaleTimestampFormat)).Debug("fetching audit logs")

	auditLogsURL := fmt.Sprintf(
		"%s/tailnet/%s/logs?start=%s&end=%s",
		apiURL, ts.tailnetName,
		startTimestamp.Format(tailscaleTimestampFormat), endTimestamp.Format(tailscaleTimestampFormat),
	)

	resp, err := ts.client.Get(auditLogsURL)
	if err != nil {
		return nil, fmt.Errorf("bad http response: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode > 399 {
		respBytes, _ := io.ReadAll(resp.Body)
		ts.logger.Debugf("%s", string(respBytes))
		return nil, fmt.Errorf("received error code: %d", resp.StatusCode)
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response: %v", err)
	}

	if ts.logger.IsLevelEnabled(logrus.TraceLevel) {
		ts.logger.Debugf("%s", string(respBytes))
	}

	var jsonResponse auditLogResponse
	if err := json.Unmarshal(respBytes, &jsonResponse); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}

	ts.logger.WithField("total_logs", len(jsonResponse.Logs)).Debug("fetched network logs")

	return jsonResponse.Logs, nil
}
