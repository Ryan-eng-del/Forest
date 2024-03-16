package commonDto

type TokensOutput struct {
	AccessToken string `json:"access_token" form:"access_token"` //access_token
	ExpiresIn int    `json:"expires_in" form:"expires_in"`     //expires_in
	TokenType string `json:"token_type" form:"token_type"`     //token_type
	Scope     string `json:"scope" form:"scope"`               //scope
}