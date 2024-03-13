package middleware

type AuthUser struct {
	Id          string `json:"sub"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Provider    string `json:"provider"`
	ProviderSub string `json:"providerSub"`
	ExpiresAt   int64  `json:"exp"`
}
