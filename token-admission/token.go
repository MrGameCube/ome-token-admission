package token_admission

import "time"

const (
	DirectionOutgoing = "outgoing"
	DirectionIncoming = "incoming"
)

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

type TokenOptions struct {
	Direction   string
	Stream      string
	Application string
	ExpiresAt   time.Time
}
