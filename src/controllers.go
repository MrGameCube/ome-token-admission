package main

import (
	"encoding/json"
	"github.com/MrGameCube/ome-token-admission/token-admission"
	ta_models "github.com/MrGameCube/ome-token-admission/token-admission/ta-models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var classLogger = log.New(os.Stdout, "Controllers", log.Ldate|log.Lshortfile|log.Ltime)

func registerControllers(router *gin.Engine) {
	classLogger.Println("Registering Controllers")
	router.POST("/api/v1/admission", handleAdmission)
	router.POST("/api/v1/stream", handleCreateStream)
}

func handleAdmission(context *gin.Context) {
	classLogger.Println(context.Request)
	resp, err := tokenAdmission.HandleAdmissionRequest(context.Request)

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

	streamResponse, err := createStreamFromBody(&bodyData)
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

func createStreamFromBody(bodyData *[]byte) (*ta_models.StreamResponse, error) {

	var streamReq ta_models.StreamRequest
	err := json.Unmarshal(*bodyData, &streamReq)
	if err != nil {
		classLogger.Println(err)
		return nil, token_admission.ErrInvalidRequest
	}

	return tokenAdmission.CreateStream(&streamReq)
}
