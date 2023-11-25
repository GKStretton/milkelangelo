package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/vote-collection", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		claims, err := verifyUserRequest(r)
		if err != nil {
			fmt.Printf("error verifying user request: %v\n", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		values := r.URL.Query()
		fmt.Printf("got collection vote request: %v\n", values)

		name := values.Get("name")
		fmt.Fprintf(w, "<h1>%s</h1><p>%s</p>", name, claims.OpaqueUserID)
	})
	http.HandleFunc("/vote-location", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		claims, err := verifyUserRequest(r)
		if err != nil {
			fmt.Printf("error verifying user request: %v\n", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		locationVote, err := getLocationVote(r)
		if err != nil {
			fmt.Printf("failed to parse location vote: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Printf("got location vote request: %v\n", locationVote)

		fmt.Fprintf(w, "<h1>test</h1>\n<p>%f, %f</p><p>%s</p>", locationVote.X, locationVote.Y, claims.OpaqueUserID)
	})
	http.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {
		//todo: verify depth/goo auth token
		//todo: open sse stream
		//todo: send every vote to listener. (this is stateless)
	})

	fmt.Println("listening...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET")
	(*w).Header().Set("Access-Control-Allow-Headers", "Authorization, X-Twitch-Extension-Client-Id")
}

type locationVote struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func getLocationVote(r *http.Request) (*locationVote, error) {
	defer r.Body.Close()

	data := &locationVote{}
	err := json.NewDecoder(r.Body).Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
