package token

import "time"

type TokenRequest struct {
	TokenOptions TokenEntity
	Expires      int
}
type TokenResponse struct {
	Success bool
	Token   string
}

type TokenEntity struct {
	ID          int64
	Token       string
	Direction   string
	Stream      string
	Application string
	ExpiresAt   time.Time
}
