package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

var profaneWords = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

func responseWithJSON(w http.ResponseWriter, code int, payload any) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

	return nil
}

func responseWithErr(w http.ResponseWriter, code int, msg string) error {
	return responseWithJSON(w, code, map[string]string{"error": msg})
}

// Should have reused the words array instead of new cleanWords array, damnit
func cleanProfanity(body string) string {
	words := strings.Split(body, " ")
	cleanWords := []string{}

	for _, word := range words {
		if _, ok := profaneWords[strings.ToLower(word)]; ok {
			cleanWords = append(cleanWords, "****")
		} else {
			cleanWords = append(cleanWords, word)
		}
	}

	return strings.Join(cleanWords, " ")
}
