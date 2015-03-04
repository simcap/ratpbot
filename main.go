package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"text/template"
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

	question := &Question{text: r.URL.Query().Get("q")}
	answer := process(question)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"answer":   answer.text,
		"question": question.text,
	})
}

type Question struct {
	text string
}

type Answer struct {
	text string
}

func process(q *Question) Answer {
	text := q.text
	metroEnquiry := regexp.MustCompile("\\b[123456789]\\b{1}|\\b10|11|12|13|14\\b")
	rerEnquiry := regexp.MustCompile("\\b(?:[Rr][Ee][Rr])?([ABCDE]){1}\\b")

	metroMatch := metroEnquiry.FindStringSubmatch(text)
	if metroMatch != nil {
		return Answer{text: fmt.Sprintf("Ligne %s", metroMatch[1])}
	}

	rerMatch := rerEnquiry.FindStringSubmatch(text)
	if rerMatch != nil {
		return Answer{text: fmt.Sprintf("RER %s", rerMatch[1])}
	}
	return Answer{text: "Hmmm..."}
}
