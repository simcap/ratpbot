package main

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/simcap/ratpbot/bot"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/api", api)
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, nil)
}

func api(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	question := &bot.Question{r.URL.Query().Get("q")}
	answer := bot.Reply(question)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"answer":   answer.Text(),
		"question": question.Text,
	})
}
