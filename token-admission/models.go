package token_admission

import "time"

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
