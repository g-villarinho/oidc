package models

import (
	"strings"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/pkg/oauth"
)

type AuthorizePayload struct {
	ClientID            string `query:"client_id" validate:"required"`
	RedirectURI         string `query:"redirect_uri" validate:"required,url"`
	ResponseType        string `query:"response_type" validate:"required,oneof=code token id_token"`
	Scope               string `query:"scope" validate:"required"`
	State               string `query:"state"`
	Nonce               string `query:"nonce"`
	CodeChallenge       string `query:"code_challenge"`
	CodeChallengeMethod string `query:"code_challenge_method" validate:"omitempty,oneof=plain S256"`
}

type ExchangeTokenPayload struct {
	GrantType    string `form:"grant_type" validate:"required,oneof=authorization_code refresh_token"`
	Code         string `form:"code" validate:"required_if=GrantType authorization_code"`
	RedirectURI  string `form:"redirect_uri" validate:"required_if=GrantType authorization_code,omitempty,url"`
	ClientID     string `form:"client_id" validate:"required"`
	ClientSecret string `form:"client_secret" validate:"omitempty"`
	CodeVerifier string `form:"code_verifier" validate:"omitempty"`
	RefreshToken string `form:"refresh_token" validate:"required_if=GrantType refresh_token"`
}

func (p *AuthorizePayload) GetScopes() []string {
	if p.Scope == "" {
		return []string{}
	}

	return strings.Fields(p.Scope)
}

func (p *AuthorizePayload) ToContinueURLParams() oauth.ContinueURLParams {
	return oauth.ContinueURLParams{
		ClientID:            p.ClientID,
		RedirectURI:         p.RedirectURI,
		ResponseType:        p.ResponseType,
		Scopes:              p.GetScopes(),
		State:               p.State,
		Nonce:               p.Nonce,
		CodeChallenge:       p.CodeChallenge,
		CodeChallengeMethod: p.CodeChallengeMethod,
	}
}

func (p *AuthorizePayload) ToAuthorizeParams() domain.AuthorizeParams {
	return domain.AuthorizeParams{
		ClientID:            p.ClientID,
		RedirectURI:         p.RedirectURI,
		ResponseType:        p.ResponseType,
		Scopes:              p.GetScopes(),
		State:               p.State,
		Nonce:               p.Nonce,
		CodeChallenge:       p.CodeChallenge,
		CodeChallengeMethod: p.CodeChallengeMethod,
	}
}
