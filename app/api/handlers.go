package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
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

func validateChirp(w http.ResponseWriter, req *http.Request) {
	chirp := struct {
		Body string `json:"body"`
	}{}
	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&chirp)
	if err != nil {
		log.Errorf("Error decoding parameters: %s", err)
		responseWithErr(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		responseWithErr(w, http.StatusBadRequest, "Chirp too long")
		return
	}

	cleanBody := cleanProfanity(chirp.Body)

	responseWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": cleanBody})
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
