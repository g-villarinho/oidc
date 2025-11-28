package domain

import (
	"slices"
	"strings"
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
	Scopes        []string
	LogoURL       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewClient(clientID, clientSecret, clientName string, redirectURIs, grantTypes, responseTypes, scopes []string, logoURL string) (*Client, error) {
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
		Scopes:        scopes,
		LogoURL:       logoURL,
	}, nil
}

type CreateClientParams struct {
	ClientName    string
	RedirectURIs  []string
	GrantTypes    []string
	ResponseTypes []string
	Scopes        []string
	LogoURL       string
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

func (c *Client) SupportsScopes(requestedScopes []string) bool {
	for _, requested := range requestedScopes {
		requested = strings.TrimSpace(requested)
		if requested == "" {
			continue
		}

		if !slices.Contains(c.Scopes, requested) {
			return false
		}
	}
	return true
}

func (c *Client) HasAnyScope(requestedScopes []string) bool {
	for _, requested := range requestedScopes {
		requested = strings.TrimSpace(requested)
		if requested == "" {
			continue
		}

		if slices.Contains(c.Scopes, requested) {
			return true
		}
	}
	return false
}
