package main

import (
	"../httpHandler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/laptop", httpHandler.LaptopHandler)
	http.HandleFunc("/api/cancel", httpHandler.CancelHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
