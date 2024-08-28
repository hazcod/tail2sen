package miro

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	miroDateFormat      = "2006-01-02T15:04:05.000-0700"
	miroQueryDateFormat = "2006-01-02T15:04:05.000Z07:00"
	miroAuditLogAPIURL  = "https://api.miro.com/v2/audit/logs"
)

type AuditLog struct {
	EventTime    time.Time `json:"TimeGenerated"`
	IPAddress    string    `json:"IPAddress"`
	Organisation string    `json:"Organisation"`
	UserEmail    string    `json:"UserEmail"`
	Details      string    `json:"Details"`
	Event        string    `json:"Event"`
	Object       string    `json:"Object"`
}

// Ref https://developers.miro.com/reference/enterprise-get-audit-logs

type data struct {
	Context struct {
		Team struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"team"`
		IP           string `json:"ip"`
		Organization struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"organization"`
	} `json:"context,omitempty"`
	ID     string `json:"id"`
	Object struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"object"`
	CreatedAt string `json:"createdAt"`
	Details   struct {
		Role          string `json:"role"`
		EffectiveRole string `json:"effectiveRole"`
		AuthType      string `json:"authType"`
		MfaFactorType string `json:"mfaFactorType"`
	} `json:"details,omitempty"`
	CreatedBy struct {
		Type  string `json:"type"`
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"createdBy"`
	Event string `json:"event"`
	Type  string `json:"type"`
}

type auditLogResponse struct {
	Limit  int    `json:"limit"`
	Size   int    `json:"size"`
	Cursor string `json:"cursor"`
	Data   []data `json:"data"`
	Type   string `json:"type"`
}

func (m *Miro) GetAccessLogs(lookbackDays uint) ([]AuditLog, error) {
	logs := make([]AuditLog, 0)

	/*
		accessToken, err := m.authenticate()
		if err != nil {
			return nil, fmt.Errorf("could not authenticate: %w", err)
		}
	*/

	httpClient := http.Client{Timeout: time.Second * 10}

	now := time.Now()
	lookbackDate := now.AddDate(0, 0, -1*int(lookbackDays))

	step := 0
	cursor := ""

	for {
		httpRequest, err := http.NewRequest(http.MethodGet, miroAuditLogAPIURL, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create http request: %v", err)
		}

		query := url.Values{}
		query.Set("createdAfter", lookbackDate.Format(miroQueryDateFormat))
		query.Set("createdBefore", now.Format(miroQueryDateFormat))
		query.Set("limit", "100")
		if cursor != "" {
			query.Set("cursor", cursor)
		}
		httpRequest.URL.RawQuery = query.Encode()

		httpRequest.Header.Set("Authorization", "Bearer "+m.accessToken)
		httpRequest.Header.Set("accept", "application/json")

		m.logger.WithField("url", httpRequest.URL.String()).Debug("fetching audit logs")

		resp, err := httpClient.Do(httpRequest)
		if err != nil {
			return nil, fmt.Errorf("http request failed: %v", err)
		}

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read response: %v", err)
		}

		if m.logger.IsLevelEnabled(logrus.TraceLevel) {
			m.logger.Tracef("%s", string(b))
		}

		resp.Body.Close()

		if resp.StatusCode != 200 {
			m.logger.Debugf("%s", string(b))
			return nil, fmt.Errorf("invalid status code: %d", resp.StatusCode)
		}

		var jsonResp auditLogResponse
		if err := json.Unmarshal(b, &jsonResp); err != nil {
			m.logger.Debugf("%s", string(b))
			return nil, fmt.Errorf("failed decode response: %v", err)
		}

		m.logger.WithField("step", step).WithField("content", len(jsonResp.Data)).
			Debug("fetched audit logs")

		for _, data := range jsonResp.Data {
			timing, err := time.Parse(miroDateFormat, data.CreatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to parse time '%s': %v", data.CreatedAt, err)
			}

			detailsBytes, err := json.Marshal(&data.Details)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal details: %v", err)
			}

			logs = append(logs, AuditLog{
				EventTime:    timing,
				IPAddress:    data.Context.IP,
				Organisation: data.Context.Organization.Name,
				UserEmail:    data.CreatedBy.Email,
				Details:      string(detailsBytes),
				Event:        data.Event,
				Object:       data.Object.Name,
			})
		}

		cursor = jsonResp.Cursor

		if cursor == "" {
			m.logger.Debug("cursor empty, stopping")
			break
		}

		step += 1
		m.logger.WithField("step", step).
			WithField("cursor", cursor).Debug("fetching next miro audit page")
	}

	m.logger.WithField("total", len(logs)).Debug("fetched audit logs")

	return logs, nil
}
