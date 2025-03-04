package api

import (
	"encoding/json"
	"net/http"

	"github.com/DisguisedTrolley/chirpy/app/internal/auth"
	"github.com/charmbracelet/log"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, req *http.Request) {
	params := struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}{}
	decoder := json.NewDecoder(req.Body)

	// Decode params
	err := decoder.Decode(&params)
	if err != nil {
		log.Errorf("Unable to decode info: %s", err)
		responseWithErr(w, http.StatusUnprocessableEntity, "Unprocessable request")
		return
	}

	// Check if user exists
	user, err := cfg.dbQueries.FindUser(req.Context(), params.Email)
	if err != nil {
		log.Errorf("Error with user: %s", err)
		responseWithErr(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Verify password
	err = auth.VerifyHashedPassword(user.HashedPassword, params.Password)
	if err != nil {
		log.Errorf("Error with password: %s", err)
		responseWithErr(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Get jwtSecret
	jwt, err := cfg.GenJWTkey(user.ID, params.ExpiresInSeconds)
	if err != nil {
		log.Errorf("Error generatign jwt: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	resp := Response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: jwt,
	}

	responseWithJSON(w, http.StatusOK, resp)
}
