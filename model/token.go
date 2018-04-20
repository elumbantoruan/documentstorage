package model

// UserToken holds the user token
type UserToken struct {
	Token string `json:"token"`
}

// NewUserToken creates UserToken payload
func NewUserToken(token string) UserToken {
	return UserToken{
		Token: token,
	}
}
