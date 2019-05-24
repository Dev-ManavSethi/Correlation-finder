package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func FatalOnError(err error, message string) {
	if err != nil {
		log.Println(message)
		log.Fatal(err)
	}
}

func main() {

	http.Handle("/socket", websocket.Handler(websocketHandler))
	error1 := http.ListenAndServe(":8000", nil)
	FatalOnError(error1, "Error starting Server")

}
