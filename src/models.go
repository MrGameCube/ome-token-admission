package main

import (
	ta_models "github.com/MrGameCube/ome-token-admission/token-admission/ta-models"
)

type StreamResponseWrapper struct {
	Response  *ta_models.StreamResponse `json:"response"`
	StreamURL string                    `json:"stream_url,string,omitempty"`
	WatchURL  string                    `json:"watch_url,string,omitempty"`
	RTMPURL   string                    `json:"rtmp_url,string,omitempty"`
}
