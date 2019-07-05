package main

import (
	"../httpHandler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/laptop", httpHandler.LaptopHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
