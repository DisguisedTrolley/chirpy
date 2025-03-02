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

func (cfg *apiConfig) handleReqCount(w http.ResponseWriter, req *http.Request) {
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

func (cfg *apiConfig) handleReqCountReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func handleChirps(w http.ResponseWriter, req *http.Request) {
	chirp := chirp{}
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
