package api

import (
	"encoding/json"
	"net/http"

	"github.com/DisguisedTrolley/chirpy/app/internal/auth"
	"github.com/DisguisedTrolley/chirpy/app/internal/database"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
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

	responseWithJSON(w, http.StatusCreated, User{
		ID:          resp.ID,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
		Email:       resp.Email,
		IsChirpyRed: resp.IsChirpyRed,
	})
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	// Get access token from header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Errorf("Malformed header: %s", err)
		responseWithErr(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	// Verify access token
	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Errorf("Invalid jwt: %s", err)
		responseWithErr(w, http.StatusUnauthorized, "Invalid access token")
		return
	}

	// Get details from req body
	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		log.Errorf("Invalid params: %s", err)
		responseWithErr(w, http.StatusUnprocessableEntity, "Missing email or password")
		return
	}

	// Get hashed password
	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Errorf("Unable to hash password: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Update user details
	newUserInfo, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userId,
		Email:          params.Email,
		HashedPassword: hashedPw,
	})
	if err != nil {
		log.Errorf("Unable to update details: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseWithJSON(w, http.StatusOK, User{
		ID:          newUserInfo.ID,
		CreatedAt:   newUserInfo.CreatedAt,
		UpdatedAt:   newUserInfo.UpdatedAt,
		Email:       newUserInfo.Email,
		IsChirpyRed: newUserInfo.IsChirpyRed,
	})
}

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}

	// Decode params
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Errorf("No body: %s", err)
		return
	}

	// Verify event type
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Parse user id to uuid
	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Upgrade user
	err = cfg.dbQueries.UpgradeUser(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
