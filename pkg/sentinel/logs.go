package sentinel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	maxChunkSize = 1 * 1024 * 1024 // 1MB
)

// estimateSize estimates the size of a log batch in bytes
func estimateSize(logs []map[string]string) int {
	data, _ := json.Marshal(logs)
	return len(data)
}

func chunkLogs(slice []map[string]string) [][]map[string]string {
	var chunks [][]map[string]string
	var currentChunk []map[string]string
	for _, logEntry := range slice {
		if estimateSize(append(currentChunk, logEntry)) > maxChunkSize && len(currentChunk) > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = nil
		}
		currentChunk = append(currentChunk, logEntry)
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

func (s *Sentinel) SendLogs(ctx context.Context, l *logrus.Logger, endpoint, ruleID, streamName string, logs []map[string]string) error {
	logger := l.WithField("module", "sentinel_logs")

	logger.WithField("stream_name", streamName).WithField("total", len(logs)).Info("shipping logs")

	chunkedLogs := chunkLogs(logs)
	for i, logsChunk := range chunkedLogs {
		l.WithField("progress", fmt.Sprintf("%d/%d", i+1, len(chunkedLogs))).Debug("ingesting log chunks")

		if len(logsChunk) == 0 {
			l.Warn("processing empty chunk")
			continue
		}

		if err := s.IngestLog(ctx, endpoint, ruleID, streamName, logsChunk); err != nil {
			return fmt.Errorf("could not ingest log: %v", err)
		}
	}

	//

	logger.WithField("stream_name", streamName).Info("shipped logs")

	return nil
}
