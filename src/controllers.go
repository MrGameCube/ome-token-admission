package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func registerControllers(router *gin.Engine) {
	router.POST("/api/admission", handleAdmission)
	router.POST("/api/stream")
	router.POST("/api/stream/:app:/:stream:/token")
}

func handleAdmission(context *gin.Context) {
	admissionReq := OMEAdmissionBody{}
	bodyBytes, _ := ioutil.ReadAll(context.Request.Body)

	json.Unmarshal(bodyBytes, &admissionReq)
	reqUrl, _ := url.Parse(admissionReq.Request.URL)
	log.Println(admissionReq)
	log.Println("Token: ", reqUrl.Query().Get("token"))
	context.JSON(http.StatusOK, OMEAdmissionResponse{
		Allowed: true,
	})
}
