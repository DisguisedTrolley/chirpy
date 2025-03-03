package api

import (
	"encoding/json"
	"net/http"

	"github.com/DisguisedTrolley/chirpy/app/internal/auth"
	"github.com/DisguisedTrolley/chirpy/app/internal/database"
	"github.com/charmbracelet/log"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, req *http.Request) {
	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(req.Body)

	// Decode params
	err := decoder.Decode(&params)
	if err != nil {
		log.Errorf("Unprocessable request: %s", err)
		responseWithErr(w, http.StatusBadRequest, "Unprocessable request")
		return
	}

	// Generate hash
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Errorf("Error generating hash: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Create user
	resp, err := cfg.dbQueries.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		log.Errorf("Unable to create user: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Unable to create user")
		return
	}

	responseWithJSON(w, http.StatusCreated, User(resp))
}
