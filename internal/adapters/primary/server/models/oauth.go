package models

import (
	"strings"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/pkg/security"
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

func (p *AuthorizePayload) GetScopes() []string {
	if p.Scope == "" {
		return []string{}
	}
	return strings.Fields(p.Scope)
}

func ToContinueURLParams(payload AuthorizePayload) security.ContinueURLParams {
	return security.ContinueURLParams{
		ClientID:            payload.ClientID,
		RedirectURI:         payload.RedirectURI,
		ResponseType:        payload.ResponseType,
		Scopes:              payload.GetScopes(),
		State:               payload.State,
		Nonce:               payload.Nonce,
		CodeChallenge:       payload.CodeChallenge,
		CodeChallengeMethod: payload.CodeChallengeMethod,
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
