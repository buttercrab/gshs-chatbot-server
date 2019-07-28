package gshs_chatbot_server

import (
	"github.com/buttercrab/gshs-chatbot-server/httpHandler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/request", httpHandler.RequestHandler)
	http.HandleFunc("/api/cancel", httpHandler.CancelHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
