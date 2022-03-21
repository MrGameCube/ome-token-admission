package main

import (
	"context"
	"errors"
	token_admission "github.com/MrGameCube/ome-token-admission/token-admission"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var tokenAdmission *token_admission.TokenAdmission

func main() {
	tokenAdmission = token_admission.New(nil)
	router := initializeGin()
	server := initializeHTTPServer(router)
	waitForShutdown(server)

	log.Println("Server stopped!")
}

func initializeGin() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles("web/public/index.html")
	router.Static("/static", "./web/public/static")
	registerControllers(router)
	return router
}

func initializeHTTPServer(router *gin.Engine) *http.Server {
	server := &http.Server{
		Addr:    ":8083",
		Handler: router,
	}
	go func() {
		log.Printf("HTTP Server Listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()
	return server
}

func waitForShutdown(server *http.Server) {
	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	// Wartet auf Stoppsignal des OS (Bsp. Strg+C oder kill)
	<-quitChannel
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
