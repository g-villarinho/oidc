package domain

type Client struct {
	ID            string
	ClientID      string
	ClientSecret  string
	ClientName    string
	RedirectURIs  []string
	GrantTypes    []string
	ResponseTypes []string
	Scope         string
}
