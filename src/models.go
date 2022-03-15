package main

import (
	"time"
)

type OMEAdmissionRequest struct {
	Direction string
	Protocol  string
	URL       string
	Time      time.Time
}

type OMEClient struct {
	Address string
	Port    int
}

type OMEAdmissionBody struct {
	Client  OMEClient
	Request OMEAdmissionRequest
}

type OMEAdmissionResponse struct {
	Allowed  bool   `json:"allowed"`
	NewURL   string `json:"new_url"`
	Lifetime int    `json:"lifetime"`
	Reason   string `json:"reason"`
}

type StreamRequest struct {
	StreamOptions StreamInfo
	Expires       int
	CreateTokens  bool
}

type StreamResponse struct {
	Success     bool
	StreamToken string
	WatchToken  string
}

type TokenRequest struct {
	TokenOptions string
	Expires      int
}
type TokenResponse struct {
	Success bool
	Token   string
}

type TokenInfo struct {
	Direction   string
	Stream      string
	Application string
}

type StreamInfo struct {
	Title           string
	StreamName      string
	ApplicationName string
	CreationDate    time.Time
	OwnerName       string
	OwnerID         string
	Public          bool
}
