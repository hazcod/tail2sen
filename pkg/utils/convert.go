package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"tail2sentinel/pkg/tailscale"
)

const (
	iso8601Format = "2006-01-02T15:04:05Z"
)

func toJson(obj interface{}) (string, error) {
	switch obj.(type) {
	case string:
		return obj.(string), nil
	}

	b, err := json.Marshal(&obj)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func getIANAProtocolFromNumber(proto int) string {
	switch proto {
	case 0:
		return "HOPORT"
	case 1:
		return "ICMP"
	case 2:
		return "IGMP"
	case 3:
		return "GGP"
	case 4:
		return "IPv4"
	case 5:
		return "ST"
	case 6:
		return "TCP"
	case 7:
		return "CBT"
	case 41:
		return "IPv6"
	case 43:
		return "IPv6"
	case 44:
		return "IPv6"
	case 47:
		return "GRE"
	case 143:
		return "Ethernet"
	default:
		return fmt.Sprintf("%d", proto)
	}
}

func ConvertTSAuditToMap(_ *logrus.Logger, logs []tailscale.AuditLog) ([]map[string]string, error) {
	output := make([]map[string]string, len(logs))

	for i, log := range logs {
		old, err := toJson(log.Old)
		if err != nil {
			return nil, fmt.Errorf("couldnt convert old: %v", err)
		}

		edited, err := toJson(log.New)
		if err != nil {
			return nil, fmt.Errorf("couldnt convert new: %v", err)
		}

		actor, err := toJson(log.Actor)
		if err != nil {
			return nil, fmt.Errorf("couldnt convert actor: %v", err)
		}

		target, err := toJson(log.Target)
		if err != nil {
			return nil, fmt.Errorf("couldnt convert target: %v", err)
		}

		output[i] = map[string]string{
			"TimeGenerated": log.EventTime.Format(iso8601Format),
			"Action":        log.Action,
			"ActionType":    log.Type,
			"Origin":        log.Origin,
			"Actor":         actor,
			"Target":        target,
			"Old":           old,
			"New":           edited,
		}
	}

	return output, nil
}

func ConvertTSNetworkToMap(l *logrus.Logger, logs []tailscale.NetworkLog) ([]map[string]string, error) {
	output := make([]map[string]string, 0)

	for _, log := range logs {
		for i, traffic := range log.VirtualTraffic {
			output = append(output, map[string]string{
				"TimeGenerated": log.Logged.Format(iso8601Format),
				"NodeID":        log.NodeID,
				"Start":         log.Start.Format(iso8601Format),
				"End":           log.End.Format(iso8601Format),
				"Index":         fmt.Sprintf("%d", i),

				"Protocol": getIANAProtocolFromNumber(traffic.Proto),
				"Src":      traffic.Src,
				"Dst":      traffic.Dst,
				"Bytes":    fmt.Sprintf("%d", traffic.RxBytes),
				"Packets":  fmt.Sprintf("%d", traffic.RxPkts),
			})
		}
	}

	return output, nil
}
