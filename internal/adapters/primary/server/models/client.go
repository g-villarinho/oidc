package models

import "github.com/g-villarinho/oidc-server/internal/core/domain"

type CreateClientPayload struct {
	ClientName    string   `json:"client_name" binding:"required"`
	RedirectURIs  []string `json:"redirect_uris" binding:"required,min=1"`
	GrantTypes    []string `json:"grant_types" binding:"required,min=1"`
	ResponseTypes []string `json:"response_types" binding:"required,min=1"`
	Scopes        []string `json:"scopes" binding:"required,min=1"`
	LogoURL       string   `json:"logo_url"`
}

type UpdateClientPayload struct {
	ClientName    string   `json:"client_name" binding:"required"`
	RedirectURIs  []string `json:"redirect_uris" binding:"required,min=1"`
	GrantTypes    []string `json:"grant_types" binding:"required,min=1"`
	ResponseTypes []string `json:"response_types" binding:"required,min=1"`
	Scopes        []string `json:"scopes" binding:"required,min=1"`
}

type ClientResponse struct {
	ID            string   `json:"id"`
	ClientID      string   `json:"client_id"`
	ClientName    string   `json:"client_name"`
	RedirectURIs  []string `json:"redirect_uris"`
	GrantTypes    []string `json:"grant_types"`
	ResponseTypes []string `json:"response_types"`
	Scopes        []string `json:"scopes"`
	LogoURL       string   `json:"logo_url"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
}

type ClientListResponse struct {
	Clients []ClientResponse `json:"clients"`
	Total   int              `json:"total"`
}

func ToCreateClientParams(req CreateClientPayload) domain.CreateClientParams {
	return domain.CreateClientParams{
		ClientName:    req.ClientName,
		RedirectURIs:  req.RedirectURIs,
		GrantTypes:    req.GrantTypes,
		ResponseTypes: req.ResponseTypes,
		Scopes:        req.Scopes,
		LogoURL:       req.LogoURL,
	}
}
