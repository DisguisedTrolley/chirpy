package api

import (
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
