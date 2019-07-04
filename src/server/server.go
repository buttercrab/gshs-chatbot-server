package main

import (
	"../httpHandler"
	"net/http"
)

func main() {
	http.HandleFunc("/api/laptop", httpHandler.LaptopHandler)

	_ = http.ListenAndServe(":3080", nil)
}
