package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/vote", voteHandler)
	http.HandleFunc("/listen", listenHandler)

	fmt.Println("listening...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Authorization, X-Twitch-Extension-Client-Id")
}
