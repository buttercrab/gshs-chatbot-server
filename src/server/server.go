package main

import (
	"../httpHandler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/request", httpHandler.RequestHandler)
	http.HandleFunc("/api/cancel", httpHandler.CancelHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
