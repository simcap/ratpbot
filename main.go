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

	if item, ok := MetroLineDetector.item(text); ok {
		return Answer{text: fmt.Sprintf("Ligne %s", item)}
	}

	if item, ok := RerLineDetector.item(text); ok {
		return Answer{text: fmt.Sprintf("Rer %s", item)}
	}

	return Answer{text: "Hmmm..."}
}

var (
	MetroLineDetector = Detector{
		regexp.MustCompile("\\b([123456789])\\b|(1[01234])\\b"),
	}
	RerLineDetector = Detector{
		regexp.MustCompile("\\b(?:[Rr][Ee][Rr])?([ABCDE]){1}\\b"),
	}
)

type Detector struct {
	reg *regexp.Regexp
}

func (d *Detector) item(text string) (string, bool) {
	if item := d.reg.FindStringSubmatch(text); len(item) > 1 {
		return item[1], true
	}
	return "", false
}
