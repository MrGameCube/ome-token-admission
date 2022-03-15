package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net/http"
)

// https://airensoft.gitbook.io/ovenmediaengine/access-control/admission-webhooks#security
func validateHMACRequest(req *http.Request, bodyBytes []byte) bool {

	hmacData, err := base64.RawURLEncoding.DecodeString(req.Header.Get("X-OME-Signature"))
	if err != nil {
		log.Println(err)
		return false
	}
	// TODO: make key configurable
	bodyHmac := hmac.New(sha1.New, []byte("1234"))
	bodyHmac.Write(bodyBytes)
	newHmac := bodyHmac.Sum(nil)
	return hmac.Equal(hmacData, newHmac)
}
