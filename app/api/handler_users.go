package api

import (
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {
	params := struct {
		Email string `json:"email"`
	}{}
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&params)
	if err != nil {
		log.Errorf("Unprocessable request: %s", err)
		responseWithErr(w, http.StatusBadRequest, "Unprocessable request")
		return
	}

	resp, err := cfg.dbQueries.CreateUser(req.Context(), params.Email)
	if err != nil {
		log.Errorf("Unable to create user: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Unable to create user")
		return
	}

	responseWithJSON(w, http.StatusCreated, User(resp))
}
