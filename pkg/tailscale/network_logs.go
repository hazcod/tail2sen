package tailscale

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

const (
	tailscaleTimestampFormat = "2006-01-02T15:04:05.000Z"
)

type networkLogResponse struct {
	Logs []NetworkLog `json:"logs"`
}

type VirtualTraffic struct {
	Proto   int    `json:"proto"`
	Src     string `json:"src"`
	Dst     string `json:"dst"`
	RxPkts  int    `json:"rxPkts"`
	RxBytes int    `json:"rxBytes"`
}

type NetworkLog struct {
	Logged         time.Time        `json:"logged"`
	NodeID         string           `json:"nodeId"`
	Start          time.Time        `json:"start"`
	End            time.Time        `json:"end"`
	VirtualTraffic []VirtualTraffic `json:"virtualTraffic"`
}

func (ts *Tailscale) GetNetworkLogs(lookback time.Duration) ([]NetworkLog, error) {
	logger := ts.logger.WithField("module", "network_logs")

	endTimestamp := time.Now()
	startTimestamp := endTimestamp.Add(-lookback)

	logger.WithField("start_time", startTimestamp.Format(tailscaleTimestampFormat)).Debug("fetching network logs")

	networkLogsURL := fmt.Sprintf(
		"%s/tailnet/%s/network-logs?start=%s&end=%s",
		apiURL, ts.tailnetName,
		startTimestamp.Format(tailscaleTimestampFormat), endTimestamp.Format(tailscaleTimestampFormat),
	)

	resp, err := ts.client.Get(networkLogsURL)
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

	var jsonResponse networkLogResponse
	if err := json.Unmarshal(respBytes, &jsonResponse); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}

	ts.logger.WithField("total_logs", len(jsonResponse.Logs)).Debug("fetched network logs")

	return jsonResponse.Logs, nil
}
