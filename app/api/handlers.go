package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DisguisedTrolley/chirpy/app/internal/database"
	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) getReqCount(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
		`, cfg.fileserverHits.Load())

	w.Write([]byte(body))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		responseWithErr(w, http.StatusForbidden, "Action forbidden")
		return
	}

	err := cfg.dbQueries.DeleteUsers(req.Context())
	if err != nil {
		log.Errorf("Unable to delete users: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Unable to delete users")
		return
	}
	cfg.fileserverHits.Store(0) // Reset visit counter

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func validateChirp(body string) (string, error) {
	if len(body) > 140 {
		return "", fmt.Errorf("Chirp too long")
	}

	cleanBody := cleanProfanity(body)

	return cleanBody, nil
}

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
