package token_admission

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const TOKEN_BYTE_LENGTH = 16

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

func parseStreamFromURL(reqURL string) (appName string, streamName string, token string) {
	parsedURL, err := url.Parse(reqURL)
	if err != nil {
		return "", "", ""
	}
	pathElements := strings.Split(parsedURL.Path, "/")
	if len(pathElements) < 3 {
		return "", "", parsedURL.Query().Get("token")
	}
	return pathElements[1], pathElements[2], parsedURL.Query().Get("token")
}

func generateToken() (string, error) {
	tokenBytes := make([]byte, TOKEN_BYTE_LENGTH)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tokenBytes), nil
}
