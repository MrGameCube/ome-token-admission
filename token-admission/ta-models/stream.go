package ta_models

import (
	"time"
)

type StreamEntity struct {
	ID              int64     `json:"id,omitempty"`
	Title           string    `json:"title,omitempty"`
	StreamName      string    `json:"stream_name,omitempty"`
	ApplicationName string    `json:"application_name,omitempty"`
	CreationDate    time.Time `json:"creation_date"`
	OwnerName       string    `json:"owner_name,omitempty"`
	OwnerID         string    `json:"owner_id,omitempty"`
	Public          bool      `json:"public,omitempty"`
}

type StreamParameters struct {
	Title           string `json:"title,omitempty"`
	StreamName      string `json:"stream_name,omitempty"`
	ApplicationName string `json:"application_name,omitempty"`
	OwnerName       string `json:"owner_name,omitempty"`
	OwnerID         string `json:"owner_id,omitempty"`
	Public          bool   `json:"public,omitempty"`
}

type StreamRequest struct {
	StreamOptions *StreamParameters `json:"stream_options,omitempty"`
	ExpireAt      time.Time         `json:"expire_at"`
	CreateTokens  bool              `json:"create_tokens,omitempty"`
}

type StreamResponse struct {
	Entity      *StreamEntity `json:"entity,omitempty"`
	StreamToken string        `json:"stream_token,omitempty"`
	WatchToken  string        `json:"watch_token,omitempty"`
}
