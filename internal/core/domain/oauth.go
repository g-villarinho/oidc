package domain

import "errors"

var (
	ErrClientNotFound          = errors.New("client not found")
	ErrInvalidRedirectURI      = errors.New("invalid redirect URI")
	ErrUnauthorizedClient      = errors.New("unauthorized client")
	ErrUnsupportedResponseType = errors.New("unsupported response type")
	ErrInvalidScope            = errors.New("invalid scope")
)

type AuthorizeParams struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	Scopes              []string
	State               string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
}

type ExchangeTokenParams struct {
	GrantType    string
	Code         string
	RedirectURI  string
	ClientID     string
	ClientSecret string
	CodeVerifier string
	RefreshToken string
}
