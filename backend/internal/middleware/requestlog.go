package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func LogRequests(next http.Handler) http.Handler {
	logFile, err := os.OpenFile("traintrack.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}

	return handlers.LoggingHandler(logFile, next)
}
