package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DisguisedTrolley/chirpy/app/internal/auth"
	"github.com/DisguisedTrolley/chirpy/app/internal/database"
	"github.com/charmbracelet/log"
)

func (cfg *apiConfig) loginUser(w http.ResponseWriter, req *http.Request) {
	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	// Get access token
	jwt, err := cfg.GenJWTkey(user.ID)
	if err != nil {
		log.Errorf("Error generatign jwt: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	// Get refresh token
	refreshToken, _ := auth.MakeRefreshToken()

	// Add refresh token to databse
	err = cfg.dbQueries.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		log.Errorf("Error inserting refresh token")
		responseWithErr(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	resp := Response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        jwt,
		RefreshToken: refreshToken,
	}

	responseWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) refreshToken(w http.ResponseWriter, req *http.Request) {
	// Get token from header
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Errorf("Invalid refresh token: %s", err)
		responseWithErr(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	// Get user refresh token details
	userDetails, err := cfg.dbQueries.GetRefreshToken(req.Context(), token)
	if err != nil {
		log.Errorf("Invalid refresh token: %s", err)
		responseWithErr(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Check if token is expired/revoked
	if time.Now().After(userDetails.ExpiresAt) || userDetails.RevokedAt.Valid {
		log.Errorf("token expired")
		responseWithErr(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Create new access token
	newAccessToken, err := auth.MakeJWT(userDetails.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Errorf("Unable to create new access token")
		responseWithErr(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: newAccessToken,
	})
}

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Errorf("Invalid refresh token: %s", err)
		responseWithErr(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	err = cfg.dbQueries.RevokeToken(req.Context(), token)
	if err != nil {
		log.Errorf("Invalid refresh token: %s", err)
		responseWithErr(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
