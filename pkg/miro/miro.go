package miro

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Miro struct {
	logger      *logrus.Logger
	accessToken string
}

func New(l *logrus.Logger, accessToken string) (*Miro, error) {
	if l == nil {
		return nil, fmt.Errorf("empty logger provided")
	}
	if accessToken == "" {
		return nil, fmt.Errorf("empty access token provided")
	}

	m := Miro{
		accessToken: accessToken,
		logger:      l,
	}

	return &m, nil
}
