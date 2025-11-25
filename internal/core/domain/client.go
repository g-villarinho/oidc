package domain

import (
	"slices"
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

func NewClient(clientID, clientSecret, clientName string, redirectURIs, grantTypes, responseTypes []string, scope string, logoURL string) (*Client, error) {
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
		LogoURL:       logoURL,
	}, nil
}

func (c *Client) HasRedirectURI(uri string) bool {
	return slices.Contains(c.RedirectURIs, uri)
}

func (c *Client) SupportsGrantType(grantType string) bool {
	return slices.Contains(c.GrantTypes, grantType)
}

func (c *Client) SupportsResponseType(responseType string) bool {
	return slices.Contains(c.ResponseTypes, responseType)
}
