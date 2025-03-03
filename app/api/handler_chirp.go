package api

import (
	"encoding/json"
	"net/http"

	"github.com/DisguisedTrolley/chirpy/app/internal/database"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func (cfg *apiConfig) addChirp(w http.ResponseWriter, req *http.Request) {
	params := struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}{}

	// decode req
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Errorf("Unable to decode params: %s", err)
		responseWithErr(w, http.StatusUnprocessableEntity, "Unprocessable request")
		return
	}

	// validate chirp
	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		log.Error("Chirp too long")
		responseWithErr(w, http.StatusUnprocessableEntity, "Chirp too long")
		return
	}

	params.Body = cleanedBody

	// create chirp
	resp, err := cfg.dbQueries.CreateChirp(req.Context(), database.CreateChirpParams(params))
	if err != nil {
		log.Errorf("Error creatingchirp: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Error creating chirp")
		return
	}

	responseWithJSON(w, http.StatusCreated, Chirp(resp))
}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, req *http.Request) {
	resp, err := cfg.dbQueries.GetChirps(req.Context())
	if err != nil {
		log.Errorf("Error getting chirps: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	// Convert return type to add json tags
	chirps := []Chirp{}
	for _, v := range resp {
		newChirp := Chirp(v)
		chirps = append(chirps, newChirp)
	}

	w.Header().Add("Content-Type", "application/json")
	responseWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) getSingleChirp(w http.ResponseWriter, req *http.Request) {
	param := req.PathValue("chirpID")
	chirpID, err := uuid.Parse(param) // Convert path param to UUID
	if err != nil {
		log.Errorf("Invalid uuid: %s", err)
		responseWithErr(w, http.StatusNotFound, "chirp not found")
		return
	}

	resp, err := cfg.dbQueries.GetSingleChirp(req.Context(), chirpID)
	if err != nil {
		log.Errorf("Chirp not found: %s", err)
		responseWithErr(w, http.StatusNotFound, "chirp not found")
		return
	}

	responseWithJSON(w, http.StatusOK, Chirp(resp))
}
