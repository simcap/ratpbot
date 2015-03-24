package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApiDetection(t *testing.T) {
	data := []struct {
		in  string
		out string
	}{
		{in: "6", out: "Ligne 6"},
		{in: "1", out: "Ligne 1"},
		{in: "14", out: "Ligne 14"},
		{in: "A", out: "Rer A"},
		{in: "rerE", out: "Rer E"},
	}

	for _, d := range data {
		q := d.in
		exp := d.out

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("", "?q="+q, nil)
		api(w, req)

		answer := struct {
			Answer   string
			Question string
		}{}

		json.Unmarshal(w.Body.Bytes(), &answer)

		if act := answer.Answer; act != exp {
			t.Errorf("got %q want %q", act, exp)
		}
	}
}
