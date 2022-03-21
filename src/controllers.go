package main

import (
	"github.com/MrGameCube/ome-token-admission/token-admission"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerControllers(router *gin.Engine) {
	router.POST("/api/admission", handleAdmission)
	router.POST("/api/stream")
	//router.POST("/api/stream/:app:/:stream:/token")
}

func handleAdmission(context *gin.Context) {
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
