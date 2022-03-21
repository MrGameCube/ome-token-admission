package token_admission

import "time"

type StreamEntity struct {
	ID              int64
	Title           string
	StreamName      string
	ApplicationName string
	CreationDate    time.Time
	OwnerName       string
	OwnerID         string
	Public          bool
}

type StreamRequest struct {
	StreamOptions StreamEntity
	Expires       int
	CreateTokens  bool
}

type StreamResponse struct {
	Success     bool
	StreamToken string
	WatchToken  string
}
