package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DisguisedTrolley/chirpy/app/internal/auth"
	"github.com/google/uuid"
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

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("Chirp too long")
	}

	cleanBody := cleanProfanity(body)

	return cleanBody, nil
}

func (cfg *apiConfig) GenJWTkey(userID uuid.UUID) (string, error) {
	exp := time.Hour

	jwt, err := auth.MakeJWT(userID, cfg.jwtSecret, exp)
	if err != nil {
		return "", err
	}

	return jwt, nil
}
