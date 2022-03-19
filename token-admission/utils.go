package token_admission

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net/http"
)

// ValidateHMACRequest conforms to the OMV: https://airensoft.gitbook.io/ovenmediaengine/access-control/admission-webhooks#security
func ValidateHMACRequest(req *http.Request, bodyBytes []byte) bool {

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
