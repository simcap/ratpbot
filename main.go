package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/api", api)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("home.html")
	if err != nil {
		log.Printf("Error parsing template %v", err)
	}
	t.Execute(w, nil)
}

func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	question := r.URL.Query().Get("q")
	reply := BotReply(question)

	json.NewEncoder(w).Encode(reply)
}
