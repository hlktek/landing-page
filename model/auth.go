package model

//GoogleUserInfo model for get user by token
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	DisplayName   string `json:"name"`
}

// SessionInfo model data save in session
type SessionInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
	DisplayName   string `json:"name"`
	Token         string `json:"token"`
	Wallet        string `json:"wallet"`
}
