package miro

// Not used when using non-expiring access tokens

/*
const (
	miroAuthURL = "https://api.miro.com/v1/oauth-token"
)

func (m *Miro) authenticate() (token string, err error) {
	httpClient := http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(http.MethodGet, miroAuthURL, nil)
	if err != nil {
		return "", fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+m.accessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("could not execute request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not execute request: status code %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}

	m.logger.Fatalf("%s", string(respBody))

	return "", nil
}
*/
