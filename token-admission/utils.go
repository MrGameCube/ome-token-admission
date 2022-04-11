package token_admission

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const TOKEN_BYTE_LENGTH = 16

// ValidateHMACRequest conforms to the OMV: https://airensoft.gitbook.io/ovenmediaengine/access-control/admission-webhooks#security
func ValidateHMACRequest(req *http.Request, bodyBytes []byte, key []byte) bool {

	hmacData, err := base64.RawURLEncoding.DecodeString(req.Header.Get("X-OME-Signature"))
	if err != nil {
		log.Println(err)
		return false
	}
	bodyHmac := hmac.New(sha1.New, key)
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

func GenerateWebRTCURL(host string, app string, stream string, token string, port uint, useHTTPS bool) *url.URL {
	var protocol string
	if useHTTPS {
		protocol = "wss"
	} else {
		protocol = "ws"
	}
	return generateOMEURL(fmt.Sprintf("%s:%d", host, port),
		fmt.Sprintf("%s/%s", app, stream), protocol, token)
}

func GenerateRTMPURL(host string, app string, stream string, token string, port uint) *url.URL {
	return generateOMEURL(fmt.Sprintf("%s:%d", host, port),
		fmt.Sprintf("%s/%s", app, stream), "rtmp", token)
}

func GenerateHLSURL(host string, app string, stream string, token string, port uint, useHTTPS bool) *url.URL {
	return generateHTTPURL(fmt.Sprintf("%s:%d", host, port),
		fmt.Sprintf("%s/%s/playlist.m3u8", app, stream), token, useHTTPS)
}

func GenerateDASHURL(host string, app string, stream string, token string, port uint, useHTTPS bool) *url.URL {
	return generateHTTPURL(fmt.Sprintf("%s:%d", host, port),
		fmt.Sprintf("%s/%s/manifest.mpd", app, stream), token, useHTTPS)
}

func GenerateLLDASHURL(host string, app string, stream string, token string, port uint, useHTTPS bool) *url.URL {
	return generateHTTPURL(fmt.Sprintf("%s:%d", host, port),
		fmt.Sprintf("%s/%s/manifest_ll.mpd", app, stream), token, useHTTPS)
}

func generateHTTPURL(host string, path string, token string, useHTTPS bool) *url.URL {
	var protocol string
	if useHTTPS {
		protocol = "https"
	} else {
		protocol = "http"
	}
	return generateOMEURL(host, path, protocol, token)
}

func generateOMEURL(host string, path string, scheme string, token string) *url.URL {
	return &url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: "token=" + token,
	}
}
