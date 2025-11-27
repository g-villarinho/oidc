package models

import (
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/pkg/security"
)

type AuthorizePayload struct {
	ClientID            string   `query:"client_id" binding:"required"`
	RedirectURI         string   `query:"redirect_uri" binding:"required,url"`
	ResponseType        string   `query:"response_type" binding:"required,oneof=code token id_token"`
	Scopes              []string `query:"scopes" binding:"dive,required"`
	State               string   `query:"state"`
	Nonce               string   `query:"nonce"`
	CodeChallenge       string   `query:"code_challenge"`
	CodeChallengeMethod string   `query:"code_challenge_method" binding:"omitempty,oneof=plain S256"`
}

func ToContinueURLParams(payload AuthorizePayload) security.ContinueURLParams {
	return security.ContinueURLParams{
		ClientID:            payload.ClientID,
		RedirectURI:         payload.RedirectURI,
		ResponseType:        payload.ResponseType,
		Scopes:              payload.Scopes,
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
		Scopes:              p.Scopes,
		State:               p.State,
		Nonce:               p.Nonce,
		CodeChallenge:       p.CodeChallenge,
		CodeChallengeMethod: p.CodeChallengeMethod,
	}
}
