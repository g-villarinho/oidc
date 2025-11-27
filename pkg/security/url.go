package security

import (
	"net/url"
	"strings"
)

type ContinueURLParams struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	Scopes              []string
	State               string
	Nonce               string
	CodeChallenge       string
	CodeChallengeMethod string
}

func GenerateContinueURL(baseURL string, params ContinueURLParams) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	q := u.Query()

	q.Set("client_id", params.ClientID)
	q.Set("redirect_uri", params.RedirectURI)
	q.Set("response_type", params.ResponseType)

	if len(params.Scopes) > 0 {
		q.Set("scope", strings.Join(params.Scopes, " "))
	}

	if params.State != "" {
		q.Set("state", params.State)
	}

	if params.Nonce != "" {
		q.Set("nonce", params.Nonce)
	}

	if params.CodeChallenge != "" {
		q.Set("code_challenge", params.CodeChallenge)

		method := params.CodeChallengeMethod
		if method == "" {
			method = "plain"
		}
		q.Set("code_challenge_method", method)
	}

	u.RawQuery = q.Encode()

	return u.String()
}
