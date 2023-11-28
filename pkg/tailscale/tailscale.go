package tailscale

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

const (
	apiURL = "https://api.tailscale.com/api/v2/"
)

var tsAPIScopes = []string{"all:read"}

type Tailscale struct {
	logger      *logrus.Logger
	client      *http.Client
	tailnetName string
}

func New(logger *logrus.Logger, tailnetName, clientID, clientSecret string) (*Tailscale, error) {
	if logger == nil {
		return nil, fmt.Errorf("nil logger provided")
	}
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("empty api token provided")
	}
	if tailnetName == "" {
		return nil, fmt.Errorf("empty tailnet name provided")
	}

	var oauthConfig = &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     apiURL + "oauth/token",
		Scopes:       tsAPIScopes,
	}

	client := oauthConfig.Client(context.Background())
	if client == nil {
		return nil, fmt.Errorf("could not create tailscale oauth client")
	}

	ts := Tailscale{
		logger:      logger,
		client:      client,
		tailnetName: tailnetName,
	}

	return &ts, nil
}
