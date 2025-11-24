package domain

const (
	AccessTokenLoginDuration = 24 * 60 * 60 // 24 hours
)

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
