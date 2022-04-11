package main

import (
	"encoding/json"
	"fmt"
	"github.com/MrGameCube/ome-token-admission/token-admission"
	ta_models "github.com/MrGameCube/ome-token-admission/token-admission/ta-models"
	"github.com/gin-gonic/gin"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var classLogger = log.New(os.Stdout, "Controllers", log.Ldate|log.Lshortfile|log.Ltime)

func registerControllers(router *gin.Engine) {
	classLogger.Println("Registering Controllers")
	router.LoadHTMLGlob("web/**/*.html")
	router.POST("/api/v1/admission", handleAdmission)
	router.POST("/api/v1/stream", handleCreateStream)
	router.GET("/play/:token", handlePlayStream)
	router.GET("/stream/:token", handleStartStream)
}

func handleAdmission(context *gin.Context) {
	resp, err := tokenAdmission.HandleAdmissionRequest(context.Request, []byte(config.OMESharedSecret))

	if err == token_admission.ErrInvalidSignature {
		context.Status(http.StatusUnauthorized)
		return
	}

	if err == token_admission.ErrTokenMissing {
		context.Status(http.StatusBadRequest)
		return
	}

	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, resp)
}

func handleCreateStream(context *gin.Context) {

	// Reading the body data
	bodyData, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		classLogger.Println(err)
		context.Status(http.StatusBadRequest)
	}

	// Authentication with shared secret
	//if !token_admission.ValidateHMACRequest(context.Request, bodyData, []byte("1234")) {
	//	context.Status(http.StatusUnauthorized)
	//	return
	//}

	streamResponse, err := createStreamFromBody(&bodyData, context.Request)
	if err == token_admission.ErrInvalidRequest {
		classLogger.Println(err)
		context.Status(http.StatusBadRequest)
		return
	}
	if err != nil {
		classLogger.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, streamResponse)
}

func createStreamFromBody(bodyData *[]byte, req *http.Request) (*StreamResponseWrapper, error) {

	var streamReq ta_models.StreamRequest
	err := json.Unmarshal(*bodyData, &streamReq)
	if err != nil {
		classLogger.Println(err)
		return nil, token_admission.ErrInvalidRequest
	}
	streamResp, err := tokenAdmission.CreateStream(&streamReq)
	if err != nil {
		return nil, err
	}
	var streamLink, watchLink, rtmpLink string

	if streamResp.StreamToken != "" {
		streamURL := config.BaseURL
		streamURL.Path = fmt.Sprintf("stream/%s", streamResp.StreamToken)
		streamLink = streamURL.String()
		rtmpURL := token_admission.GenerateRTMPURL(config.OmeURL, streamResp.Entity.ApplicationName, streamResp.Entity.StreamName, streamResp.StreamToken, config.RTMPPort)
		rtmpLink = rtmpURL.String()
	}
	if streamResp.WatchToken != "" {
		watchURL := config.BaseURL
		watchURL.Path = fmt.Sprintf("play/%s", streamResp.WatchToken)
		watchLink = watchURL.String()
	}
	return &StreamResponseWrapper{
		Response:  streamResp,
		StreamURL: streamLink,
		WatchURL:  watchLink,
		RTMPURL:   rtmpLink,
	}, nil
}

func handlePlayStream(context *gin.Context) {
	token := html.EscapeString(context.Param("token"))
	tokenInfo, err := tokenAdmission.GetTokenInfo(token)
	if err != nil || tokenInfo.Direction != ta_models.DirectionOutgoing {
		context.Status(http.StatusNotFound)
		return
	}
	streamInfo, err := tokenAdmission.GetStreamInfo(tokenInfo.Application, tokenInfo.Stream)
	if err != nil || strings.TrimSpace(token) == "" {
		classLogger.Println("No stream for token:", tokenInfo)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.HTML(http.StatusOK, "player.html", gin.H{
		"targetStream": gin.H{
			"title":     streamInfo.Title,
			"webRTCURL": token_admission.GenerateWebRTCURL(config.OmeURL, streamInfo.ApplicationName, streamInfo.StreamName, token, config.WebRTCPort, config.UseHTTPS).String(),
		},
	})
}

func handleStartStream(context *gin.Context) {
	token := html.EscapeString(context.Param("token"))
	tokenInfo, err := tokenAdmission.GetTokenInfo(token)
	if err != nil || tokenInfo.Direction != ta_models.DirectionIncoming {
		context.Status(http.StatusNotFound)
		return
	}
	streamInfo, err := tokenAdmission.GetStreamInfo(tokenInfo.Application, tokenInfo.Stream)
	if err != nil || strings.TrimSpace(token) == "" {
		classLogger.Println("No stream for token:", tokenInfo)
		context.Status(http.StatusInternalServerError)
		return
	}

	context.HTML(http.StatusOK, "streamer.html", gin.H{
		"targetStream": gin.H{
			"webRTCSendURL": token_admission.GenerateWebRTCURL(config.OmeURL, streamInfo.ApplicationName,
				streamInfo.StreamName, token, config.WebRTCPort, config.UseHTTPS).String() + "&direction=send",
		},
	})
}
