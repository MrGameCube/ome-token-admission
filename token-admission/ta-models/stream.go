package ta_models

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
	ExpireAt      time.Time
	CreateTokens  bool
}

type StreamResponse struct {
	Entity      *StreamEntity
	StreamToken string
	WatchToken  string
}
