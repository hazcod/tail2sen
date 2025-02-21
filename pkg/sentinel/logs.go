package sentinel

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	maxChunkSize = 1000 * 1000 // 1MB
)

// estimateSize estimates the size of a log batch in bytes
func (s *Sentinel) estimateSize(logs []map[string]string) (int, error) {
	data, err := json.Marshal(logs)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (s *Sentinel) chunkLogs(slice []map[string]string) ([][]map[string]string, error) {
	var chunks [][]map[string]string
	var currentChunk []map[string]string
	var currentSize int

	for _, logEntry := range slice {
		// Estimate size with the new log added
		logSize, err := s.estimateSize([]map[string]string{logEntry}) // Estimate single log size
		if err != nil {
			return nil, fmt.Errorf("error estimating log size: %v", err)
		}

		// If adding this log exceeds maxChunkSize, save the current chunk
		if currentSize+logSize > maxChunkSize {
			if len(currentChunk) == 0 {
				// Handle case where a single log is larger than the max chunk size
				return nil, fmt.Errorf("single log exceeds max chunk size")
			}
			chunks = append(chunks, currentChunk)
			currentChunk = nil
			currentSize = 0
		}

		// Add log to the current chunk
		currentChunk = append(currentChunk, logEntry)
		currentSize += logSize
	}

	// Append any remaining logs
	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks, nil
}

func (s *Sentinel) SendLogs(ctx context.Context, l *logrus.Logger, endpoint, ruleID, streamName string, logs []map[string]string) error {
	logger := l.WithField("module", "sentinel_logs")

	logger.WithField("stream_name", streamName).WithField("total", len(logs)).Info("chunking logs")

	chunkedLogs, err := s.chunkLogs(logs)
	if err != nil {
		return fmt.Errorf("failed to chunk logs: %v", err)
	}

	logger.WithField("chunked_logs", len(chunkedLogs)).Info("sending chunked logs")

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
