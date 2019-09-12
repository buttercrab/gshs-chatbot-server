package main

import (
	"github.com/buttercrab/gshs-chatbot-server/httpHandler"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/api/laptop", httpHandler.LaptopHandler)
	http.HandleFunc("/api/debateInform", httpHandler.DebateInformHandler)
	http.HandleFunc("/api/debate", httpHandler.DebateHandler)
	http.HandleFunc("/api/cancel", httpHandler.CancelHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
