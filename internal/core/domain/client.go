package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID            uuid.UUID
	ClientID      string
	ClientSecret  string
	ClientName    string
	RedirectURIs  []string
	GrantTypes    []string
	ResponseTypes []string
	Scope         string
	LogoURL       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewClient(clientID, clientSecret, clientName string, redirectURIs, grantTypes, responseTypes []string, scope string) (*Client, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Client{
		ID:            id,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		ClientName:    clientName,
		RedirectURIs:  redirectURIs,
		GrantTypes:    grantTypes,
		ResponseTypes: responseTypes,
		Scope:         scope,
	}, nil
}
