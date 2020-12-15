package dto

import "time"

// author zhasulan
// created on 27.10.20 14:05

type UserTokenData struct {
	Hash      string    `json:"hash"`
	ExpiresAt time.Time `json:"expires_at"`
	//Subject   string    `json:"subject,omitempty"`
}

type UserTokenPair struct {
	Access  UserTokenData `json:"access"`
	Refresh UserTokenData `json:"refresh"`
}

type LoginResult struct {
	UserTokenPair
	PassChangeNeeded bool `json:"pass_change_needed"`
}

type Token struct {
	TokenHash string `json:"token_hash"`
}
