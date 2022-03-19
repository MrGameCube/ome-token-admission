package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"ts-stream/token-admission"
)

func registerControllers(router *gin.Engine) {
	router.POST("/api/admission", handleAdmission)
	router.POST("/api/stream")
	//router.POST("/api/stream/:app:/:stream:/token")
}

func handleAdmission(context *gin.Context) {
	admissionReq := token_admission.OMEAdmissionBody{}
	bodyBytes, _ := ioutil.ReadAll(context.Request.Body)
	if !token_admission.ValidateHMACRequest(context.Request, bodyBytes) {
		context.Status(http.StatusUnauthorized)
		return
	}
	json.Unmarshal(bodyBytes, &admissionReq)
	reqUrl, _ := url.Parse(admissionReq.Request.URL)
	log.Println(admissionReq)
	log.Println("Token: ", reqUrl.Query().Get("token"))
	context.JSON(http.StatusOK, token_admission.OMEAdmissionResponse{
		Allowed: true,
	})
}
