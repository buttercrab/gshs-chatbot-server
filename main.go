package main

import (
	"github.com/buttercrab/gshs-chatbot-server/httpHandler"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/api/request", httpHandler.LaptopHandler)
	http.HandleFunc("/api/cancel", httpHandler.CancelHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
