package models

type CreateClientRequest struct {
	ClientSecret  string   `json:"client_secret" binding:"required,min=8"`
	ClientName    string   `json:"client_name" binding:"required"`
	RedirectURIs  []string `json:"redirect_uris" binding:"required,min=1"`
	GrantTypes    []string `json:"grant_types" binding:"required,min=1"`
	ResponseTypes []string `json:"response_types" binding:"required,min=1"`
	Scopes        []string `json:"scopes" binding:"required,min=1"`
	LogoURL       string   `json:"logo_url"`
}

type UpdateClientRequest struct {
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
